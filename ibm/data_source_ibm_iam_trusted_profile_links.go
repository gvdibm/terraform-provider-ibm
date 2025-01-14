// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func dataSourceIBMIamTrustedProfileLinks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileLinkListRead,

		Schema: map[string]*schema.Schema{
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile.",
			},
			"links": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of links to a trusted profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the unique identifier of the claim rule.",
						},
						"entity_tag": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "version of the claim rule.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If set contains a date time string of the creation date in ISO format.",
						},
						"modified_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If set contains a date time string of the last modification date in ISO format.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Optional name of the Link.",
						},
						"cr_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The compute resource type. Valid values are VSI, IKS_SA, ROKS_SA.",
						},
						"link": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN of the compute resource.",
									},
									"namespace": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The compute resource namespace, only required if cr_type is IKS_SA or ROKS_SA.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the compute resource, only required if cr_type is IKS_SA or ROKS_SA.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIamTrustedProfileLinkListRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listLinkOptions := &iamidentityv1.ListLinksOptions{}

	listLinkOptions.SetProfileID(d.Get("profile_id").(string))

	profileLinkList, response, err := iamIdentityClient.ListLinks(listLinkOptions)
	if err != nil {
		log.Printf("[DEBUG] ListLink failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListLink failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMIamTrustedProfileLinkListID(d))

	if profileLinkList.Links != nil {
		err = d.Set("links", dataSourceProfileLinkListFlattenLinks(profileLinkList.Links))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting links %s", err))
		}
	}

	return nil
}

// dataSourceIBMIamTrustedProfileLinkListID returns a reasonable ID for the list.
func dataSourceIBMIamTrustedProfileLinkListID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceProfileLinkListFlattenLinks(result []iamidentityv1.ProfileLink) (links []map[string]interface{}) {
	for _, linksItem := range result {
		links = append(links, dataSourceProfileLinkListLinksToMap(linksItem))
	}

	return links
}

func dataSourceProfileLinkListLinksToMap(linksItem iamidentityv1.ProfileLink) (linksMap map[string]interface{}) {
	linksMap = map[string]interface{}{}

	if linksItem.ID != nil {
		linksMap["id"] = linksItem.ID
	}
	if linksItem.EntityTag != nil {
		linksMap["entity_tag"] = linksItem.EntityTag
	}
	if linksItem.CreatedAt != nil {
		linksMap["created_at"] = linksItem.CreatedAt.String()
	}
	if linksItem.ModifiedAt != nil {
		linksMap["modified_at"] = linksItem.ModifiedAt.String()
	}
	if linksItem.Name != nil {
		linksMap["name"] = linksItem.Name
	}
	if linksItem.CrType != nil {
		linksMap["cr_type"] = linksItem.CrType
	}
	if linksItem.Link != nil {
		linkList := []map[string]interface{}{}
		linkMap := dataSourceProfileLinkListLinksLinkToMap(*linksItem.Link)
		linkList = append(linkList, linkMap)
		linksMap["link"] = linkList
	}

	return linksMap
}

func dataSourceProfileLinkListLinksLinkToMap(linkItem iamidentityv1.ProfileLinkLink) (linkMap map[string]interface{}) {
	linkMap = map[string]interface{}{}

	if linkItem.CRN != nil {
		linkMap["crn"] = linkItem.CRN
	}
	if linkItem.Namespace != nil {
		linkMap["namespace"] = linkItem.Namespace
	}
	if linkItem.Name != nil {
		linkMap["name"] = linkItem.Name
	}

	return linkMap
}
