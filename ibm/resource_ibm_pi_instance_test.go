// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
)

func testAccCheckIBMPIInstanceConfig(name, instanceHealthStatus string) string {
	return fmt.Sprintf(`
	resource "ibm_pi_key" "key" {
		pi_cloud_instance_id = "%[1]s"
		pi_key_name          = "%[2]s"
		pi_ssh_key           = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR"
	  }
	  data "ibm_pi_image" "power_image" {
		pi_image_name        = "%[3]s"
		pi_cloud_instance_id = "%[1]s"
	  }
	  data "ibm_pi_network" "power_networks" {
		pi_cloud_instance_id = "%[1]s"
		pi_network_name      = "%[4]s"
	  }
	  resource "ibm_pi_volume" "power_volume" {
		pi_volume_size       = 20
		pi_volume_name       = "%[2]s"
		pi_volume_shareable  = true
		pi_volume_pool       = data.ibm_pi_image.power_image.storage_pool
		pi_cloud_instance_id = "%[1]s"
	  }
	  resource "ibm_pi_instance" "power_instance" {
		pi_memory             = "2"
		pi_processors         = "0.25"
		pi_instance_name      = "%[2]s"
		pi_proc_type          = "shared"
		pi_image_id           = data.ibm_pi_image.power_image.id
		pi_key_pair_name      = ibm_pi_key.key.key_id
		pi_sys_type           = "s922"
		pi_cloud_instance_id  = "%[1]s"
		pi_storage_pool       = data.ibm_pi_image.power_image.storage_pool
		pi_health_status      = "%[5]s"
		pi_volume_ids         = [ibm_pi_volume.power_volume.volume_id]
		pi_network {
			network_id = data.ibm_pi_network.power_networks.id
		}
	  }
	`, pi_cloud_instance_id, name, pi_image, pi_network_name, instanceHealthStatus)
}

func testAccIBMPIInstanceNetworkConfig(name, privateNetIP string) string {
	return fmt.Sprintf(`
	resource "ibm_pi_key" "key" {
		pi_cloud_instance_id = "%[1]s"
		pi_key_name          = "%[2]s"
		pi_ssh_key           = "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEArb2aK0mekAdbYdY9rwcmeNSxqVCwez3WZTYEq+1Nwju0x5/vQFPSD2Kp9LpKBbxx3OVLN4VffgGUJznz9DAr7veLkWaf3iwEil6U4rdrhBo32TuDtoBwiczkZ9gn1uJzfIaCJAJdnO80Kv9k0smbQFq5CSb9H+F5VGyFue/iVd5/b30MLYFAz6Jg1GGWgw8yzA4Gq+nO7HtyuA2FnvXdNA3yK/NmrTiPCdJAtEPZkGu9LcelkQ8y90ArlKfjtfzGzYDE4WhOufFxyWxciUePh425J2eZvElnXSdGha+FCfYjQcvqpCVoBAG70U4fJBGjB+HL/GpCXLyiYXPrSnzC9w=="
	}
	resource "ibm_pi_instance" "power_instance" {
		pi_memory             = "2"
		pi_processors         = "0.25"
		pi_instance_name      = "%[2]s"
		pi_proc_type          = "shared"
		pi_image_id           = "f4501cad-d0f4-4517-9eea-85402309d90d"
		pi_key_pair_name      = ibm_pi_key.key.key_id
		pi_sys_type           = "e980"
		pi_storage_type 	  = "tier3"
		pi_cloud_instance_id  = "%[1]s"
		pi_network {
			network_id = "tf-cloudconnection-23"
			ip_address = "%[3]s"
		}
	}
	`, pi_cloud_instance_id, name, privateNetIP)
}

