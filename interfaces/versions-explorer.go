package interfaces

type VersionsExplorer interface {

	// Returns the names of all the different chains that the user can install
	ListChainNames() []string
}
