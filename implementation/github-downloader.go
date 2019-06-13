package implementation

import (
	"fmt"
	"github.com/commercionetwork/chain-installer/apis"
	"github.com/commercionetwork/chain-installer/types"
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

func (downloader GithubBasedDownloader) DownloadGenesisFile(chainName string) string {
	fmt.Println("===> Getting the proper genesis file")

	// Get the genesis file information
	genesisRemotePath := fmt.Sprintf("%s/genesis.json", downloader.getReleaseFolder(chainName))

	var genesisData types.FileData
	apis.GetUrlContents(genesisRemotePath, &genesisData)

	fmt.Println("===> Downloading the genesis file")
	return apis.GetUrlContentsAsString(genesisData.DownloadUrl)
}

func (downloader GithubBasedDownloader) GetSeeds(chainName string) string {
	fmt.Println("===> Downloading the seeds")

	// Get the URL where to find the .seeds file
	seedsRemotePath := fmt.Sprintf("%s/.seeds", downloader.getReleaseFolder(chainName))

	// Download the seeds data
	var seedsData types.FileData
	apis.GetUrlContents(seedsRemotePath, &seedsData)
	return apis.GetUrlContentsAsString(seedsData.DownloadUrl)
}

func (downloader GithubBasedDownloader) DownloadExecutable(chainName string, installationDir string) {
	fmt.Println("===> Downloading the chain executable")

	// Get the release version information
	releaseFileRemotePath := fmt.Sprintf("%s/.release", downloader.getReleaseFolder(chainName))

	var releaseVersion types.FileData
	apis.GetUrlContents(releaseFileRemotePath, &releaseVersion)

	// === STEP 1 ===
	zipName, asset := downloader.getAssetsInfo(releaseVersion)

	// === STEP 2 ===
	downloadPath := downloader.downloadFiles(zipName, asset)

	// === STEP 3 ===
	downloader.unzipAndSetup(downloadPath, asset)

	// === STEP 4 ===
	//cleanupInstallationFiles(downloadPath, downloadedFolderPath)

	// === STEP 5 ===
	fmt.Println("===> Executable downloaded successfully")
}
