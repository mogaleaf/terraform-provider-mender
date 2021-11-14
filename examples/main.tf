terraform {
  required_providers {
    mender = {
      source  = "mogaleaf.com/mogaleaf/mender"
    }
  }
}

provider "mender" {
  username = "<your mender connect email>"
  password = "<your mender password>"
  host     = "https://hosted.mender.io"
}

resource "mender_artifact" "my_release" {
  source_file = "<path of the mender artifact>"
  description = "here is a description of my artifact"
}

output "id" {
  value = mender_artifact.my_release.id
}
output "md5" {
  value = mender_artifact.my_release.md5
}
