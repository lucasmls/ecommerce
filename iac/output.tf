output "project_id" {
  description = "Google Cloud Project ID"
  value       = google_project.ecommerce-microservices.project_id
}

output "ecommerce-instance" {
  description = "The name of the ecommerce instance created"
  value       = google_sql_database_instance.ecommerce-instance.name
}

output "instance_address" {
  description = "The IPv4 address of the ecommerce database instance"
  value       = google_sql_database_instance.ecommerce-instance.ip_address.0.ip_address
}

output "generated_user_password" {
  description = "The auto generated default user password if no input password was provided"
  value       = random_id.user-password.hex
  sensitive   = true
}
