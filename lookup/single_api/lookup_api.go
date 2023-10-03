package single_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/proxy"
)

func LookupAPI(proxy proxy.Dialer, number string, apiKey string) (PhoneNumberInfo, int, error) {
	httpClient := &http.Client{Transport: &http.Transport{Dial: proxy.Dial}}
	lookupURL := fmt.Sprintf("https://api.telnyx.com/v2/number_lookup/%s", number)
	req, err := http.NewRequest("GET", lookupURL, nil)
	if err != nil {
		return PhoneNumberInfo{}, 0, err
	}
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Accept":        {"application/json"},
		"Authorization": {"Bearer " + apiKey},
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return PhoneNumberInfo{}, 0, err
	}
	if res.StatusCode == 200 {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return PhoneNumberInfo{}, res.StatusCode, err
		}
		var lookupInfo PhoneNumberInfo
		err = json.Unmarshal(b, &lookupInfo)
		if err != nil {
			return PhoneNumberInfo{}, res.StatusCode, err
		}
		return lookupInfo, res.StatusCode, nil
	}
	err = errors.New("invalid number or api key")
	return PhoneNumberInfo{}, res.StatusCode, err
}
