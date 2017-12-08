package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/heroku/terraform-provider-librato/librato"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: librato.Provider})
}
