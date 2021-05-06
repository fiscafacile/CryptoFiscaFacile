package kraken

import (
	"fmt"
	"strings"

	scribble "github.com/nanobox-io/golang-scribble"
)

func (api *api) getAPIAssets() {
	const SOURCE = "Kraken API Assets :"
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
		_, err := api.clientAssets.R().
			SetHeaders(headers).
			SetResult(&api.assets).
			Post(api.basePath + resource)
		if err != nil || len(api.assets.Error) > 0 {
			fmt.Println(SOURCE, "Error Requesting AssetsInfo"+strings.Join(api.assets.Error, ""))
		}
		if useCache {
			err = db.Write("Kraken/public", "Assets", api.assets.Result)
			if err != nil {
				fmt.Println(SOURCE, "Error Caching AssetsInfo")
			}
		}
	}
}
