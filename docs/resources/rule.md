---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "panther_rule Resource - terraform-provider-panther"
subcategory: ""
description: |-
  
---

# panther_rule (Resource)



## Example Usage

```terraform
# Manage detection rule for log analysis
resource "panther_rule" "example" {
  display_name         = ""
  body                 = ""
  severity             = ""
  description          = ""
  enabled              = true
  dedup_period_minutes = 60
  log_types = [
    ""
  ]
  tags = [
    ""
  ]
  runbook = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `body` (String) The python body of the rule
- `severity` (String)

### Optional

- `dedup_period_minutes` (Number) The amount of time in minutes for grouping alerts
- `description` (String) The description of the rule
- `display_name` (String) The display name of the rule
- `enabled` (Boolean) Determines whether or not the rule is active
- `inline_filters` (String) The filter for the rule represented in YAML
- `log_types` (List of String) log types
- `managed` (Boolean) Determines if the rule is managed by panther
- `output_ids` (List of String) Destination IDs that override default alert routing based on severity
- `reports` (Map of List of String) reports
- `runbook` (String) How to handle the generated alert
- `summary_attributes` (List of String) A list of fields in the event to create top 5 summaries for
- `tags` (List of String) The tags for the rule
- `tests` (Attributes List) Unit tests for the Rule. Best practice is to include a positive and negative case (see [below for nested schema](#nestedatt--tests))
- `threshold` (Number) the number of events that must match before an alert is triggered

### Read-Only

- `created_at` (String)
- `created_by` (Attributes) The actor who created the rule (see [below for nested schema](#nestedatt--created_by))
- `created_by_external` (String) The text of the user-provided CreatedBy field when uploaded via CI/CD
- `id` (String) The ID of this resource.
- `last_modified` (String)

<a id="nestedatt--tests"></a>
### Nested Schema for `tests`

Required:

- `expected_result` (Boolean) The expected result
- `name` (String) name
- `resource` (String) resource

Optional:

- `mocks` (List of Map of String) mocks


<a id="nestedatt--created_by"></a>
### Nested Schema for `created_by`

Read-Only:

- `id` (String)
- `type` (String)
