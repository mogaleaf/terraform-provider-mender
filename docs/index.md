---
page_title: "Provider: Mender"
subcategory: ""
description: |-
Terraform provider for interacting with mender API.
---

# Mender Provider

-> visit the [Mender api](https://docs.mender.io/api/)
Mender api provider is used to interact with the managements apis of mender.

## Example Usage

```terraform
provider "mender" {
  username = "myEmail@gmail.com"
  password = "PasswordToGetToMenderUi"
  host     = "https://hosted.mender.io"
}
```

## Schema

### Optional


- **username** (String, Optional) Username to authenticate to Mender API, can be set by environment variable `MENDER_USERNAME`
- **password** (String, Optional) Password to authenticate to Mender API, can be set by environment variable `MENDER_PASSWORD`
- **host** (String, Optional) Mender API address (defaults to `https://hosted.mender.io`), can be set by environment variable `MENDER_HOST`
