resource "sigsci_site_rule" "invalid_host_header" {
  site_short_name = var.site_name
  type            = "request"
  group_operator  = "any"
  requestlogging  = "sampled"
  enabled         = true
  reason          = "blocks requests with invalid host headers"
  expiration      = ""

  conditions {
    type     = "single"
    field    = "domain"
    operator = "notInList"
    value    = sigsci_site_list.allowed_hosts.id
  }

  actions {
    type   = "addSignal"
    signal = sigsci_site_signal_tag.invalid_host_header.id
  }

  actions {
    type   = "block"
    signal = ""
  }
}
