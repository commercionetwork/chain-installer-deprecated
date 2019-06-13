package apis

import (
	"github.com/commercionetwork/chain-installer/types"
	"github.com/commercionetwork/chain-installer/utils"
	"strings"
)

// GetChainsVersions allows to retrieve the list of all the chains that should be displayed to the user as the
// possibly installable chains.
func GetChainsVersions() []string {

	var content []types.RepoContent
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
