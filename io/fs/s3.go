package fs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/datarhei/core/v16/glob"
	"github.com/datarhei/core/v16/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Config struct {
	// Namee is the name of the filesystem
	Name            string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
	UseSSL          bool

	Logger log.Logger
}

type s3Filesystem struct {
	metadata map[string]string
	metaLock sync.RWMutex

	name string

	endpoint        string
	accessKeyID     string
	secretAccessKey string
	region          string
	bucket          string
	useSSL          bool

	client *minio.Client

	logger log.Logger
}

var fakeDirEntry = "..."

func NewS3Filesystem(config S3Config) (Filesystem, error) {
	fs := &s3Filesystem{
		metadata:        make(map[string]string),
		name:            config.Name,
		endpoint:        config.Endpoint,
		accessKeyID:     config.AccessKeyID,
		secretAccessKey: config.SecretAccessKey,
		region:          config.Region,
		bucket:          config.Bucket,
		useSSL:          config.UseSSL,
		logger:          config.Logger,
	}

	if fs.logger == nil {
		fs.logger = log.New("")
	}

	client, err := minio.New(fs.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(fs.accessKeyID, fs.secretAccessKey, ""),
		Region: fs.region,
		Secure: fs.useSSL,
	})

	if err != nil {
		return nil, fmt.Errorf("can't connect to s3 endpoint %s: %w", fs.endpoint, err)
	}

	fs.logger = fs.logger.WithFields(log.Fields{
		"name":     fs.name,
		"type":     "s3",
		"bucket":   fs.bucket,
		"region":   fs.region,
		"endpoint": fs.endpoint,
	})

	fs.logger.Debug().Log("Connected")

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	exists, err := client.BucketExists(ctx, fs.bucket)
	if err != nil {
		fs.logger.WithError(err).Log("Can't access bucket")
		return nil, fmt.Errorf("can't access bucket %s: %w", fs.bucket, err)
	}

	if exists {
		fs.logger.Debug().Log("Bucket already exists")
	} else {
		fs.logger.Debug().Log("Bucket doesn't exists")
		err = client.MakeBucket(ctx, fs.bucket, minio.MakeBucketOptions{Region: fs.region})
		if err != nil {
			fs.logger.WithError(err).Log("Can't create bucket")
			return nil, fmt.Errorf("can't create bucket %s: %w", fs.bucket, err)
		} else {
			fs.logger.Debug().Log("Bucket created")
		}
	}

	fs.client = client

	return fs, nil
}

func (fs *s3Filesystem) Name() string {
	return fs.name
}

func (fs *s3Filesystem) Type() string {
	return "s3"
}

func (fs *s3Filesystem) Metadata(key string) string {
	fs.metaLock.RLock()
	defer fs.metaLock.RUnlock()

	return fs.metadata[key]
}

func (fs *s3Filesystem) SetMetadata(key, data string) {
	fs.metaLock.Lock()
	defer fs.metaLock.Unlock()

	fs.metadata[key] = data
}

func (fs *s3Filesystem) Size() (int64, int64) {
	size := int64(0)

	files := fs.List("/", "")

	for _, file := range files {
		size += file.Size()
	}

	return size, -1
}

func (fs *s3Filesystem) Files() int64 {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := fs.client.ListObjects(ctx, fs.bucket, minio.ListObjectsOptions{
		WithVersions: false,
		WithMetadata: false,
		Prefix:       "",
		Recursive:    true,
		MaxKeys:      0,
		StartAfter:   "",
		UseV1:        false,
	})

	nfiles := int64(0)

	for object := range ch {
		if object.Err != nil {
			fs.logger.WithError(object.Err).Log("Listing object failed")
		}

		if strings.HasSuffix("/"+object.Key, "/"+fakeDirEntry) {
			// Skip fake entries (see MkdirAll)
			continue
		}

		nfiles++
	}

	return nfiles
}

func (fs *s3Filesystem) Symlink(oldname, newname string) error {
	return fmt.Errorf("not implemented")
}

func (fs *s3Filesystem) Stat(path string) (FileInfo, error) {
	path = fs.cleanPath(path)

	if len(path) == 0 {
		return &s3FileInfo{
			name:         "/",
			size:         0,
			dir:          true,
			lastModified: time.Now(),
		}, nil
	}

	ctx := context.Background()

	object, err := fs.client.GetObject(ctx, fs.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		if fs.isDir(path) {
			return &s3FileInfo{
				name:         "/" + path,
				size:         0,
				dir:          true,
				lastModified: time.Now(),
			}, nil
		}

		fs.logger.Debug().WithField("key", path).WithError(err).Log("Not found")
		return nil, err
	}

	defer object.Close()

	stat, err := object.Stat()
	if err != nil {
		if fs.isDir(path) {
			return &s3FileInfo{
				name:         "/" + path,
				size:         0,
				dir:          true,
				lastModified: time.Now(),
			}, nil
		}

		fs.logger.Debug().WithField("key", path).WithError(err).Log("Stat failed")
		return nil, err
	}

	return &s3FileInfo{
		name:         "/" + stat.Key,
		size:         stat.Size,
		lastModified: stat.LastModified,
	}, nil
}

