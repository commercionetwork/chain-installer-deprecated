package utils

import (
	"errors"
	"github.com/commercionetwork/chain-installer/types"
	"os"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetUserHome() string {
	home, err := os.UserHomeDir()
	CheckError(err)
	return home
}

func ReplaceLast(original, old, replace string) string {
	i := strings.LastIndex(original, old)

	if i >= 0 {
		return original[:i] + strings.Replace(original[i:], old, replace, 1)
	} else {
		return original
	}
}

// === Contents ===

func FilterContent(items []types.RepoContent, test func(types.RepoContent) bool) (ret []types.RepoContent) {
	for _, item := range items {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return
}

func MapContent(items []types.RepoContent, mapper func(content types.RepoContent) string) (ret []string) {
	for _, item := range items {
		ret = append(ret, mapper(item))
	}
	return
}

// === Releases ===

func FindReleaseByTagName(items []types.Release, tagName string) types.Release {
	var release types.Release

	found := false
	for _, item := range items {
		if item.TagName == tagName {
			release = item
			found = true
		}
	}

	if !found {
		CheckError(errors.New("item not found"))
	}

	return release
}

// === Assets ===

func FindReleaseAssetByName(items []types.Asset, name string) types.Asset {
	var asset types.Asset
	found := false

	for _, item := range items {
		if strings.EqualFold(item.Name, name) {
			asset = item
			found = true
		}
	}

	if !found {
		CheckError(errors.New("item not found"))
	}

	return asset
}
