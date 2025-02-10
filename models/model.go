package models

type Configurations struct {
	Env     string
	Prefix  string
	Server  ServerConfig
	ExtUrls ExtUrls
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

type ExtUrls struct {
	CountiresFetchUrl string
	CtxTimeoutUrl     int
}
