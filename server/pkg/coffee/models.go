package coffee

type SupporterList struct {
	Data  []Supporter `json:"data"`
	Total int32       `json:"total"`
}

type Supporter struct {
	Qty      int32  `json:"support_coffees"`
	Price    string `json:"support_coffee_price"`
	Currency string `json:"support_currency"`
	Name     string `json:"payer_name"`
	Note     string `json:"support_note"`
}
