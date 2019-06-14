package implementation

import (
	"bufio"
	"fmt"
	"github.com/commercionetwork/chain-installer/apis"
	"github.com/commercionetwork/chain-installer/types"
	"github.com/commercionetwork/chain-installer/utils"
	otiai10copy "github.com/otiai10/copy"
	"os"
	"regexp"
	"runtime"
	"strings"
)

func (downloader GithubBasedDownloader) readDataFileLines(contents string) types.ChainInfo {

	var chainInfo types.ChainInfo

	// Start reading from the file with a reader.
	reader := bufio.NewReader(strings.NewReader(contents))

	for {
		line, err := reader.ReadString('\n')
		parts := regexp.MustCompile("\\s{2,}").Split(line, 2)

		if len(parts) == 2 {
			value := strings.TrimSpace(parts[1])

			if contains(parts[0], "release") {
				chainInfo.ReleaseTag = value
			} else if contains(parts[0], "seeds") {
				chainInfo.Seeds = value
			} else if contains(parts[0], "persistent") && contains(parts[0], "peers") {
				chainInfo.PersistentPeers = value
			} else if contains(parts[0], "genesis") && contains(parts[0], "checksum") {
				chainInfo.GenesisChecksum = value
			}
		}

		if err != nil {
			break
		}
	}

	return chainInfo
}

func contains(original, search string) bool {
	return strings.Contains(strings.ToLower(original), strings.ToLower(search))
}

func (downloader GithubBasedDownloader) getAssetsInfo(releaseVersion types.FileData) (string, types.Asset) {
	fmt.Println("====> Getting the .release file")

	// Get the .release and the genesis.json files contents
	tagName := apis.GetUrlContentsAsString(releaseVersion.DownloadUrl)

	// Get all the releases and find the one having the given tag name
	fmt.Println(fmt.Sprintf("====> Searching the release with tag name %s", tagName))

	// Get the release API URL
	releaseApiPath := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases",
		downloader.ApiInfo.ExecutablesRepository.User,
		downloader.ApiInfo.ExecutablesRepository.RepoName)

	// Get the releases list
	var releases types.Releases
	apis.GetUrlContents(releaseApiPath, &releases)

	// Find the release with the given tag name
	release := releases.FindByTagName(tagName)

	// Get the asset representing the zip file inside which there are the executables for the given OS and Architecture
	zipName := fmt.Sprintf("%s-%s.zip", runtime.GOOS, runtime.GOARCH)
	fmt.Println(fmt.Sprintf("====> Searching the asset with name %s inside release %s", zipName, tagName))

	asset := release.Assets.FindByName(zipName)

	return zipName, asset
}

func (downloader GithubBasedDownloader) downloadFiles(asset types.Asset, zipName string) string {
	fmt.Println("====> Downloading the executable files")

	// Download the asset inside the installation directory
	downloadPath := fmt.Sprintf("%s/executables.zip", downloader.InstallationDir)

	err := apis.DownloadFile(downloadPath, asset.DownloadUrl)
	utils.CheckError(err)

	return downloadPath
}

func (downloader GithubBasedDownloader) unzipAndSetup(downloadPath string, asset types.Asset) string {
	fmt.Println("====> Unzipping the executable files")

	// Unzip the file inside the installation directory
	fmt.Println(fmt.Sprintf("====> Unzipping %s into %s", downloadPath, downloader.InstallationDir))
	err := apis.Unzip(downloadPath, downloader.InstallationDir)
	utils.CheckError(err)

	downloadedFolder := utils.ReplaceLast(asset.Name, ".zip", "")
	downloadedFolderPath := downloader.InstallationDir + "/" + downloadedFolder

	err = otiai10copy.Copy(downloadedFolderPath, downloader.InstallationDir)

	return downloadedFolderPath
}

func cleanupInstallationFiles(downloadPath string, downloadedFolderPath string) {
	fmt.Println("====> Performing cleanup")

	// Delete the useless folders
	err := os.RemoveAll(downloadPath)
	utils.CheckError(err)

	err = os.RemoveAll(downloadedFolderPath)
	utils.CheckError(err)
}
