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

dependency "example_project" {
  config_path  = find_in_parent_folders(local.env.dependencies.example_project_dependency_path)
  mock_outputs = local.env.dependencies.example_project_mock_outputs

  mock_outputs_allowed_terraform_commands = ["init", "plan", "validate"]
}

dependency "with_sql_instance" {
  config_path  = find_in_parent_folders(local.env.dependencies.with_sql_instance_dependency_path)
  mock_outputs = local.env.dependencies.with_sql_instance_mock_outputs

  mock_outputs_allowed_terraform_commands = ["init", "plan", "validate"]
}

terraform {
  source = "git::git@github.com:terraform-google-modules/terraform-google-cloud-operations//modules/ops-agent-policy?ref=${local.versions.google_module_cloud_ops}"
}

inputs = {
  agent_rules     = local.inputs.agent_rules
  assignment_id   = local.inputs.assignment_id
  instance_filter = local.inputs.instance_filter
  project         = dependency.example_project.outputs.project_id
  zone            = format("%s-%s", local.region.region, local.region.zone_preference)
}
