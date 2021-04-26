package main

import (
	"github.com/IanCassTwo/terraform-provider-iancass/iancass"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: iancass.Provider,
	})
}
