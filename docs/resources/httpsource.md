---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "panther_httpsource Resource - terraform-provider-panther"
subcategory: ""
description: |-
  
---

# panther_httpsource (Resource)



## Example Usage

```terraform
# Manage Http Log Source integration
resource "panther_httpsource" "example_http_source" {
  integration_label = ""
  log_stream_type   = "JSON"
  log_types         = ""
  auth_method       = "SharedSecret"
  auth_header_key   = ""
  auth_secret_value = ""
  auth_username     = ""
  auth_password     = ""
  auth_hmac_alg     = ""
  auth_bearer_token = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `auth_method` (String) The authentication method of the http source
- `integration_label` (String) The integration label (name)
- `log_stream_type` (String) The log stream type. Supported log stream types: Auto, JSON, JsonArray, Lines, CloudWatchLogs, XML
- `log_types` (List of String) The log types of the integration

### Optional

- `auth_bearer_token` (String) The authentication bearer token value of the http source. Used for Bearer auth method
- `auth_header_key` (String) The authentication header key of the http source. Used for HMAC and SharedSecret auth methods
- `auth_hmac_alg` (String) The authentication algorithm of the http source. Used for HMAC auth method
- `auth_password` (String) The authentication header password of the http source. Used for Basic auth method
- `auth_secret_value` (String) The authentication header secret value of the http source. Used for HMAC and SharedSecret auth methods
- `auth_username` (String) The authentication header username of the http source. Used for Basic auth method
- `id` (String) ID of the http source to fetch
- `log_stream_type_options` (Attributes) (see [below for nested schema](#nestedatt--log_stream_type_options))

<a id="nestedatt--log_stream_type_options"></a>
### Nested Schema for `log_stream_type_options`

Optional:

- `json_array_envelope_field` (String) Path to the array value to extract elements from, only applicable if logStreamType is JsonArray. Leave empty if the input JSON is an array itself
