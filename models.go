package main

type Domain struct {
	Name                  string `json:"name"`
	ContainerName         string `json:"container_name"`
	IP                    string `json:"ip"`
	HttpEnabled           bool   `json:"http_enabled"`
	HttpsEnabled          bool   `json:"https_enabled"`
	HttpWebsocketsEnabled bool   `json:"http_websockets_enabled"`
}

type Config struct {
	Domains []Domain `json:"domains"`
}
