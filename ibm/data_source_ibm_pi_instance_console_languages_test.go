// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMPIInstanceConsoleLanguages(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPIInstanceConsoleLanguagesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_pi_console_languages.example", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_pi_console_languages.example", "console_languages.#"),
				),
			},
		},
	})
}

func testAccCheckIBMPIInstanceConsoleLanguagesConfig() string {
	return fmt.Sprintf(`
	data "ibm_pi_console_languages" "example" {
		pi_cloud_instance_id = "%s"
		pi_instance_name     = "%s"
	}`, pi_cloud_instance_id, pi_instance_name)
}
