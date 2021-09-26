data "tss" "mock-data" {
  username   = "testing"
  password   = "test123"
  server_url = "https://my-thycotic-server.example.com/SecretServer"

  secret_id = "500"
  secret_fields = [
    "password",
    "username",
  ]
}
