---
page_title: "order Resource - terraform-provider-mender"
subcategory: ""
description: |-
The artifact resource allows you to upload an artifact on mender.
---

# Resource `mender_artifact`

-> visit the [artifact upload](https://docs.mender.io/api/#management-api-deployments-upload-artifact)

The artifact resource allows you to upload a mender artifact on your mender account.


## Example Usage

```terraform
resource "mender_artifact" "my_release" {
  source_file = "/var/build/artifact.mender"
  description = "here is a description of my artifact"
}
```


## Argument Reference

 - `source_file` - (Required) path the the binary to upload
 - `description` - description of the artifact
