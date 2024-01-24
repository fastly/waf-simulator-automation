module "app1_example_com" {
  source    = "./sites/app1.example.com"
  site_name = module.corp.app1_example_com
}

module "app2_example_com" {
  source    = "./sites/app2.example.com"
  site_name = module.corp.app2_example_com
}
