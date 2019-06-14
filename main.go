package main

import (
	"errors"
	"fmt"
	"github.com/commercionetwork/chain-installer/implementation"
	"github.com/commercionetwork/chain-installer/interfaces"
	"github.com/commercionetwork/chain-installer/types"
	"github.com/commercionetwork/chain-installer/utils"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	// =========================
	// === Build the ApiInfo ===
	apiInfo := types.ApiInfo{
		ExecutablesRepository: types.ExecutablesRepoInfo{
			User:     "commercionetwork",
			RepoName: "commercionetwork",
		},
		ChainsRepository: types.ChainsRepoInfo{
			User:                 "commercionetwork",
			RepoName:             "chains",
			ValidChainNamePrefix: "commercio-",
		},
	}

	// =========================================
	// === Ask the user to select a chain id ===

	explorer := implementation.GithubVersionsExplorer{
		ApiInfo: apiInfo,
	}
	chainId := getChainId(explorer)

	// === Ask the user where to install the things ===
	installationDir := getInstallationDirectory()

	// Get the absolute path to the installation dir
	installationDir, err := filepath.Abs(installationDir)
	utils.CheckError(err)

	// Create the installation dir if it does not exist
	err = os.MkdirAll(installationDir, os.ModePerm)
	utils.CheckError(err)

	// =============================
	// === Build the coordinator ===

	application := types.Application{
		DaemonName: "cnd",
	}

	coordinator := interfaces.Coordinator{
		InstallationDir: installationDir,
		ChainName:       chainId,
		Application:     application,
		Downloader: implementation.GithubBasedDownloader{
			InstallationDir: installationDir,
			Application:     application,
			ApiInfo:         apiInfo,
		},
	}

	// =========================
	// === Download the data ===

	coordinator.PerformChainDownloadAndSetup()
	fmt.Println(fmt.Sprintf("Software installed into %s", installationDir))

	// =============================================
	// === Ask the user to start the node or not ===

	if askStartCnd() {
		cmd := exec.Command(installationDir+"/cnd", "start")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		utils.CheckError(err)
	}
}

// getChainId allows us to ask the user which chain he would like to install from the different chains available.
// Once the user has selected the chain, the selected id is returned.
func getChainId(versionsExplorer interfaces.VersionsExplorer) string {
	chainPrompt := promptui.Select{
		Label: "Select the chain version you wish to install",
		Templates: &promptui.SelectTemplates{
			Selected: "Chain to be installed: {{ . }}",
		},
		Items: versionsExplorer.ListChainNames(),
	}

	_, chain, err := chainPrompt.Run()
	utils.CheckError(err)

	return chain
}

// getInstallationDirectory asks the user where to install al the executable files and returns the prompted
// installation path.
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

// askStartCnd asks the user if he wants to start the installed full node or not, and returns a boolean indicating the
// chosen option
func askStartCnd() bool {
	chainPrompt := promptui.Select{
		Label: "Do you wish to start your node now?",
		Items: []string{"Yes", "No"},
	}

	_, answer, err := chainPrompt.Run()
	utils.CheckError(err)

	return strings.EqualFold(answer, "yes")
}
