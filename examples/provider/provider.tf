terraform {
  required_providers {
    wiz = {
      source  = "tiagoasousa/wiz"
      version = "0.0.1"
    }
  }
}

provider "wiz" {
  wiz_url                = var.wiz_url
  wiz_auth_client_id     = var.wiz_auth_client_id
  wiz_auth_client_secret = var.wiz_auth_client_secret
  wiz_auth_audience      = "wiz-api"
}
