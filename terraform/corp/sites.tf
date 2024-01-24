resource "sigsci_site" "app1_example_com" {
  short_name             = "app1.example.com"
  display_name           = "app1.example.com"
  block_duration_seconds = 600
  agent_anon_mode        = ""
  agent_level            = "block"
}

resource "sigsci_site" "app2_example_com" {
  short_name             = "app2.example.com"
  display_name           = "app2.example.com"
  block_duration_seconds = 600
  agent_anon_mode        = ""
  agent_level            = "block"
}