func (fs *s3Filesystem) Open(path string) File {
	path = fs.cleanPath(path)
	ctx := context.Background()

	object, err := fs.client.GetObject(ctx, fs.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		fs.logger.Debug().WithField("key", path).Log("Not found")
		return nil
	}

	stat, err := object.Stat()
	if err != nil {
		fs.logger.Debug().WithField("key", path).Log("Stat failed")
		return nil
	}

	file := &s3File{
		data:         object,
		name:         "/" + stat.Key,
		size:         stat.Size,
		lastModified: stat.LastModified,
	}

	fs.logger.Debug().WithField("key", stat.Key).Log("Opened")

	return file
}

func (fs *s3Filesystem) ReadFile(path string) ([]byte, error) {
	path = fs.cleanPath(path)
	file := fs.Open(path)
	if file == nil {
		return nil, os.ErrNotExist
	}

	defer file.Close()

	buf := &bytes.Buffer{}

	_, err := buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (fs *s3Filesystem) write(path string, r io.Reader) (int64, bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	overwrite := false

	_, err := fs.client.StatObject(ctx, fs.bucket, path, minio.StatObjectOptions{})
	if err == nil {
		overwrite = true
	}

	info, err := fs.client.PutObject(ctx, fs.bucket, path, r, -1, minio.PutObjectOptions{
		UserMetadata:            map[string]string{},
		UserTags:                map[string]string{},
		Progress:                nil,
		ContentType:             "",
		ContentEncoding:         "",
		ContentDisposition:      "",
		ContentLanguage:         "",
		CacheControl:            "",
		Mode:                    "",
		RetainUntilDate:         time.Time{},
		ServerSideEncryption:    nil,
		NumThreads:              0,
		StorageClass:            "",
		WebsiteRedirectLocation: "",
		PartSize:                0,
		LegalHold:               "",
		SendContentMd5:          false,
		DisableContentSha256:    false,
		DisableMultipart:        false,
		Internal:                minio.AdvancedPutOptions{},
	})
	if err != nil {
		fs.logger.WithError(err).WithField("key", path).Log("Failed to store file")
		return -1, false, err
	}

	fs.logger.Debug().WithFields(log.Fields{
		"key":       path,
		"overwrite": overwrite,
	}).Log("Stored")

	return info.Size, !overwrite, nil
}

func (fs *s3Filesystem) WriteFileReader(path string, r io.Reader) (int64, bool, error) {
	path = fs.cleanPath(path)
	return fs.write(path, r)
}

func (fs *s3Filesystem) WriteFile(path string, data []byte) (int64, bool, error) {
	return fs.WriteFileReader(path, bytes.NewBuffer(data))
}

func (fs *s3Filesystem) WriteFileSafe(path string, data []byte) (int64, bool, error) {
	return fs.WriteFileReader(path, bytes.NewBuffer(data))
}

func (fs *s3Filesystem) Rename(src, dst string) error {
	src = fs.cleanPath(src)
	dst = fs.cleanPath(dst)

	err := fs.Copy(src, dst)
	if err != nil {
		return err
	}

	res := fs.Remove(src)
	if res == -1 {
		return fmt.Errorf("failed to remove source file: %s", src)
	}

	return nil
}

func (fs *s3Filesystem) Copy(src, dst string) error {
	src = fs.cleanPath(src)
	dst = fs.cleanPath(dst)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := fs.client.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: fs.bucket,
		Object: dst,
	}, minio.CopySrcOptions{
		Bucket: fs.bucket,
		Object: src,
	})

	return err
}

func (fs *s3Filesystem) MkdirAll(path string, perm os.FileMode) error {
	if path == "/" {
		return nil
	}

	info, err := fs.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return os.ErrExist
		}

		return nil
	}

	path = filepath.Join(path, fakeDirEntry)

	_, _, err = fs.write(path, strings.NewReader(""))
	if err != nil {
		return fmt.Errorf("can't create directory")
	}

	return nil
}

func (fs *s3Filesystem) Remove(path string) int64 {
	path = fs.cleanPath(path)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stat, err := fs.client.StatObject(ctx, fs.bucket, path, minio.StatObjectOptions{})
	if err != nil {
		fs.logger.Debug().WithField("key", path).Log("Not found")
		return -1
	}

	err = fs.client.RemoveObject(ctx, fs.bucket, path, minio.RemoveObjectOptions{
		GovernanceBypass: true,
	})
	if err != nil {
		fs.logger.WithError(err).WithField("key", stat.Key).Log("Failed to delete file")
		return -1
	}

	fs.logger.Debug().WithField("key", stat.Key).Log("Deleted")

	return stat.Size
}

