package get_geo_by_ip

import "time"

type Response struct {
	Ip            string  `json:"ip"`
	Success       bool    `json:"success"`
	Type          string  `json:"type"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continent_code"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	City          string  `json:"city"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	IsEu          bool    `json:"is_eu"`
	Postal        string  `json:"postal"`
	CallingCode   string  `json:"calling_code"`
	Capital       string  `json:"capital"`
	Borders       string  `json:"borders"`
	Flag          struct {
		Img          string `json:"img"`
		Emoji        string `json:"emoji"`
		EmojiUnicode string `json:"emoji_unicode"`
	} `json:"flag"`
	Connection struct {
		Asn    int    `json:"asn"`
		Org    string `json:"org"`
		Isp    string `json:"isp"`
		Domain string `json:"domain"`
	} `json:"connection"`
	Timezone struct {
		Id          string    `json:"id"`
		Abbr        string    `json:"abbr"`
		IsDst       bool      `json:"is_dst"`
		Offset      int       `json:"offset"`
		Utc         string    `json:"utc"`
		CurrentTime time.Time `json:"current_time"`
	} `json:"timezone"`
	Currency struct {
		Name         string `json:"name"`
		Code         string `json:"code"`
		Symbol       string `json:"symbol"`
		Plural       string `json:"plural"`
		ExchangeRate int    `json:"exchange_rate"`
	} `json:"currency"`
	Security struct {
		Anonymous bool `json:"anonymous"`
		Proxy     bool `json:"proxy"`
		Vpn       bool `json:"vpn"`
		Tor       bool `json:"tor"`
		Hosting   bool `json:"hosting"`
	} `json:"security"`
	Rate struct {
		Limit     int `json:"limit"`
		Remaining int `json:"remaining"`
	} `json:"rate"`
}
