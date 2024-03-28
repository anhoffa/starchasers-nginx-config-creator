package main

type Domain struct {
	Name                  string `json:"name"`
	Owner                 string `json:"owner"`
	IP                    string `json:"ip"`
	HTTPEnabled           bool   `json:"http_enabled"`
	HTTPSEnabled          bool   `json:"https_enabled"`
	HTTPWebsocketsEnabled bool   `json:"http_websockets_enabled"`
}

type Config struct {
	Domains []Domain `json:"domains"`
}
