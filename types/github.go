package types

import (
	"errors"
	"github.com/commercionetwork/chain-installer/utils"
	"strings"
)

/// RepoContent represents a single item inside a Github repository
type RepoContents []RepoContent
type RepoContent struct {
	Name string `json:"name"`
	Type string `json:"type"` // Either dir if directory, or file if file
}

// Allows to filter a list of contents
func (items RepoContents) Filter(test func(RepoContent) bool) (ret RepoContents) {
	for _, item := range items {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return
}

// Allows to map a list of contents into a list of strings
func (items RepoContents) Map(mapper func(content RepoContent) string) (ret []string) {
	for _, item := range items {
		ret = append(ret, mapper(item))
	}
	return
}

// =====================================================

type Releases []Release
type Release struct {
	TagName string `json:"tag_name"`
	Url     string `json:"url"`
	Assets  Assets `json:"assets"`
}

// Allows to find a specific release having a given tag name associated
func (items Releases) FindByTagName(tagName string) Release {
	var release Release

	found := false
	for _, item := range items {
		if item.TagName == tagName {
			release = item
			found = true
		}
	}

	if !found {
		utils.CheckError(errors.New("item not found"))
	}

	return release
}

// =====================================================

type FileData struct {
	Name        string `json:"name"`
	DownloadUrl string `json:"download_url"`
}

// =====================================================

type Assets []Asset
type Asset struct {
	Name        string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
}

// Allows to find an asset having a given name inside the list
func (items Assets) FindByName(name string) Asset {
	var asset Asset
	found := false

	for _, item := range items {
		if strings.EqualFold(item.Name, name) {
			asset = item
			found = true
		}
	}

	if !found {
		utils.CheckError(errors.New("item not found"))
	}

	return asset
}
