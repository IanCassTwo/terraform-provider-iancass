package iancass

import (
	"errors"
	"fmt"
	//	"log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"strings"
)

const (
	Version = "0.0.1"
)

type Config struct {
}

func getConfigOptions(section string) *schema.Resource {
	section = strings.ToUpper(section)

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_HOST"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_HOST"); v != "" {
						return v, nil
					}

					return nil, errors.New("host not set")
				},
			},
			"access_token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_ACCESS_TOKEN"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_ACCESS_TOKEN"); v != "" {
						return v, nil
					}

					return nil, errors.New("access_token not set")
				},
			},
			"client_token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_CLIENT_TOKEN"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_CLIENT_TOKEN"); v != "" {
						return v, nil
					}

					return nil, errors.New("client_token not set")
				},
			},
			"client_secret": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_CLIENT_SECRET"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_CLIENT_SECRET"); v != "" {
						return v, nil
					}

					return nil, errors.New("client_secret not set")
				},
			},
			"max_body": {
				Type:     schema.TypeInt,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_MAX_SIZE"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_MAX_SIZE"); v != "" {
						return v, nil
					}

					return 131072, nil
				},
			},
		},
	}
}

// Provider returns the Akamai terraform.Resource provider.
func Provider() *schema.Provider {
	client.UserAgent = client.UserAgent + " terraform-iancass/" + Version

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"edgerc": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
			},
			"section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
			},
			"auth": &schema.Schema{
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     getConfigOptions("property"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"iancass_alb_activation": resourceALBActivation(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config, err := getConfig(d)
	if err != nil {
		return nil, fmt.Errorf("at least one configuration must be defined")
	}

	return config, nil
}

type resourceData interface {
	GetOk(string) (interface{}, bool)
	Get(string) interface{}
}

type set interface {
	List() []interface{}
}

func getConfig(d resourceData) (*edgegrid.Config, error) {
	var config edgegrid.Config
	var err error
	if _, ok := d.GetOk("auth"); ok {
		attributes := d.Get("auth").(set).List()[0].(map[string]interface{})

		config = edgegrid.Config{
			Host:         attributes["host"].(string),
			AccessToken:  attributes["access_token"].(string),
			ClientToken:  attributes["client_token"].(string),
			ClientSecret: attributes["client_secret"].(string),
			MaxBody:      attributes["max_body"].(int),
		}

		return &config, nil
	}

	// default
	edgerc := d.Get("edgerc").(string)
	section := d.Get("section").(string)
	config, err = edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
