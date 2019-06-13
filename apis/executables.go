package apis

import (
	"fmt"
	"github.com/commercionetwork/chain-installer/types"
	"github.com/commercionetwork/chain-installer/utils"
	"os"
	"runtime"
)

// DownloadChainExecutable allows to download the executables files for the given chainId inside the specified
// installationDir
// All the files are properly fetched from the different Github repositories.
func DownloadChainExecutable(chainId string, installationDir string) {

	fmt.Println("===> Downloading the chain executable <===")

	// Get the contents of the specific chain folder
	tagFolder := fmt.Sprintf("https://api.github.com/repos/commercionetwork/chains/contents/%s", chainId)

	// Get the data of the file containing the release version
	var releaseVersion types.FileData
	getUrlContents(tagFolder+"/.release", &releaseVersion)

	// === STEP 1 ===
	zipName, asset := getAssetsInfo(releaseVersion)

	// === STEP 2 ===
	downloadPath := downloadFiles(installationDir, zipName, asset)

	// === STEP 3 ===
	downloadedFolderPath := unzipAndSetup(downloadPath, installationDir, asset)

	// === STEP 4 ===
	cleanupInstallationFiles(downloadPath, downloadedFolderPath)

	// === STEP 5 ===
	fmt.Println("====> Executable downloaded successfully! <====")
}

func getAssetsInfo(releaseVersion types.FileData) (string, types.Asset) {
	fmt.Println("\n ====> Step 1 - Obtain the info <====")

	// Get the .release and the genesis.json files contents
	fmt.Println("Reading the .release file contents")

	tagName := getUrlContentsAsString(releaseVersion.DownloadUrl)
	fmt.Println("Contents of the .release file successfully read!")

	// Get all the releases and find the one having the given tag name
	fmt.Println(fmt.Sprintf("Searching the release with tag name %s", tagName))
	var releases []types.Release
	getUrlContents("https://api.github.com/repos/commercionetwork/commercionetwork/releases", &releases)

	release := utils.FindReleaseByTagName(releases, tagName)
	fmt.Println(fmt.Sprintf("Release with tag name %s found!", tagName))

	// Get the asset representing the zip file inside which there are the executables for the given OS and Architecture
	zipName := fmt.Sprintf("%s-%s.zip", runtime.GOOS, runtime.GOARCH)
	fmt.Println(fmt.Sprintf("Searching the asset with name %s inside release %s", zipName, tagName))

	asset := utils.FindReleaseAssetByName(release.Assets, zipName)
	fmt.Println(fmt.Sprintf("Asset with name %s found!", zipName))

	return zipName, asset
}

func downloadFiles(installationDir string, zipName string, asset types.Asset) string {
	fmt.Println("\n ====> Step 2 - Downloading the files <====")

	// Download the asset inside the installation directory
	downloadPath := installationDir + "/executables.zip"

	fmt.Printf("Downloading %s into %s \n", zipName, downloadPath)
	err := DownloadFile(downloadPath, asset.DownloadUrl)
	utils.CheckError(err)

	fmt.Println(fmt.Sprintf("%s successfully downloaded", zipName))

	return downloadPath
}

func unzipAndSetup(downloadPath string, installationDir string, asset types.Asset) string {
	fmt.Println("\n ====> Step 3 - Unzipping and setup <====")

	// Unzip the file inside the installation directory
	fmt.Println(fmt.Sprintf("Unzipping %s into %s", downloadPath, installationDir))
	err := Unzip(downloadPath, installationDir)
	utils.CheckError(err)

	fmt.Println(fmt.Sprintf("%s unzipped!", downloadPath))

	downloadedFolder := utils.ReplaceLast(asset.Name, ".zip", "")
	downloadedFolderPath := installationDir + "/" + downloadedFolder

	// Move the cncli
	fmt.Println(fmt.Sprintf("Moving cncli into %s", installationDir))
	err = os.Rename(downloadedFolderPath+"/cncli", installationDir+"/cncli")
	utils.CheckError(err)

	// Move the cnd
	fmt.Println(fmt.Sprintf("Moving cnd into %s", installationDir))
	err = os.Rename(downloadedFolderPath+"/cnd", installationDir+"/cnd")
	utils.CheckError(err)

	return downloadedFolderPath
}

func cleanupInstallationFiles(downloadPath string, downloadedFolderPath string) {
	fmt.Println("\n ====> Step 4 - Cleanup <====")

	// Delete the useless folders
	fmt.Println("Deleting useless folders")
	err := os.Remove(downloadPath)
	utils.CheckError(err)

	err = os.Remove(downloadedFolderPath)
	utils.CheckError(err)
}
