resource "random_id" "ecommerce-instance-suffix" {
  byte_length = 2
}

resource "google_sql_database_instance" "ecommerce-instance" {
  name             = "ecommerce-${random_id.ecommerce-instance-suffix.hex}"
  project          = google_project.ecommerce-microservices.project_id
  region           = var.project_region
  database_version = "POSTGRES_14"

  deletion_protection = false

  settings {
    tier              = "db-f1-micro"
    activation_policy = "ALWAYS"
    disk_autoresize   = true

    disk_size         = 10
    disk_type         = "PD_SSD"
    availability_type = "ZONAL"
  }
}

resource "google_sql_database" "ecommerce-database" {
  name     = "ecommerce"
  instance = google_sql_database_instance.ecommerce-instance.name
  project  = google_project.ecommerce-microservices.project_id
}

resource "random_id" "user-password" {
  byte_length = 8
}

resource "google_sql_user" "application" {
  name     = "application"
  password = var.user-password-ecommerce-instance == "" ? random_id.user-password.hex : var.user-password-ecommerce-instance
  instance = google_sql_database_instance.ecommerce-instance.name
  project  = google_project.ecommerce-microservices.project_id
}
