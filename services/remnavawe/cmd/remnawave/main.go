package main

import (
	"remnawave/config"
	"remnawave/server/api/rest"
)

func main() {
	api := rest.NewAPIServer(config.NewConfig())
	api.Start()
}
