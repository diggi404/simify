package single_api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func CheckBalance(apiKey string) (DataInfo, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.telnyx.com/v2/balance", nil)
	if err != nil {
		return DataInfo{}, err
	}
	req.Header = http.Header{
		"Authorization": {"Bearer " + apiKey},
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return DataInfo{}, err
	}
	if res.StatusCode == 200 {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return DataInfo{}, err
		}
		var balanceInfo DataInfo
		json.Unmarshal(b, &balanceInfo)
		return balanceInfo, nil
	}
	err = errors.New("your api key is invalid")
	return DataInfo{}, err

}
