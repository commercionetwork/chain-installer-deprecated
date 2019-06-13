package types

// Represents the Github repository inside which all the different chains versions data are hosted.
// The repo should have a series of folders all starting with the specified ValidChainNamePrefix, inside which
// there should be:
// 1. A .data file containing the information about the associated executable files
// 2. A genesis.json file
type ChainsRepoInfo struct {
	User                 string
	RepoName             string
	ValidChainNamePrefix string
}

// Represents the Github repository inside which all the different executables of your chain can be found.
// This repository must have a series of releases with the same tag name as the ones specified inside the various
// .data files that can be found in the chains repository.
type ExecutablesRepoInfo struct {
	User     string
	RepoName string
}

// Contains all the information about the different Github repositories that will be used while downloading the
// executables and the useful data of your chain.
type ApiInfo struct {
	ChainsRepository      ChainsRepoInfo
	ExecutablesRepository ExecutablesRepoInfo
}
