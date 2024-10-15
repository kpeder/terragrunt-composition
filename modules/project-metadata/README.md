## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [google_compute_project_metadata_item.this](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_project_metadata_item) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_metadata_items"></a> [metadata\_items](#input\_metadata\_items) | Dictionary of metadata keys and values to set | `map(string)` | `{}` | no |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The ID of the project in which to set metadata values | `string` | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_metadata_items"></a> [metadata\_items](#output\_metadata\_items) | Map of metadata items set on the project |
