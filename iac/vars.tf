variable "project" {
  type    = string
  default = "ecommerce-"
}

variable "project_region" {
  type    = string
  default = "us-central1"
}

variable "project_zone" {
  type    = string
  default = "us-central1-c"
}

variable "user-password-ecommerce-instance" {
  type        = string
  default     = ""
  description = "The password that will be used for the default user in ecommerce Postgres instance. If none is provided, a random one will be generated."
}

variable "billing_account" {
  type        = string
  default     = ""
  sensitive   = true
  description = "Google Cloud billing account id"

  validation {
    condition     = length(var.billing_account) == 0
    error_message = "The billing_account must be specified"
  }
}
