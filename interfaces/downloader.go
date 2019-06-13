package interfaces

type Downloader interface {

	// Downloads the genesis file for the given chain name and returns its content
	DownloadGenesisFile(chainName string) string

	// Downloads the executables for the given chain name inside the specified directory
	DownloadExecutable(chainName string, installationDir string)

	// Returns the list of seed nodes that should be put inside the config.toml file for the given chain name
	GetSeeds(chainName string) string
}
