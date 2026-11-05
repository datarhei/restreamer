package api

// About is some general information about the API
type About struct {
	App       string   `json:"app"`
	Auths     []string `json:"auths"`
	Name      string   `json:"name"`
	ID        string   `json:"id"`
	CreatedAt string   `json:"created_at"`
	Uptime    uint64   `json:"uptime_seconds"`
	Version   Version  `json:"version"`
}

// Version is some information about the binary
type Version struct {
	Number   string `json:"number"`
	Commit   string `json:"repository_commit"`
	Branch   string `json:"repository_branch"`
	Build    string `json:"build_date"`
	Arch     string `json:"arch"`
	Compiler string `json:"compiler"`
}

// MinimalAbout is the minimal information about the API
type MinimalAbout struct {
	App     string         `json:"app"`
	Auths   []string       `json:"auths"`
	Version VersionMinimal `json:"version"`
}

type VersionMinimal struct {
	Number string `json:"number"`
}
