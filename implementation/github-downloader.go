package implementation

import (
	sha2562 "crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/commercionetwork/chain-installer/apis"
	"github.com/commercionetwork/chain-installer/types"
	"github.com/commercionetwork/chain-installer/utils"
)

type GithubBasedDownloader struct {
	InstallationDir string

	Application types.Application
	ApiInfo     types.ApiInfo
}

func (downloader GithubBasedDownloader) getReleaseFolder(chainName string) string {
	// Get the contents of the specific chain folder
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s",
		downloader.ApiInfo.ChainsRepository.User,
		downloader.ApiInfo.ChainsRepository.RepoName,
		chainName)
}

func (downloader GithubBasedDownloader) GetChainInfo(chainName string) types.ChainInfo {
	// Get the URL where to find the .data file
	dataRemotePath := fmt.Sprintf("%s/.data", downloader.getReleaseFolder(chainName))

	// Download the seeds data
	var seedsData types.FileData
	apis.GetUrlContents(dataRemotePath, &seedsData)
	dataFileContents := apis.GetUrlContentsAsString(seedsData.DownloadUrl)

	// Parse the .data file content into a types.ChainInfo object
	chainInfo := downloader.readDataFileLines(dataFileContents)
	chainInfo.ChainName = chainName

	err := chainInfo.CheckValidity()
	utils.CheckError(err)

	return chainInfo
}

func (downloader GithubBasedDownloader) DownloadGenesisFile(info types.ChainInfo) string {
	fmt.Println("===> Getting the proper genesis file")

	// Get the genesis file information
	genesisRemotePath := fmt.Sprintf("%s/genesis.json", downloader.getReleaseFolder(info.ChainName))

	var genesisData types.FileData
	apis.GetUrlContents(genesisRemotePath, &genesisData)

	// Download the genesis.json file contents
	fmt.Println("===> Downloading the genesis file")
	genesisContents := apis.GetUrlContentsAsString(genesisData.DownloadUrl)

	// Check the validity of the contents
	sha256 := sha2562.Sum256([]byte(genesisContents))
	hexString := hex.EncodeToString(sha256[:])

	if hexString != info.GenesisChecksum {
		message := fmt.Sprintf("genesis.json checksum does not match downloaded genesis.json SHA256. Required %s but got %s instead",
			info.GenesisChecksum, hexString)
		panic(errors.New(message))
	}

	return genesisContents
}

func (downloader GithubBasedDownloader) DownloadExecutable(info types.ChainInfo, installationDir string) {
	fmt.Println("===> Downloading the chain executable")

	// Get the release version information
	releaseFileRemotePath := fmt.Sprintf("%s/.release", downloader.getReleaseFolder(info.ChainName))

	var releaseVersion types.FileData
	apis.GetUrlContents(releaseFileRemotePath, &releaseVersion)

	// === STEP 1 ===
	zipName, asset := downloader.getAssetsInfo(releaseVersion)

	// === STEP 2 ===
	downloadPath := downloader.downloadFiles(asset, zipName)

	// === STEP 3 ===
	downloadedFolderPath := downloader.unzipAndSetup(downloadPath, asset)

	// === STEP 4 ===
	cleanupInstallationFiles(downloadPath, downloadedFolderPath)

	// === STEP 5 ===
	fmt.Println("===> Executable downloaded successfully")
}
