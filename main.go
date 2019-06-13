package main

import (
	"errors"
	"fmt"
	"github.com/commercionetwork/chain-installer/apis"
	"github.com/commercionetwork/chain-installer/utils"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"strings"
)

func main() {

	fmt.Print("Welcome to the Commercio.network chain installer")

	// Ask the user to select a chain id
	chainId := getChainId()

	// Ask the user where to install the things
	installationDir := getInstallationDirectory()
	installationDir = utils.ReplaceLast(installationDir, "/", "")

	// Download the executable for the given chain id inside the given directory
	apis.DownloadChainExecutable(chainId, installationDir)

	// Download the genesis file inside the proper dir
	apis.DownloadGenesisFile(chainId, installationDir)

	// Ask the user to start the cnd or not
	if askStartCnd() {
		cmd := exec.Command(installationDir+"/cnd", "start")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		utils.CheckError(err)
	}
}

func getChainId() string {
	chains := apis.GetChainsVersions()

	chainPrompt := promptui.Select{
		Label: "Select the chain version you wish to install",
		Templates: &promptui.SelectTemplates{
			Selected: "Chain to be installed: {{ . }}",
		},
		Items: chains,
	}

	_, chain, err := chainPrompt.Run()
	utils.CheckError(err)

	return chain
}

func getInstallationDirectory() string {

	dirPrompt := promptui.Prompt{
		Label: "Installation directory",
		Templates: &promptui.PromptTemplates{
			Success: "Installation directory: ",
		},
		Validate: func(s string) error {
			if len(s) == 0 {
				return errors.New("The installation directory cannot be null")
			}
			return nil
		},
	}

	directory, err := dirPrompt.Run()
	utils.CheckError(err)

	return directory
}

func askStartCnd() bool {
	chainPrompt := promptui.Select{
		Label: "Do you wish to start your node now?",
		Items: []string{"Yes", "No"},
	}

	_, answer, err := chainPrompt.Run()
	utils.CheckError(err)

	return strings.EqualFold(answer, "yes")
}
