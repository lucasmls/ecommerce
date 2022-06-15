terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.89.0"
    }

    random = {
      source  = "hashicorp/random"
      version = "3.2.0"
    }
  }
}

provider "google" {
  project = var.project
  region  = var.project_region
  zone    = var.project_zone
}

resource "random_id" "id" {
  prefix      = var.project
  byte_length = 4
}

data "google_billing_account" "billing_account" {
  billing_account = var.billing_account
}

resource "google_project" "ecommerce-microservices" {
  name            = "e-commerce microservices"
  project_id      = random_id.id.hex
  billing_account = data.google_billing_account.billing_account.id
}
