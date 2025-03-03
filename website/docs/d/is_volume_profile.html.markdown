---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : Volume Profile"
description: |-
  Manages IBM Cloud virtual server volume profile.
---

# ibm_is_volume_profile
Retrieve information of an existing IBM Cloud virtual server volume profile as a read-only data source. For more information, about virtual server volume profile, see [restoring a volume from a snapshot](https://cloud.ibm.com/docs/vpc?topic=vpc-snapshots-vpc-restore).

**Note:** 
VPC infrastructure services are a regional specific based endpoint, by default targets to `us-south`. Please make sure to target right region in the provider block as shown in the `provider.tf` file, if VPC service is created in region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```

## Example usage

```terraform

data "ibm_is_volume_profile" "volprofile"{
  name = "general-purpose"
}

```

## Argument reference
Review the argument references that you can specify for your data source. 

- `name` - (Required, String) The name for the virtual server volume profile.

## Attribute reference
In addition to the argument reference list, you can access the following attribute references after your data source is created. 

- `family` - (String) The family of the virtual server volume profile.
