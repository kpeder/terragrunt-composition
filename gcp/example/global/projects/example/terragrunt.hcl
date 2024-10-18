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

dependency "example_folder" {
  config_path  = find_in_parent_folders(local.env.dependencies.example_folder_dependency_path)
  mock_outputs = local.env.dependencies.example_folder_mock_outputs

  mock_outputs_allowed_terraform_commands = ["init", "plan", "validate"]
}

terraform {
  source = "git::git@github.com:terraform-google-modules/terraform-google-project-factory?ref=${local.versions.google_module_project}"
}

inputs = {
  activate_apis           = local.inputs.activate_apis
  auto_create_network     = local.inputs.auto_create_network
  billing_account         = coalesce(local.inputs.billing_account_override, local.platform.billing_account)
  default_service_account = local.inputs.default_service_account
  folder_id               = dependency.example_folder.outputs.id
  labels                  = local.env.labels
  name                    = format("%s-%s-%s", local.platform.prefix, local.env.environment, local.inputs.project_name)
  org_id                  = local.platform.organization_id
  project_id              = local.inputs.project_id_override
  random_project_id       = local.inputs.random_project_id
}
