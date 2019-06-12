package apis

import (
	"encoding/json"
	"fmt"
	"github.com/commercionetwork/chain-installer/types"
	"github.com/commercionetwork/chain-installer/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func getUrl(url string) io.ReadCloser {
	// Get the list of all the items inside the chains repository
	response, err := http.Get(url)
	utils.CheckError(err)
	return response.Body
}

func getUrlContentsAsString(url string) string {
	responseBody := getUrl(url)
	contents, err := ioutil.ReadAll(responseBody)
	utils.CheckError(err)
	return string(contents)
}

func getUrlContents(url string, destination interface{}) {
	responseBody := getUrl(url)
	err := json.NewDecoder(responseBody).Decode(&destination)
	utils.CheckError(err)
}

func GetChainsVersions() []string {

	var content [] types.RepoContent
	getUrlContents("https://api.github.com/repos/commercionetwork/chains/contents", &content)

	// Filter all the items that are a directory and have the name starting with
	folders := utils.FilterContent(content, func(c types.RepoContent) bool {
		return c.Type == "dir" && strings.HasPrefix(c.Name, "commercio-")
	})

	chains := utils.MapContent(folders, func(content types.RepoContent) string {
		return content.Name
	})

	return chains
}


func DownloadChainExecutable(chainId string, installationDir string) {

	fmt.Println("===> Downloading the chain executable <===")

	installationDir = utils.ReplaceLast(installationDir, "/", "")

	// Get the contents of the specific chain folder
	tagFolder := fmt.Sprintf("https://api.github.com/repos/commercionetwork/chains/contents/%s", chainId)

	// Get the data of the file containing the release version
	var releaseVersion types.FileData
	getUrlContents(tagFolder + "/.release", &releaseVersion)


	// === STEP 1 ===
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


	// === STEP 2 ===
	fmt.Println("\n ====> Step 2 - Downloading the files <====")

	// Download the asset inside the installation directory
	downloadPath := installationDir + "/executables.zip"
	fmt.Printf("Downloading %s into %s \n", zipName, downloadPath)
	err := DownloadFile(downloadPath, asset.DownloadUrl)
	utils.CheckError(err)
	fmt.Println(fmt.Sprintf("%s successfully downloaded", zipName))


	// === STEP 3 ===
	fmt.Println("\n ====> Step 3 - Unzipping and setup <====")

	// Unzip the file inside the installation directory
	fmt.Println(fmt.Sprintf("Unzipping %s into %s", downloadPath, installationDir))
	err = Unzip(downloadPath, installationDir)
	utils.CheckError(err)
	fmt.Println(fmt.Sprintf("%s unzipped!", downloadPath))

	downloadedFolder := utils.ReplaceLast(asset.Name, ".zip", "")
	downloadedFolderPath := installationDir + "/" + downloadedFolder

	// Move the cncli
	fmt.Println(fmt.Sprintf("Moving cncli into %s", installationDir))
	err = os.Rename(downloadedFolderPath + "/cncli", installationDir + "/cncli")
	utils.CheckError(err)

	// Move the cnd
	fmt.Println(fmt.Sprintf("Moving cnd into %s", installationDir))
	err = os.Rename(downloadedFolderPath + "/cnd", installationDir + "/cnd")
	utils.CheckError(err)


	// === STEP 4 ===
	fmt.Println("\n ====> Step 4 - Cleanup <====")

	// Delete the useless folders
	fmt.Println("Deleting useless folders")
	err = os.Remove(downloadPath)
	utils.CheckError(err)

	err = os.Remove(downloadedFolderPath)
	utils.CheckError(err)


	// === STEP 5 ===
	fmt.Println("====> Step 5 - End <====")
	fmt.Println("Executable downloaded successfully!")
}

func DownloadGenesisFile(chainId, installationDir string) {

	installationDir = utils.ReplaceLast(installationDir, "/", "")

	// Get the user home
	home, err := os.UserHomeDir()
	utils.CheckError(err)

	// === STEP 0 ===
	fmt.Println("\n ====> Step 1 - Cleanup <====")
	fmt.Println("Removing old data")

	// Remove the old cnd folder
	err = os.RemoveAll(home + "/.cnd")
	utils.CheckError(err)

	// Run the init command to create the folders again
	command := fmt.Sprintf("%s/cnd", installationDir)
	_, err = exec.Command(command, "init", "pippo").CombinedOutput()
	utils.CheckError(err)

	// === STEP 1 ===
	fmt.Println("\n ====> Step 1 - Getting the info <====")
	fmt.Println("Retrieving the genesis file to be downloaded")

	// Get the contents of the specific chain folder
	tagFolder := fmt.Sprintf("https://api.github.com/repos/commercionetwork/chains/contents/%s", chainId)

	// Get the data of the file containing the genesis information
	var genesisData types.FileData
	getUrlContents(tagFolder + "/genesis.json", &genesisData)
	fmt.Println("Genesis file info read successfully!")


	// === STEP 2 ===
	fmt.Println("\n ====> Step 2 - Downloading the genesis file <====")
	genesisContents := getUrlContentsAsString(genesisData.DownloadUrl)

	// Create the config folder
	configFolder := home + "/.cnd/config"
	err = os.Mkdir(configFolder, os.ModePerm)

	// Create the local genesis file
	genesisFile, err := os.Create(configFolder + "/genesis.json")
	utils.CheckError(err)

	// Write the data inside the config file
	_, err = genesisFile.WriteString(genesisContents)
	utils.CheckError(err)

	fmt.Println("Genesis file downloaded successfully")

	// === STEP 3 ===
	fmt.Println("\n ====> Step 3 - Setup <=====")
	fmt.Println("Downloading the seeds")

	// Get the seeds
	var seedsData types.FileData
	getUrlContents(tagFolder + "/.seeds", &seedsData)
	seeds := getUrlContentsAsString(seedsData.DownloadUrl)
	fmt.Println("Seeds downloaded successfully! ")

	// Read the current config.toml file contents
	fmt.Println("Writing config.toml")
	configTomlFile := home + "/.cnd/config/config.toml"
	configContents, err := ioutil.ReadFile(configTomlFile)
	utils.CheckError(err)

	// Replace the seeds inside the config.toml file
	seedsValue := fmt.Sprintf("seeds = \"%s\"", seeds)
	configContentsWithSeeds := strings.ReplaceAll(string(configContents), "seeds = \"\"", seedsValue)
	err = ioutil.WriteFile(configTomlFile, []byte(configContentsWithSeeds), os.ModePerm)
	utils.CheckError(err)
	fmt.Println("Config.toml updated successfully")
}