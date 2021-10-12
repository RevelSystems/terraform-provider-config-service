terraform {
  required_providers {
    config-service = {
      version = "0.1.0"
      source  = "revelsystems.com/revel/config-service"
    }
  }
}

provider "config-service" {
  token     = "<token_value>"
  base_url  = "<base_url>"
}

resource "configuration" "webhook_service_configuration" {
  provider = config-service
  client = "webhook-service-test7"
  attributes_json = jsonencode({
    url         = "https://google.com",
    some_flag   = false,
    some_number = 123,
    more        = null
  })
}

output "attributes_json" {
  value = configuration.webhook_service_configuration.attributes_json
}
output "client" {
 value = configuration.webhook_service_configuration.client
}

output "created_on" {
 value = configuration.webhook_service_configuration.created_on
}

output "updated_on" {
 value = configuration.webhook_service_configuration.updated_on
}
