package twitter

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
)

const endpoint = "https://api.twitter.com/1.1"

func UserShow(name string) {
	url := fmt.Sprintf("%s/users/show.json?screen_name=%s", endpoint, name)

	response, _ := util.Get(url, headers)
}

func UserTimeline(name string, count int, useCache bool) {

}
