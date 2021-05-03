package kraken

import (
	"fmt"
	"strings"

	scribble "github.com/nanobox-io/golang-scribble"
)

func (api *api) getAPIAssets() {
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Kraken/public", "Assets", &api.assets.Result)
	}
	if !useCache || err != nil {
		resource := "/0/public/Assets"
		headers := make(map[string]string)
		headers["Content-Type"] = "application/json"
		resp, err := api.clientAssets.R().
			SetHeaders(headers).
			SetResult(&AssetsInfo{}).
			Post(api.basePath + resource)
		if err != nil || len((*resp.Result().(*AssetsInfo)).Error) > 0 {
			fmt.Println("Kraken API assets : Error Requesting AssetsInfo" + strings.Join((*resp.Result().(*AssetsInfo)).Error, ""))
		}
		result := (*resp.Result().(*AssetsInfo)).Result.(map[string]interface{})
		if useCache {
			err = db.Write("Kraken/public", "Assets", result)
			if err != nil {
				fmt.Println("Kraken API assets : Error Caching AssetsInfo")
			}
		}
	}
}
