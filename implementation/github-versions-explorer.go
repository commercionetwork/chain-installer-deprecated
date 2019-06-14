package implementation

import (
	"fmt"
	"github.com/commercionetwork/chain-installer/apis"
	"github.com/commercionetwork/chain-installer/types"
	"strings"
)

// Implementation of the VersionsExplorer interface allowing to list the available chains from a GitHub repository
type GithubVersionsExplorer struct {
	ApiInfo types.ApiInfo
}

// listChainNames implements the VersionsExplorer interface
func (explorer GithubVersionsExplorer) ListChainNames() []string {

	chainApiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents",
		explorer.ApiInfo.ChainsRepository.User,
		explorer.ApiInfo.ChainsRepository.RepoName)

	var content types.RepoContents
	apis.GetUrlContents(chainApiUrl, &content)

	// Filter all the items that are a directory and have the name starting with
	folders := content.Filter(func(c types.RepoContent) bool {
		return c.Type == "dir" && strings.HasPrefix(c.Name, explorer.ApiInfo.ChainsRepository.ValidChainNamePrefix)
	})

	chains := folders.Map(func(content types.RepoContent) string {
		return content.Name
	})

	return chains
}
