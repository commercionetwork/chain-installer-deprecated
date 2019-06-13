package apis

import (
	"encoding/json"
	"github.com/commercionetwork/chain-installer/utils"
	"io"
	"io/ioutil"
	"net/http"
)

// getUrls allows to perform a GET request to the given url, and return the answer as an io.ReadCloser.
func getUrl(url string) io.ReadCloser {
	// Get the list of all the items inside the chains repository
	response, err := http.Get(url)
	utils.CheckError(err)
	return response.Body
}

// getUrlContentsAsString allows to perform a GET request to the given url and return the server answer as a string.
func getUrlContentsAsString(url string) string {
	responseBody := getUrl(url)
	contents, err := ioutil.ReadAll(responseBody)
	utils.CheckError(err)
	return string(contents)
}

// getUrlContents performs a GET request to the given url and later parses the response and puts it into the destination
// interface
func getUrlContents(url string, destination interface{}) {
	responseBody := getUrl(url)
	err := json.NewDecoder(responseBody).Decode(&destination)
	utils.CheckError(err)
}
