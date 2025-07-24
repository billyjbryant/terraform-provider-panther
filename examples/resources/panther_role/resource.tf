# Manage user role and permissions
resource "panther_role" "example" {
  name = "example-role"
  permissions = [
    "RuleRead",
    "PolicyRead"
  ]
}