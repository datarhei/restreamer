package fs

import (
	"bytes"

	"github.com/puzpuzpuz/xsync/v3"
)

type memStorage struct {
	lock  *xsync.RBMutex
	files *xsync.MapOf[string, *memFile]
}

func newMemStorage() *memStorage {
	m := &memStorage{
		lock:  xsync.NewRBMutex(),
		files: xsync.NewMapOf[string, *memFile](),
	}

	return m
}

func (m *memStorage) Delete(key string) (*memFile, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.files.LoadAndDelete(key)
}

func (m *memStorage) Store(key string, value *memFile) (*memFile, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.files.LoadAndStore(key, value)
}

func (m *memStorage) Load(key string) (*memFile, bool) {
	token := m.lock.RLock()
	defer m.lock.RUnlock(token)

	return m.files.Load(key)
}

func (m *memStorage) LoadAndCopy(key string) (*memFile, bool) {
	token := m.lock.RLock()
	defer m.lock.RUnlock(token)

	v, ok := m.files.Load(key)
	if !ok {
		return nil, false
	}

	f := &memFile{
		memFileInfo: memFileInfo{
			name:    v.name,
			size:    v.size,
			dir:     v.dir,
			lastMod: v.lastMod,
			linkTo:  v.linkTo,
		},
	}

	if v.data != nil {
		f.data = bytes.NewBuffer(v.data.Bytes())
	}

	return f, true
}

func (m *memStorage) Has(key string) bool {
	token := m.lock.RLock()
	defer m.lock.RUnlock(token)

	_, ok := m.files.Load(key)

	return ok
}

func (m *memStorage) Range(f func(key string, value *memFile) bool) {
	token := m.lock.RLock()
	defer m.lock.RUnlock(token)

	m.files.Range(f)
}