func (fs *s3Filesystem) RemoveAll() int64 {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	totalSize := int64(0)

	objectsCh := make(chan minio.ObjectInfo)

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)

		for object := range fs.client.ListObjects(ctx, fs.bucket, minio.ListObjectsOptions{
			Recursive: true,
		}) {
			if object.Err != nil {
				fs.logger.WithError(object.Err).Log("Listing object failed")
				continue
			}
			totalSize += object.Size
			objectsCh <- object
		}
	}()

	for err := range fs.client.RemoveObjects(context.Background(), fs.bucket, objectsCh, minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}) {
		fs.logger.WithError(err.Err).WithField("key", err.ObjectName).Log("Deleting object failed")
	}

	fs.logger.Debug().Log("Deleted all files")

	return totalSize
}

func (fs *s3Filesystem) List(path, pattern string) []FileInfo {
	path = fs.cleanPath(path)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := fs.client.ListObjects(ctx, fs.bucket, minio.ListObjectsOptions{
		WithVersions: false,
		WithMetadata: false,
		Prefix:       path,
		Recursive:    true,
		MaxKeys:      0,
		StartAfter:   "",
		UseV1:        false,
	})

	files := []FileInfo{}

	for object := range ch {
		if object.Err != nil {
			fs.logger.WithError(object.Err).Log("Listing object failed")
			continue
		}

		key := "/" + object.Key
		if strings.HasSuffix(key, "/"+fakeDirEntry) {
			// filter out fake directory entries (see MkdirAll)
			continue
		}

		if len(pattern) != 0 {
			if ok, _ := glob.Match(pattern, key, '/'); !ok {
				continue
			}
		}

		f := &s3FileInfo{
			name:         key,
			size:         object.Size,
			lastModified: object.LastModified,
		}

		files = append(files, f)
	}

	return files
}

func (fs *s3Filesystem) LookPath(file string) (string, error) {
	if strings.Contains(file, "/") {
		file = fs.cleanPath(file)
		info, err := fs.Stat(file)
		if err == nil {
			if !info.Mode().IsRegular() {
				return file, os.ErrNotExist
			}
			return file, nil
		}
		return "", os.ErrNotExist
	}
	path := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(path) {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		path := filepath.Join(dir, file)
		path = fs.cleanPath(path)
		if info, err := fs.Stat(path); err == nil {
			if !filepath.IsAbs(path) {
				return path, os.ErrNotExist
			}
			if !info.Mode().IsRegular() {
				return path, os.ErrNotExist
			}
			return path, nil
		}
	}
	return "", os.ErrNotExist
}

func (fs *s3Filesystem) isDir(path string) bool {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	if path == "/" {
		return true
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := fs.client.ListObjects(ctx, fs.bucket, minio.ListObjectsOptions{
		WithVersions: false,
		WithMetadata: false,
		Prefix:       path,
		Recursive:    true,
		MaxKeys:      1,
		StartAfter:   "",
		UseV1:        false,
	})

	files := uint64(0)

	for object := range ch {
		if object.Err != nil {
			fs.logger.WithError(object.Err).Log("Listing object failed")
			continue
		}

		files++
	}

	return files > 0
}

func (fs *s3Filesystem) cleanPath(path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join("/", path)
	}

	path = strings.TrimSuffix(path, "/"+fakeDirEntry)

	return filepath.Join("/", filepath.Clean(path))[1:]
}

type s3FileInfo struct {
	name         string
	size         int64
	dir          bool
	lastModified time.Time
}

func (f *s3FileInfo) Name() string {
	return f.name
}

func (f *s3FileInfo) Size() int64 {
	return f.size
}

func (f *s3FileInfo) Mode() os.FileMode {
	return fs.FileMode(fs.ModePerm)
}

func (f *s3FileInfo) ModTime() time.Time {
	return f.lastModified
}

func (f *s3FileInfo) IsLink() (string, bool) {
	return "", false
}

func (f *s3FileInfo) IsDir() bool {
	return f.dir
}

type s3File struct {
	data         io.ReadCloser
	name         string
	size         int64
	lastModified time.Time
}

func (f *s3File) Read(p []byte) (int, error) {
	return f.data.Read(p)
}

func (f *s3File) Close() error {
	return f.data.Close()
}

func (f *s3File) Name() string {
	return f.name
}

func (f *s3File) Stat() (FileInfo, error) {
	return &s3FileInfo{
		name:         f.name,
		size:         f.size,
		lastModified: f.lastModified,
	}, nil
}
