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
  region   = yamldecode(file(find_in_parent_folders("region.yaml")))
  versions = yamldecode(file(find_in_parent_folders("versions.yaml")))
}

dependency "private_network" {
  config_path  = find_in_parent_folders(local.env.dependencies.private_network_dependency_path)
  mock_outputs = local.env.dependencies.private_network_mock_outputs

  mock_outputs_allowed_terraform_commands = ["init", "plan", "validate"]
}

dependency "secondary_subnets" {
  config_path  = find_in_parent_folders(local.env.dependencies.secondary_subnets_dependency_path)
  mock_outputs = local.env.dependencies.secondary_subnets_mock_outputs

  mock_outputs_allowed_terraform_commands = ["init", "plan", "validate"]
}

terraform {
  source = "git::git@github.com:terraform-google-modules/terraform-google-cloud-router?ref=${local.versions.google_module_router}"
}

inputs = {
  name    = format("%s-%s", local.region.region, local.inputs.name)
  network = dependency.private_network.outputs.network_id
  project = dependency.private_network.outputs.project_id
  region  = local.region.region

  nats = local.inputs.nat ? [{
    log_config                         = local.inputs.log_config
    name                               = format("%s-%s", local.region.region, local.inputs.name)
    source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
    subnetworks = [for s in local.inputs.source_subnets:
      {
        name                     = dependency.secondary_subnets.outputs.subnets[format("%s/%s", local.region.region, s)].id
        source_ip_ranges_to_nat  = local.inputs.source_subnet_ranges
      }]
  }] : []

}
