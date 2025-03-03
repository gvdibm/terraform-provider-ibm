// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
)

func TestAccIBMPIDhcpbasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMPIDhcpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPIDhcpConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPIDhcpExists("ibm_pi_dhcp.dhcp_service"),
					resource.TestCheckResourceAttrSet(
						"ibm_pi_dhcp.dhcp_service", "dhcp_id"),
				),
			},
		},
	})
}

func testAccCheckIBMPIDhcpDestroy(s *terraform.State) error {
	sess, err := testAccProvider.Meta().(ClientSession).IBMPISession()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_pi_dhcp" {
			continue
		}

		cloudInstanceID, dhcpID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := st.NewIBMPIDhcpClient(context.Background(), sess, cloudInstanceID)
		_, err = client.Get(dhcpID)
		if err == nil {
			return fmt.Errorf("PI DHCP still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
func testAccCheckIBMPIDhcpExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}

		sess, err := testAccProvider.Meta().(ClientSession).IBMPISession()
		if err != nil {
			return err
		}

		cloudInstanceID, dhcpID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := st.NewIBMPIDhcpClient(context.Background(), sess, cloudInstanceID)

		_, err = client.Get(dhcpID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckIBMPIDhcpConfig() string {
	return fmt.Sprintf(`
	resource "ibm_pi_dhcp" "dhcp_service" {
		pi_cloud_instance_id = "%s"
	}
	`, pi_cloud_instance_id)
}
