package provider

import (
        "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/IanCassTwo/terraform-provider-iancass/api/siteshield"
	"log"
)

func dataSourceAkamaiSiteShield() *schema.Resource {
	return &schema.Resource{
		Read: dataAkamaiSiteShieldRead,
		Schema: map[string]*schema.Schema{
			"mapid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"acknowledged": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"currentcidrs": {
                                Type:         schema.TypeSet,
                                Elem: &schema.Schema {
                                        Type: schema.TypeString,
                                },
                                Computed:     true,
                        },
			"proposedcidrs": {
                                Type:         schema.TypeSet,
                                Elem: &schema.Schema {
                                        Type: schema.TypeString,
                                },
                                Computed:     true,
                        },
			"contacts": {
                                Type:         schema.TypeSet,
                                Elem: &schema.Schema {
                                        Type: schema.TypeString,
                                },
                                Computed:     true,
                        },
		},
	}
}

func dataAkamaiSiteShieldRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] entered dataAkamaiSiteShieldRead")


	id := d.Get("mapid").(string)
	siteshieldmapresponse, err := siteshield.GetMap(id)
	if err != nil {
		return err
	}

	acknowledged := siteshieldmapresponse.Acknowledged
	d.Set("acknowledged", acknowledged)

	currentcidrs := siteshieldmapresponse.CurrentCidrs
	d.Set("currentcidrs", currentcidrs)

	proposedcidrs := siteshieldmapresponse.ProposedCidrs
	d.Set("proposedcidrs", proposedcidrs)

	contacts := siteshieldmapresponse.Contacts
	d.Set("contacts", contacts)
	
	d.SetId(id)

	return nil
}
