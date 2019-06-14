package interfaces

import "github.com/commercionetwork/chain-installer/types"

type Downloader interface {

	// Returns the information about the chain having the given chainName
	GetChainInfo(chainName string) types.ChainInfo

	// Downloads the genesis file for the given chain name and returns its content
	DownloadGenesisFile(info types.ChainInfo) string

	// Downloads the executables for the given chain name inside the specified directory
	DownloadExecutable(info types.ChainInfo, installationDir string)
}
