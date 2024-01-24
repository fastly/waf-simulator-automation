resource "sigsci_site_rule" "sensitive_account_api" {
  site_short_name = var.site_name
  type            = "request"
  group_operator  = "all"
  requestlogging  = "sampled"
  enabled         = true
  reason          = "monitor sensitive account api endpoints"
  expiration      = ""

  conditions {
    type     = "single"
    field    = "path"
    operator = "inList"
    value    = sigsci_site_list.sensitive_account_api.id
  }

  actions {
    type   = "addSignal"
    signal = sigsci_site_signal_tag.sensitive_account_api.id
  }
}