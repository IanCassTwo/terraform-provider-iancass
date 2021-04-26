package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/IanCassTwo/terraform-provider-iancass/iancass"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: iancass.Provider,
	})
}
