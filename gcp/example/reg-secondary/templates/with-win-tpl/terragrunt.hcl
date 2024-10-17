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

dependency "secondary_subnets" {
  config_path  = find_in_parent_folders(local.env.dependencies.secondary_subnets_dependency_path)
  mock_outputs = local.env.dependencies.secondary_subnets_mock_outputs

  mock_outputs_allowed_terraform_commands = ["init", "plan", "validate"]
}

terraform {
  source = "git::git@github.com:terraform-google-modules/terraform-google-vm//modules/instance_template?ref=${local.versions.google_module_vm}"
}

inputs = {
  additional_disks     = [for disk in local.inputs.additional_disks: merge(
                            {for k, v in disk: k => v if k != "disk_labels"},
                            {for k, v in disk: k => merge(local.env.labels, local.inputs.labels, v) if k == "disk_labels" })]
  auto_delete          = local.inputs.auto_delete
  description          = local.inputs.description
  disk_labels          = merge(local.env.labels, local.inputs.labels, local.inputs.disk_labels)
  disk_size_gb         = local.inputs.disk_size_gb
  disk_type            = local.inputs.disk_type
  gpu                  = local.inputs.gpu
  instance_description = local.inputs.instance_description
  labels               = merge(local.env.labels, local.inputs.labels)
  machine_type         = local.inputs.machine_type
  name_prefix          = format("%s-%s-%s", local.platform.prefix, local.env.environment, local.inputs.name_prefix)
  project_id           = dependency.example_project.outputs.project_id
  region               = local.region.region
  service_account      = {
                           email  = coalesce(local.inputs.service_account.email, dependency.example_project.outputs.service_account_email)
                           scopes = toset(local.inputs.service_account.scopes)
                         }
  source_image         = local.inputs.source_image.image
  source_image_family  = local.inputs.source_image.family
  source_image_project = local.inputs.source_image.project
  spot                 = local.inputs.spot
  startup_script       = local.inputs.startup_script
  subnetwork           = dependency.secondary_subnets.outputs.subnets[format("%s/%s", local.region.region, local.inputs.subnetwork)].id
  tags                 = local.inputs.tags
}
