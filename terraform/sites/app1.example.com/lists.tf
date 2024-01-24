resource "sigsci_site_list" "sensitive_account_api" {
  site_short_name = var.site_name
  name            = "sensitive account api endpoints"
  type            = "string"
  description     = "sensitive account api endpoints"
  entries = [
    "/api/v1/account/reset_password",
    "/api/v1/account/update_profile",
    "/api/v1/account/delete_profile",
    "/api/v1/account/reset_api_key"
  ]
}
