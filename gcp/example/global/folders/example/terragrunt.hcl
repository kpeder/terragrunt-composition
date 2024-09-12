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

terraform {
  source = "git::git@github.com:terraform-google-modules/terraform-google-folders?ref=${local.versions.google_module_folder}"
}

inputs = {
  names  = [for name in local.inputs.names: format("%s-%s", local.env.environment, name)]
  parent = format("%s/%s",
                  coalesce(local.inputs.parent.type, local.platform.parent.type),
                  coalesce(local.inputs.parent.id, local.platform.parent.id))
  prefix = try(local.platform.prefix, null)
}
