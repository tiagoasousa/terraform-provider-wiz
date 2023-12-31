---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "wiz_control_associations Resource - terraform-provider-wiz"
subcategory: ""
description: |-
  Manage associations between controls and security sub-categories. Associations defined outside this resouce will remain untouched through the lifecycle of this resource. Wiz managed controls cannot be associated to Wiz managed security sub-categories. This resource does not support imports; it can, however, overlay existing resources to bring them under management.
---

# wiz_control_associations (Resource)

Manage associations between controls and security sub-categories. Associations defined outside this resouce will remain untouched through the lifecycle of this resource. Wiz managed controls cannot be associated to Wiz managed security sub-categories. This resource does not support imports; it can, however, overlay existing resources to bring them under management.

## Example Usage

```terraform
resource "wiz_control_associations" "test" {
  security_sub_category_ids = [
    "2e5bc0d5-835b-4b4c-99cf-b1c6ace90a52",
    "708ec4a1-1a5c-4cb3-9c52-511229c5bb35",
  ]
  control_ids = [
    "301e5fd0-6a1a-42a7-99f5-3b0436d55a7f",
    "a5fbd955-ed78-445a-827a-06d6cbe5aab2",
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `control_ids` (List of String) List of control IDs.
- `security_sub_category_ids` (List of String) List of security sub-category IDs.

### Optional

- `details` (String) Details of the association. This information is not used to manage resources but can serve as notes or documentation for the associations.
    - Defaults to `undefined`.

### Read-Only

- `id` (String) Internal identifier for the association.
