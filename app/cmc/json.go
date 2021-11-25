package cmc

type Price struct {
	Price float64 `json:"price"`
}

type CMCQuote struct {
	EUR Price `json:"EUR"`
}

type CMCCryptocurrency struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Symbol string   `json:"symbol"`
	Quote  CMCQuote `json:"quote"`
	Logo   string   `json:"logo"`
}

type CMCResponse struct {
	Data []CMCCryptocurrency `json:"data"`
}

type CMCMapResponse struct {
	Data map[string]CMCCryptocurrency `json:"data"`
}
