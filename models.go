package main

type Domain struct {
	Name                  string `json:"name"`
	ContainerName         string `json:"container_name"`
	Ip                    string `json:"ip"`
	HttpEnabled           bool   `json:"http_enabled"`
	HttpsEnabled          bool   `json:"https_enabled"`
	HttpWebsocketsEnabled bool   `json:"http_websockets_enabled"`
}

type Config struct {
	Domains []Domain `json:"domains"`
}

type Container struct {
	ContainerName string
	Ip            string
}

type TemplateParams struct {
	Domains    []Domain
	Containers []Container
}
