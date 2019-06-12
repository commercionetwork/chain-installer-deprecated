package main

import (
	"errors"
	"fmt"
	"github.com/commercionetwork/chain-installer/apis"
	"github.com/commercionetwork/chain-installer/utils"
	"github.com/manifoldco/promptui"
	"os/exec"
	"strings"
)

func main() {

	fmt.Print("Welcome to the Commercio.network chain installer")

	// Ask the user to select a chain id
	chainId := getChainId()

	// Ask the user where to install the things
	installationDir := getInstallationDirectory()

	// Download the executable for the given chain id inside the given directory
	apis.DownloadChainExecutable(chainId, installationDir)

	// Download the genesis file inside the proper dir
	apis.DownloadGenesisFile(chainId, installationDir)

	// Ask the user to start the cnd or not
	if askStartCnd() {
		out, err := exec.Command(installationDir + "/cnd", "start").StdoutPipe()
		buff := make([]byte,10)
		var n int
		for err == nil {
			n,err = out.Read(buff)
			if n > 0{
				fmt.Printf("taken %d chars %s",n,string(buff[:n]))
			}
		}
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
		Items: []string { "Yes", "No" },
	}

	_, answer, err := chainPrompt.Run()
	utils.CheckError(err)

	return strings.EqualFold(answer, "yes")
}