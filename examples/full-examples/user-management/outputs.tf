output "analyst_role_id" {
  description = "ID of the security analyst role"
  value       = panther_role.security_analyst.id
}

output "admin_role_id" {
  description = "ID of the security admin role"
  value       = panther_role.security_admin.id
}

output "analyst_user_id" {
  description = "ID of the analyst user"
  value       = panther_user.analyst.id
}

output "admin_user_id" {
  description = "ID of the admin user"
  value       = panther_user.admin.id
}