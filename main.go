package gist4737109

import (
	"encoding/json"
	"fmt"

	. "gist.github.com/4668739.git"
)

// This assumes there's only one file in the gist
func GistIdToGistContents(gistId string) string {
	return GistIdCommitIdToGistContents(gistId, "")
}

// This assumes there's only one file in the gist
func GistIdCommitIdToGistContents(gistId, commitId string) string {
	gistUrl := "https://api.github.com/gists/" + gistId
	if "" != commitId {
		gistUrl += "/" + commitId
	}

	var gistJson struct {
		Files map[string]struct{ Raw_Url string }
	}
	err := json.Unmarshal(HttpGetB(gistUrl), &gistJson)
	if err != nil {
		fmt.Println("error:", err)
	}
	for _, v := range gistJson.Files {
		return HttpGet(v.Raw_Url)
	}
	return ""
}

func GistIdToUsername(gistId string) (string, error) {
	gistUrl := "https://api.github.com/gists/" + gistId

	var gistJson struct {
		User struct{ Login string }
	}
	err := json.Unmarshal(HttpGetB(gistUrl), &gistJson)
	if err != nil {
		return "", err
	}
	return gistJson.User.Login, nil
}

func main() {
	gistId := "4737109"
	commitId := "1b4e4b0e6f469d5e5a91b49028fbf2ab936bfcd4"

	fmt.Println(GistIdToUsername(gistId))
	println(GistIdCommitIdToGistContents(gistId, commitId))
	println(GistIdToGistContents(gistId))
}
