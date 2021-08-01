packer {
  required_plugins {
    scaffolding = {
      version = ">=v0.1.0"
      source  = "github.com/hashicorp/scaffolding"
    }
  }
}

source "scaffolding-my-builder" "username-example" {
  mock = local.username
}

source "scaffolding-my-builder" "password-example" {
  mock = local.password
}

build {
  sources = [
    "source.scaffolding-my-builder.foo-example",
  ]

  source "source.scaffolding-my-builder.bar-example" {
    name = "bar"
  }

  provisioner "scaffolding-my-provisioner" {
    only = ["scaffolding-my-builder.foo-example"]
    mock = "foo: ${local.foo}"
  }

  provisioner "scaffolding-my-provisioner" {
    only = ["scaffolding-my-builder.bar"]
    mock = "bar: ${local.bar}"
  }

  post-processor "scaffolding-my-post-processor" {
    mock = "post-processor mock-config"
  }
}
