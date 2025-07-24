variable "token" {
  description = "Panther API token"
  type        = string
}

variable "url" {
  description = "Panther API URL"
  type        = string
}

variable "admin_email" {
  description = "Email address for the admin user"
  type        = string
  default     = "admin@company.com"
}

variable "analyst_email" {
  description = "Email address for the analyst user"
  type        = string
  default     = "analyst@company.com"
}