package softlayer

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccSoftLayerLbLocal_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSoftLayerLbLocalConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"softlayer_lb_local.testacc_foobar_lb", "connections", "15000"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_local.testacc_foobar_lb", "datacenter", "tok02"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_local.testacc_foobar_lb", "ha_enabled", "false"),
				),
			},
		},
	})
}

func TestAccSoftLayerLbLocalDedicated_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSoftLayerLbLocalDedicatedConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"softlayer_lb_local.testacc_foobar_lb", "connections", "15000"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_local.testacc_foobar_lb", "datacenter", "tok02"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_local.testacc_foobar_lb", "ha_enabled", "false"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_local.testacc_foobar_lb", "dedicated", "true"),
				),
			},
		},
	})
}

const testAccCheckSoftLayerLbLocalConfig_basic = `
resource "softlayer_lb_local" "testacc_foobar_lb" {
    connections = 15000
    datacenter    = "tok02"
    ha_enabled  = false
}`

const testAccCheckSoftLayerLbLocalDedicatedConfig_basic = `
resource "softlayer_lb_local" "testacc_foobar_lb" {
    connections = 15000
    datacenter    = "tok02"
    ha_enabled  = false
    dedicated = true	
}`
