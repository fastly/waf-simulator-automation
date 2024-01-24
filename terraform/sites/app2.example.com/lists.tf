resource "sigsci_site_list" "allowed_hosts" {
  site_short_name = var.site_name
  name            = "allowed hosts"
  type            = "string"
  description     = "a list of domains this site can serve"
  entries = [
    "www.app2.example.com",
    "app2.example.com",
  ]
}
