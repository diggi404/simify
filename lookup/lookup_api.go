package lookup

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/proxy"
)

type Portability struct {
	LRN                string `json:"lrn"`
	PortedStatus       string `json:"ported_status"`
	PortedDate         string `json:"ported_date"`
	OCN                string `json:"ocn"`
	LineType           string `json:"line_type"`
	SPID               string `json:"spid"`
	SPIDCarrierName    string `json:"spid_carrier_name"`
	SPIDCarrierType    string `json:"spid_carrier_type"`
	AltSPID            string `json:"altspid"`
	AltSPIDCarrierName string `json:"altspid_carrier_name"`
	AltSPIDCarrierType string `json:"altspid_carrier_type"`
	City               string `json:"city"`
	State              string `json:"state"`
}

type Data struct {
	CountryCode    string      `json:"country_code"`
	NationalFormat string      `json:"national_format"`
	PhoneNumber    string      `json:"phone_number"`
	Fraud          interface{} `json:"fraud"`
	Carrier        interface{} `json:"carrier"`
	CallerName     interface{} `json:"caller_name"`
	NNIDOverride   interface{} `json:"nnid_override"`
	Portability    Portability `json:"portability"`
	ValidNumber    bool        `json:"valid_number"`
	RecordType     string      `json:"record_type"`
}

type PhoneNumberInfo struct {
	Data Data `json:"data"`
}

func LookupAPI(proxy proxy.Dialer, number string) error {
	httpClient := &http.Client{Transport: &http.Transport{Dial: proxy.Dial}}
	// lookupURL := fmt.Sprintf("https://api.telnyx.com/v2/number_lookup/%s", number)
	req, err := http.NewRequest("GET", "https://api.telnyx.com/v2/account", nil)
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
