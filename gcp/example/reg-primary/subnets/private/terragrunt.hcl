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

terraform {
  source = "git::git@github.com:terraform-google-modules/terraform-google-network//modules/subnets?ref=${local.versions.google_module_network}"
}

inputs = {
  network_name = dependency.private_network.outputs.network_id
  project_id   = dependency.private_network.outputs.project_id
  subnets      = [for s in local.inputs.subnets: tomap({
                    subnet_name = s.name
                    subnet_ip = s.range
                    subnet_region = local.region.region})]
}
