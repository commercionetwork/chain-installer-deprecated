package apis

import (
	"fmt"
	"github.com/commercionetwork/chain-installer/types"
	"github.com/commercionetwork/chain-installer/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// DownloadGenesisFile allows to download and properly setup the genesis.json file for the given chainId.
// The given installationDir is used in order to properly locate the cnd executable that should be run to setup
// everything before copying the chain genesis.json to avoid any error
func DownloadGenesisFile(chainId, installationDir string) {

	// Get the user home
	home, err := os.UserHomeDir()
	utils.CheckError(err)

	// === STEP 0 ===
	cleanupExistingData(home, installationDir)

	// === STEP 1 ===
	tagFolder, genesisData := getGenesisFileInfo(chainId)

	// === STEP 2 ===
	downloadGenesisFile(genesisData, home)

	// === STEP 3 ===
	setupConfigFile(tagFolder, home, err)
}

func cleanupExistingData(home string, installationDir string) {
	fmt.Println("\n ====> Step 1 - Cleanup <====")
	fmt.Println("Removing old data")

	// Remove the old cnd folder
	err := os.RemoveAll(home + "/.cnd")
	utils.CheckError(err)

	// Run the init command to create the folders again
	command := fmt.Sprintf("%s/cnd", installationDir)
	_, err = exec.Command(command, "init", "pippo").CombinedOutput()
	utils.CheckError(err)
}

func getGenesisFileInfo(chainId string) (string, types.FileData) {
	fmt.Println("\n ====> Step 1 - Getting the info <====")
	fmt.Println("Retrieving the genesis file to be downloaded")

	// Get the contents of the specific chain folder
	tagFolder := fmt.Sprintf("https://api.github.com/repos/commercionetwork/chains/contents/%s", chainId)

	// Get the data of the file containing the genesis information
	var genesisData types.FileData
	getUrlContents(tagFolder+"/genesis.json", &genesisData)

	fmt.Println("Genesis file info read successfully!")

	return tagFolder, genesisData
}

func downloadGenesisFile(genesisData types.FileData, home string) {

	fmt.Println("\n ====> Step 2 - Downloading the genesis file <====")
	genesisContents := getUrlContentsAsString(genesisData.DownloadUrl)

	// Create the config folder
	configFolder := home + "/.cnd/config"
	err := os.Mkdir(configFolder, os.ModePerm)

	// Create the local genesis file
	genesisFile, err := os.Create(configFolder + "/genesis.json")
	utils.CheckError(err)

	// Write the data inside the config file
	_, err = genesisFile.WriteString(genesisContents)
	utils.CheckError(err)

	fmt.Println("Genesis file downloaded successfully")
}

func setupConfigFile(tagFolder string, home string, err error) {
	fmt.Println("\n ====> Step 3 - Setup <=====")
	fmt.Println("Downloading the seeds")

	// Get the seeds
	var seedsData types.FileData
	getUrlContents(tagFolder+"/.seeds", &seedsData)
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

	// Replace the persistent peers inside the config.toml file
	// TODO - Remove this when the seed node will be successfully setup
	persistentPeers := "persistent_peers = \"337aae4a3976ad9fd28c051330d2eee58c721c23@3.121.145.63:26656,195486770d8de78846589ab288c1f2224e8429d6@52.29.124.231:26656,26c17a2102f5337aac02067559a3390ca73c6c42@3.122.146.155:26656\""
	configContentsWithPeers := strings.ReplaceAll(configContentsWithSeeds, "persistent_peers = \"\"", persistentPeers)

	// Write the contents back into the file
	err = ioutil.WriteFile(configTomlFile, []byte(configContentsWithPeers), os.ModePerm)
	utils.CheckError(err)

	fmt.Println("Config.toml updated successfully")
}
