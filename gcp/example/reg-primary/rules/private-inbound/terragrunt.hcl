# Terragrunt will copy the Terraform configurations specified by the source parameter, along with any files in the
# working directory, into a temporary folder, and execute your Terraform commands in that folder.

# Include all settings from the root terragrunt.hcl file
include {
  path = find_in_parent_folders("example_terragrunt.hcl")
}

# Resources should not be destroyed without careful consideration of effects
prevent_destroy = false

locals {
  env      = yamldecode(file(find_in_parent_folders("env.yaml")))
  inputs   = yamldecode(file("inputs.yaml"))
  platform = fileexists(find_in_parent_folders("local.gcp.yaml")) ? yamldecode(file(find_in_parent_folders("local.gcp.yaml"))) : yamldecode(file(find_in_parent_folders("gcp.yaml")))
  region   = yamldecode(file(find_in_parent_folders("reg-multi/region.yaml")))
  versions = yamldecode(file(find_in_parent_folders("versions.yaml")))
}

dependency "private_network" {
  config_path  = find_in_parent_folders(local.env.dependencies.private_network_dependency_path)
  mock_outputs = local.env.dependencies.private_network_mock_outputs

  mock_outputs_allowed_terraform_commands = ["init", "plan", "validate"]
}

dependency "primary_subnets" {
  config_path  = find_in_parent_folders(local.env.dependencies.primary_subnets_dependency_path)
  mock_outputs = local.env.dependencies.primary_subnets_mock_outputs

  mock_outputs_allowed_terraform_commands = ["init", "plan", "validate"]
}

terraform {
  source = "git::git@github.com:terraform-google-modules/terraform-google-network//modules/firewall-rules?ref=${local.versions.google_module_network}"
}

inputs = {
  project_id   = dependency.private_network.outputs.project_id
  network_name = dependency.private_network.outputs.network_id

  ingress_rules = [for rule in local.inputs.rules:
    {
      name                    = rule["name"]
      description             = rule["description"]
      direction               = rule["direction"]
      disabled                = rule["disabled"]
      enable_logging          = rule["enable_logging"]
      log_config              = rule["log_config"]
      priority                = rule["priority"]
      source_ranges           = concat(rule["source_ranges"], [for s in dependency.primary_subnets.outputs.subnets: s["ip_cidr_range"]])
      target_tags             = rule["target_tags"]
      allow = [{
        protocol = rule["protocol"]
        ports    = rule["ports"]
      }]
    }]
}
