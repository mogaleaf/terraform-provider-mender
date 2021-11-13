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

resource "mender_artifact" "reboot" {
  source_file = "<path of the mender artifact>"
}

output "id" {
  value = mender_artifact.reboot.id
}
output "md5" {
  value = mender_artifact.reboot.md5
}
