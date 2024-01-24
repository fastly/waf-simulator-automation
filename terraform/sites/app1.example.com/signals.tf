resource "sigsci_site_signal_tag" "sensitive_account_api" {
  site_short_name = var.site_name
  name            = "sensitive account api"
  description     = "sensitive account api endpoint"
}