func testAccIBMPIInstanceVTLConfig(name string) string {
	return fmt.Sprintf(`
	resource "ibm_pi_key" "vtl_key" {
		pi_cloud_instance_id = "%[1]s"
		pi_key_name          = "%[2]s"
		pi_ssh_key           = "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEArb2aK0mekAdbYdY9rwcmeNSxqVCwez3WZTYEq+1Nwju0x5/vQFPSD2Kp9LpKBbxx3OVLN4VffgGUJznz9DAr7veLkWaf3iwEil6U4rdrhBo32TuDtoBwiczkZ9gn1uJzfIaCJAJdnO80Kv9k0smbQFq5CSb9H+F5VGyFue/iVd5/b30MLYFAz6Jg1GGWgw8yzA4Gq+nO7HtyuA2FnvXdNA3yK/NmrTiPCdJAtEPZkGu9LcelkQ8y90ArlKfjtfzGzYDE4WhOufFxyWxciUePh425J2eZvElnXSdGha+FCfYjQcvqpCVoBAG70U4fJBGjB+HL/GpCXLyiYXPrSnzC9w=="
	}
	
	resource "ibm_pi_network" "vtl_network" {
		pi_cloud_instance_id = "%[1]s"
		pi_network_name      = "%[2]s"
		pi_network_type      = "pub-vlan"
	}

	resource "ibm_pi_instance" "vtl_instance" {
		pi_memory             = "22"
		pi_processors         = "2"
		pi_instance_name      = "%[2]s"
		pi_license_repository_capacity = "3"
		pi_proc_type          = "shared"
		pi_image_id           = "ca4ea55f-b329-4cf5-bdce-d2f38cfc6da3"
		pi_key_pair_name      = ibm_pi_key.vtl_key.key_id
		pi_sys_type           = "s922"
		pi_cloud_instance_id  = "%[1]s"
		pi_storage_type 	  = "tier1"
		pi_network {
			network_id = ibm_pi_network.vtl_network.network_id
		}
	  }
	
	`, pi_cloud_instance_id, name)
}

func testAccCheckIBMPIInstanceDestroy(s *terraform.State) error {
	sess, err := testAccProvider.Meta().(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_pi_instance" {
			continue
		}
		cloudInstanceID, instanceID, err := splitID(rs.Primary.ID)
		if err == nil {
			return err
		}
		client := st.NewIBMPIInstanceClient(context.Background(), sess, cloudInstanceID)
		_, err = client.Get(instanceID)
		if err == nil {
			return fmt.Errorf("PI Instance still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
func testAccCheckIBMPIInstanceExists(n string) resource.TestCheckFunc {
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

		cloudInstanceID, instanceID, err := splitID(rs.Primary.ID)
		if err == nil {
			return err
		}
		client := st.NewIBMPIInstanceClient(context.Background(), sess, cloudInstanceID)

		_, err = client.Get(instanceID)
		if err != nil {
			return err
		}

		return nil
	}
}

func TestAccIBMPIInstanceBasic(t *testing.T) {
	instanceRes := "ibm_pi_instance.power_instance"
	name := fmt.Sprintf("tf-pi-instance-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMPIInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPIInstanceConfig(name, helpers.PIInstanceHealthWarning),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPIInstanceExists(instanceRes),
					resource.TestCheckResourceAttr(instanceRes, "pi_instance_name", name),
				),
			},
		},
	})
}

func TestAccIBMPIInstanceNetwork(t *testing.T) {
	instanceRes := "ibm_pi_instance.power_instance"
	name := fmt.Sprintf("tf-pi-instance-%d", acctest.RandIntRange(10, 100))
	privateNetIP := "192.112.111.220"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMPIInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIBMPIInstanceNetworkConfig(name, privateNetIP),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPIInstanceExists(instanceRes),
					resource.TestCheckResourceAttr(instanceRes, "pi_instance_name", name),
					resource.TestCheckResourceAttrSet(instanceRes, "pi_network.0.network_id"),
					resource.TestCheckResourceAttrSet(instanceRes, "pi_network.0.mac_address"),
					resource.TestCheckResourceAttr(instanceRes, "pi_network.0.ip_address", privateNetIP),
				),
			},
		},
	})
}

func TestAccIBMPIInstanceVTL(t *testing.T) {
	instanceRes := "ibm_pi_instance.vtl_instance"
	name := fmt.Sprintf("tf-pi-vtl-instance-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMPIInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIBMPIInstanceVTLConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPIInstanceExists(instanceRes),
					resource.TestCheckResourceAttr(instanceRes, "pi_instance_name", name),
					resource.TestCheckResourceAttr(instanceRes, "pi_license_repository_capacity", "3"),
				),
			},
		},
	})
}

