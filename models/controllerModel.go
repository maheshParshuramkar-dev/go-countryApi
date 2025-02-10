package models

type CountryData struct {
	Name       name                    `json:"name"`
	Currencies map[string]intCurruncey `json:"currencies"`
	Capital    []string                `json:"capital"`
	Population int64                   `json:"population"`
}

type name struct {
	Common string `json:"common"`
}

type intCurruncey struct {
	Symbol string `json:"symbol"`
}

type ApiRes struct {
	Status bool   `json:"status"`
	Result Result `json:"result"`
}

type Result struct {
	Name       string `json:"name,omitempty"`
	Capital    string `json:"capital,omitempty"`
	Currency   string `json:"currency,omitempty"`
	Population int64  `json:"population,omitempty"`
	Error      string `json:"error,omitempty"`
}
