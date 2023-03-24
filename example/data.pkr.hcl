packer {
  required_plugins {
    tss = {
      version = "0.3.2"
      source = "github.com/breed808/tss"
    }
  }
}

data "tss" "mock-data" {
  username   = "testing"
  password   = "test123"
  server_url = "https://my-thycotic-server.example.com/SecretServer"
  domain     = "testing.example.com"

  secret_id = "500"
  secret_fields = [
    "password",
    "username",
  ]
}
