package api

// FileInfo represents informatiion about a file on a filesystem
type FileInfo struct {
	Name    string `json:"name" jsonschema:"minLength=1"`
	Size    int64  `json:"size_bytes" jsonschema:"minimum=0" format:"int64"`
	LastMod int64  `json:"last_modified" jsonschema:"minimum=0" format:"int64"`
}

// FilesystemInfo represents information about a filesystem
type FilesystemInfo struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Mount string `json:"mount"`
}
