package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Activation struct {
	Network string `json:"network"`
	Dryrun  bool   `json:"dryrun"`
	Version int    `json:"version"`
}

type ActivationResponse struct {
	Network string `json:"network"`
	Dryrun  bool   `json:"dryrun"`
	Version int    `json:"version"`
	Status  string `json:"status"`
}

func resourceALBActivation() *schema.Resource {
	return &schema.Resource{
		Create: resourceALBActivationCreate,
		Read:   resourceALBActivationRead,
		Update: resourceALBActivationCreate,
		Delete: resourceALBActivationDelete,

		Schema: map[string]*schema.Schema{
			"origin_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceALBActivationDelete(d *schema.ResourceData, meta interface{}) error {
	return schema.Noop(d, meta)
}

func resourceALBActivationCreate(d *schema.ResourceData, meta interface{}) error {

	network := d.Get("network").(string)
	version := d.Get("version").(int)
	originId := d.Get("origin_id").(string)
	config := meta.(*edgegrid.Config)

	var activation Activation
	activation.Network = network
	activation.Dryrun = false
	activation.Version = version

	a, _ := json.Marshal(activation)

	req, _ := client.NewRequest(
		*config,
		"POST",
		fmt.Sprintf("/cloudlets/api/v2/origins/%s/activations", originId),
		bytes.NewBuffer(a),
	)
	resp, err := client.Do(*config, req)
	if err != nil {
		log.Fatal(err)
	}

	byt, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal(string(byt))
	}

	var activations Activation
	err = json.Unmarshal(byt, &activations)

	d.SetId(fmt.Sprintf("%s:%d", originId, version))

	// Iterate until it's deployed
	for i := 1; i < 10; i++ {
		var activations = getActivations(network, version, originId, config)
		err = json.Unmarshal(byt, &activations)
		if isActive(activations, network, version) {
			break
		}
		time.Sleep(20 * time.Second)
	}

	return resourceALBActivationRead(d, meta)
}

func isActive(activations []ActivationResponse, network string, version int) bool {
	for _, activation := range activations {
		if strings.ToLower(activation.Network) == strings.ToLower(network) && activation.Status == "active" && activation.Version == version {
			return true
		}
	}
	return false
}

func resourceALBActivationRead(d *schema.ResourceData, meta interface{}) error {

	network := d.Get("network").(string)
	version := d.Get("version").(int)
	originId := d.Get("origin_id").(string)
	config := meta.(*edgegrid.Config)

	var activations = getActivations(network, version, originId, config)

	for _, activation := range activations {
		// TODO: if we find an activation "pending", wait for it to become "active" before returning
		if strings.ToLower(activation.Network) == strings.ToLower(network) && activation.Status == "active" {
			d.SetId(fmt.Sprintf("%s:%d", originId, version))
			d.Set("version", version)
		}
	}

	return nil
}

func getActivations(network string, version int, originId string, config *edgegrid.Config) []ActivationResponse {
	req, _ := client.NewRequest(
		*config,
		"GET",
		fmt.Sprintf("/cloudlets/api/v2/origins/%s/activations", originId),
		nil,
	)
	resp, err := client.Do(*config, req)
	if err != nil {
		log.Fatal(err)
	}

	byt, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal(string(byt))
	}

	var activations []ActivationResponse
	err = json.Unmarshal(byt, &activations)
	return activations
}