func TestAccIBMPISAPInstance(t *testing.T) {
	instanceRes := "ibm_pi_instance.sap"
	name := fmt.Sprintf("tf-pi-sap-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMPIInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIBMPISAPInstanceConfig(name, "tinytest-1x4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPIInstanceExists(instanceRes),
					resource.TestCheckResourceAttr(instanceRes, "pi_instance_name", name),
					resource.TestCheckResourceAttr(instanceRes, "pi_sap_profile_id", "tinytest-1x4"),
				),
			},
			{
				Config: testAccIBMPISAPInstanceConfig(name, "tinytest-1x8"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPIInstanceExists(instanceRes),
					resource.TestCheckResourceAttr(instanceRes, "pi_instance_name", name),
					resource.TestCheckResourceAttr(instanceRes, "pi_sap_profile_id", "tinytest-1x8"),
				),
			},
		},
	})
}
func testAccIBMPISAPInstanceConfig(name, sapProfile string) string {
	return fmt.Sprintf(`
	resource "ibm_pi_network" "power_network" {
		pi_cloud_instance_id	= "%[1]s"
		pi_network_name			= "%[2]s"
		pi_network_type			= "pub-vlan"
	}
	resource "ibm_pi_instance" "sap" {
		pi_cloud_instance_id  	= "%[1]s"
		pi_instance_name      	= "%[2]s"
		pi_sap_profile_id       = "%[3]s"
		pi_image_id           	= "ef9a2f2e-6b36-48cb-aa06-223040ddb9d2"
		pi_storage_type			= "tier1"
		pi_network {
			network_id = ibm_pi_network.power_network.network_id
		}
		pi_health_status		= "OK"
	}
	`, pi_cloud_instance_id, name, sapProfile)
}

func TestAccIBMPIInstanceMixedStorage(t *testing.T) {
	instanceRes := "ibm_pi_instance.instance"
	name := fmt.Sprintf("tf-pi-mixedstorage-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMPIInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIBMPIInstanceMixedStorage(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPIInstanceExists(instanceRes),
					resource.TestCheckResourceAttr(instanceRes, "pi_instance_name", name),
					resource.TestCheckResourceAttr(instanceRes, "pi_storage_pool_affinity", "false"),
				),
			},
		},
	})
}

func testAccIBMPIInstanceMixedStorage(name string) string {
	return fmt.Sprintf(`
	resource "ibm_pi_key" "key" {
		pi_cloud_instance_id = "%[1]s"
		pi_key_name          = "%[2]s"
		pi_ssh_key           = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR"
	}
	resource "ibm_pi_network" "power_network" {
		pi_cloud_instance_id = "%[1]s"
		pi_network_name      = "%[2]s"
		pi_network_type      = "pub-vlan"
	}
	resource "ibm_pi_volume" "power_volume" {
		pi_cloud_instance_id = "%[1]s"
		pi_volume_size       = 20
		pi_volume_name       = "%[2]s"
		pi_volume_shareable  = true
		pi_volume_type       = "tier3"
	}
	resource "ibm_pi_instance" "instance" {
		pi_cloud_instance_id     = "%[1]s"
		pi_memory                = "2"
		pi_processors            = "0.25"
		pi_instance_name         = "%[2]s"
		pi_proc_type             = "shared"
		pi_image_id              = "ca4ea55f-b329-4cf5-bdce-d2f38cfc6da3"
		pi_key_pair_name         = ibm_pi_key.key.key_id
		pi_sys_type              = "s922"
		pi_storage_type          = "tier1"
		pi_storage_pool_affinity = false
		pi_network {
			network_id = ibm_pi_network.power_network.network_id
		}
	}
	resource "ibm_pi_volume_attach" "power_attach_volume"{
		pi_cloud_instance_id = "%[1]s"
		pi_volume_id         = ibm_pi_volume.power_volume.volume_id
		pi_instance_id       = ibm_pi_instance.instance.instance_id
	}
	`, pi_cloud_instance_id, name)
}
