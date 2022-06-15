resource "google_project_service" "compute" {
  service                    = "compute.googleapis.com"
  project                    = google_project.ecommerce-microservices.project_id
  disable_dependent_services = true
  depends_on = [
    google_project.ecommerce-microservices
  ]
}

resource "google_project_service" "container" {
  service                    = "container.googleapis.com"
  disable_dependent_services = true
  project                    = google_project.ecommerce-microservices.project_id
  depends_on = [
    google_project.ecommerce-microservices
  ]
}

resource "google_container_cluster" "ecommerce-cluster" {
  name     = "ecommerce-cluster"
  location = var.project_zone
  project  = google_project.ecommerce-microservices.project_id

  remove_default_node_pool = true
  initial_node_count       = 1

  depends_on = [google_project_service.compute, google_project_service.container]
}

resource "google_container_node_pool" "ecommerce-pool" {
  name       = "ecommerce-pool"
  cluster    = google_container_cluster.ecommerce-cluster.name
  location   = var.project_zone
  project    = google_project.ecommerce-microservices.project_id
  node_count = 1

  depends_on = [
    google_container_cluster.ecommerce-cluster
  ]

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    preemptible  = true
    machine_type = "n1-standard-2"
    tags         = ["gke-node", "${var.project}-gke"]
    metadata = {
      disable-legacy-endpoints = "true"
    }
  }
}
