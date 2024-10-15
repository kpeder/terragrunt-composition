variable "project_id" {
    description = "The ID of the project in which to set metadata values"
    type        = string

    default = null
}

variable "metadata_items" {
    description = "Dictionary of metadata keys and values to set"
    type        = map(string)

    default = {}
}

resource "google_compute_project_metadata_item" "this" {
    for_each = var.metadata_items

    project = var.project_id
    key     = each.key
    value   = each.value
}

output "metadata_items" {
    description = "Map of metadata items set on the project"
    value       = {for k, v in google_compute_project_metadata_item.this: k => v}
}
