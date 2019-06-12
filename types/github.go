package types

/// RepoContent represents a single item inside a Github repository
type RepoContent struct {
	Name string `json:"name"`
	Type string `json:"type"`	// Either dir if directory, or file if file
}

type Release struct {
	TagName	string `json:"tag_name"`
	Url 	string `json:"url"`
	Assets 	[]Asset `json:"assets"`
}

type FileData struct {
	Name 		string `json:"name"`
	DownloadUrl string `json:"download_url"`
}

type Asset struct {
	Name 		string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
}