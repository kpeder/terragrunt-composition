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
  versions = yamldecode(file(find_in_parent_folders("versions.yaml")))
}

dependency "example_project" {
  config_path  = find_in_parent_folders(local.env.dependencies.example_project_dependency_path)
  mock_outputs = local.env.dependencies.example_project_mock_outputs

  mock_outputs_allowed_terraform_commands = ["init", "plan", "validate"]
}

terraform {
  source = "git::git@github.com:terraform-google-modules/terraform-google-iam//modules/member_iam?ref=${local.versions.google_module_iam}"
}

inputs = {
  prefix                  = local.inputs.prefix
  project_id              = dependency.example_project.outputs.project_id
  project_roles           = local.inputs.roles
  service_account_address = dependency.example_project.outputs.service_account_email
}
