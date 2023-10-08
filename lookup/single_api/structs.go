package single_api

// HRL Lookup Response
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

// for Account Balance
type BalanceInfo struct {
	Balance         string `json:"balance"`
	CreditLimit     string `json:"credit_limit"`
	AvailableCredit string `json:"available_credit"`
	Currency        string `json:"currency"`
	RecordType      string `json:"record_type"`
}

type DataInfo struct {
	Data BalanceInfo `json:"data"`
}
