source "null" "basic-example" {
  ssh_host= "127.0.0.1"
  ssh_username = "foo"
  ssh_password = "bar"
}

build {
  sources = ["sources.null.basic-example"]
}
