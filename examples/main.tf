terraform {
  required_providers {
    config-service = {
      version = "0.1.2"
      source  = "RevelSystems/config-service"
    }
  }
}

provider "config-service" {
  token     = "<token_value>"
  base_url  = "<base_url>"
}

resource "configuration" "config" {
  provider = config-service
  client = "client_name"
  attributes_json = jsonencode({
    url         = "https://google.com",
    some_flag   = false,
    some_number = 123,
    more        = null
  })
}

output "attributes_json" {
  value = configuration.config.attributes_json
}
output "client" {
 value = configuration.config.client
}

output "created_on" {
 value = configuration.config.created_on
}

output "updated_on" {
 value = configuration.config.updated_on
}
