/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_POOL_RESOURCE3 = `
resource "bigip_ltm_node" "test-node" {
        name = "` + TEST_NODE_NAME + `"
        address = "10.10.10.10"
        connection_limit = "0"
        dynamic_ratio = "1"
        monitor = "default"
        rate_limit = "disabled"
}

resource "bigip_ltm_pool" "test-pool" {
        name = "` + TEST_POOL_NAME + `"
        monitors = ["/Common/http"]
        allow_nat = "yes"
        allow_snat = "yes"
        description = "Test-Pool-Sample"
        load_balancing_mode = "round-robin"
        slow_ramp_time = "5"
        service_down_action = "reset"
        reselect_tries = "2"
}
`

var TEST_POOL_RESOURCE1 = `
resource "bigip_ltm_node" "test-node" {
	name = "` + TEST_NODE_NAME + `"
	address = "10.10.10.10"
	connection_limit = "0"
	dynamic_ratio = "1"
	monitor = "default"
	rate_limit = "disabled"
        fqdn {
    address_family = "ipv4"
    interval       = "3000"
  }
}

resource "bigip_ltm_pool" "test-pool" {
	name = "` + TEST_POOL_NAME + `"
	monitors = ["/Common/http"]
	allow_nat = "yes"
	allow_snat = "yes"
	description = "Test-Pool-Sample"
	load_balancing_mode = "round-robin"
	slow_ramp_time = "5"
	service_down_action = "reset"
	reselect_tries = "2"
}

resource "bigip_ltm_pool_attachment" "test-pool_test-node" {
	pool = bigip_ltm_pool.test-pool.name
	node = "${bigip_ltm_node.test-node.name}:443"
}
`
var TEST_POOL_RESOURCE2 = `
resource "bigip_ltm_node" "test-node" {
        name = "` + TEST_NODE_NAME + `"
        address = "10.10.10.10"
        connection_limit = "0"
        dynamic_ratio = "1"
        monitor = "default"
        rate_limit = "disabled"
        fqdn {
    address_family = "ipv4"
    interval       = "3000"
  }
}

resource "bigip_ltm_pool" "test-pool" {
        name = "` + TEST_POOL_NAME + `"
        monitors = ["/Common/http"]
        allow_nat = "yes"
        allow_snat = "yes"
        description = "Test-Pool-Sample"
        load_balancing_mode = "round-robin"
        slow_ramp_time = "5"
        service_down_action = "reset"
        reselect_tries = "2"
}
`

func TestAccBigipLtmPoolAttachment_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_POOL_RESOURCE1,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TEST_POOL_NAME, true),
					testCheckPoolAttachment(TEST_POOL_NAME, TEST_POOLNODE_NAMEPORT, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "pool", TEST_POOL_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "node", TEST_POOLNODE_NAMEPORT),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_POOL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TEST_POOL_NAME, true),
					testCheckPoolAttachment(TEST_POOL_NAME, TEST_POOLNODE_NAMEPORT, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "pool", TEST_POOL_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "node", TEST_POOLNODE_NAMEPORT),
				),
			},
		},
	})
}

func TestAccBigipLtmPoolAttachment_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_POOL_RESOURCE1,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TEST_POOL_NAME, true),
					testCheckPoolAttachment(TEST_POOL_NAME, TEST_POOLNODE_NAMEPORT, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "pool", TEST_POOL_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "node", TEST_POOLNODE_NAMEPORT),
				),
			},
			{
				Config: TEST_POOL_RESOURCE2,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TEST_POOL_NAME, true),
					testCheckPoolAttachment(TEST_POOL_NAME, TEST_POOLNODE_NAMEPORT, false),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_POOL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TEST_POOL_NAME, true),
					testCheckPoolAttachment(TEST_POOL_NAME, TEST_POOLNODE_NAMEPORT, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "pool", TEST_POOL_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "node", TEST_POOLNODE_NAMEPORT),
				),
			},
			{
				Config: TEST_POOL_RESOURCE3,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TEST_POOL_NAME, true),
					testCheckPoolAttachment(TEST_POOL_NAME, TEST_POOLNODE_NAMEPORT, false),
				),
			},
		},
	})
}
func testCheckPoolAttachment(poolName string, expected string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		pool, err := client.GetPool(poolName)
		if err != nil {
			return err
		}
		if pool == nil {
			return fmt.Errorf("Pool %s does not exist.", poolName)
		}

		nodes, err := client.PoolMembers(poolName)
		if err != nil {
			return fmt.Errorf("Error retrieving pool (%s) members: %s", poolName, err)
		}
		if nodes == nil {
			return fmt.Errorf("Pool member %s does not exist.", expected)
		}
		found := false
		for _, node := range nodes.PoolMembers {
			if expected == node.FullPath {
				found = true
				break
			}
		}

		if !found && exists {
			return fmt.Errorf("Node %s is not a member of pool %s", expected, poolName)
		}
		if found && !exists {
			return fmt.Errorf("Node %s is still  a member of pool %s", expected, poolName)
		}

		return nil
	}
}
