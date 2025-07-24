# Define roles first
resource "panther_role" "security_analyst" {
  name        = "SecurityAnalyst"
  description = "Security analyst with read access to security data"

  permissions = [
    "RuleRead",
    "PolicyRead",
    "AlertRead",
    "DataModelRead"
  ]
}

resource "panther_role" "security_admin" {
  name        = "SecurityAdmin"
  description = "Security administrator with full access"

  permissions = [
    "RuleRead",
    "RuleWrite",
    "RuleDelete",
    "PolicyRead",
    "PolicyWrite",
    "PolicyDelete",
    "AlertRead",
    "AlertWrite",
    "UserRead",
    "UserWrite"
  ]
}

# Create users with role assignments
resource "panther_user" "analyst" {
  email       = var.analyst_email
  given_name  = "Security"
  family_name = "Analyst"
  role        = panther_role.security_analyst.id
}

resource "panther_user" "admin" {
  email       = var.admin_email
  given_name  = "Security"
  family_name = "Admin"
  role        = panther_role.security_admin.id
}