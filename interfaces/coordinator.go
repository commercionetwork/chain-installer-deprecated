package interfaces

import (
	"fmt"
	"github.com/commercionetwork/chain-installer/types"
	"github.com/commercionetwork/chain-installer/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Represents the coordinator of the download and setup of all the chain executables and configuration files
type Coordinator struct {
	InstallationDir string
	ChainName       string
	Application     types.Application
	Downloader      Downloader
}

// Downloads all the necessary chain files and setups everything so that it can work properly
func (coordinator Coordinator) PerformChainDownloadAndSetup() {

	// Download the executables inside the installation dir
	coordinator.downloadExecutables()

	// Reset the daemon folder
	coordinator.resetDaemonFolder()

	// Download the new genesis file
	coordinator.downloadGenesisFile()

	// Setup the config.toml file
	coordinator.setupConfigTomlFile()
}

func (coordinator Coordinator) downloadExecutables() {
	coordinator.Downloader.DownloadExecutable(coordinator.ChainName, coordinator.InstallationDir)
}

func (coordinator Coordinator) downloadGenesisFile() {
	genesisContents := coordinator.Downloader.DownloadGenesisFile(coordinator.ChainName)

	// Create the config folder
	configFolderPath := fmt.Sprintf("%s/.%s/config", utils.GetUserHome(), coordinator.Application.DaemonName)
	err := os.Mkdir(configFolderPath, os.ModePerm)

	// Create the local genesis file
	genesisFilePath := fmt.Sprintf("%s/genesis.json", configFolderPath)
	genesisFile, err := os.Create(genesisFilePath)
	utils.CheckError(err)

	// Write the data inside the config file
	_, err = genesisFile.WriteString(genesisContents)
	utils.CheckError(err)

	fmt.Println("===> Genesis file downloaded successfully")
}

func (coordinator Coordinator) resetDaemonFolder() {
	fmt.Println("===> Removing the existing node data")

	// Remove the old daemon folder
	daemonDataFolder := fmt.Sprintf("%s/.%s", utils.GetUserHome(), coordinator.Application.DaemonName)
	err := os.RemoveAll(daemonDataFolder)
	utils.CheckError(err)

	// Run the init command to create the folders again
	command := fmt.Sprintf("%s/%s", coordinator.InstallationDir, coordinator.Application.DaemonName)

	moniker := utils.GetRandomMoniker()
	_, err = exec.Command(command, "init", moniker).CombinedOutput()
	utils.CheckError(err)

	fmt.Println("===> Removed the existing node data")
}

func (coordinator Coordinator) setupConfigTomlFile() {
	fmt.Println("===> Writing config.toml")

	// Get the config.toml file path
	configTomlFilePath := fmt.Sprintf("%s/.%s/config/config.toml",
		utils.GetUserHome(),
		coordinator.Application.DaemonName)

	// Read the config.toml file contents
	configContents, err := ioutil.ReadFile(configTomlFilePath)
	utils.CheckError(err)

	// Get the seeds
	seeds := coordinator.Downloader.GetSeeds(coordinator.ChainName)

	// Replace the seeds inside the config.toml file
	seedsValue := fmt.Sprintf("seeds = \"%s\"", seeds)
	configContentsWithSeeds := strings.ReplaceAll(string(configContents), "seeds = \"\"", seedsValue)

	// Write the contents back into the file
	err = ioutil.WriteFile(configTomlFilePath, []byte(configContentsWithSeeds), os.ModePerm)
	utils.CheckError(err)

	fmt.Println("Config.toml updated successfully")
}
