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

