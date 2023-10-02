package lookup

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/proxy"
)

func LookupAPI(proxy proxy.Dialer, number string) error {
	httpClient := &http.Client{Transport: &http.Transport{Dial: proxy.Dial}}
	lookupURL := fmt.Sprintf("https://api.telnyx.com/v2/number_lookup/%s", number)
	req, err := http.NewRequest("GET", lookupURL, nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Accept":        {"application/json"},
		"Authorization": {"Bearer KEY018AF21010C0D85BE32E33A29235B9EC_oUoYAaDqQaIiO6tIHn7GQm"},
	}

	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	defer res.Body.Close()
	fmt.Printf("res.StatusCode: %v\n", res.StatusCode)
	fmt.Printf("body: %v\n", string(body))
	return nil
}
