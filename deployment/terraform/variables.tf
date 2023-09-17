variable "gcp_project_id" {
  type        = string
  description = "GCP project id"
}

variable "gcp_region" {
  type        = string
  description = "GCP project region"
}

variable "gke_cluster_name" {
  type        = string
  description = "GKE cluster name"
}

variable "gke_location" {
  type        = string
  description = "GKE location"
}

variable "gke_pool_name" {
  type        = string
  description = "GKE node pool name"
}

variable "gke_node_preemptible" {
  type        = bool
  description = "GKE node preemptible"
}

variable "gke_node_machine_type" {
  type        = string
  description = "GKE node machine type"
}

variable "gcr_image_name" {
  type        = string
  description = "GCR image name"
}

variable "gke_deployment_name" {
  type        = string
  description = "GKE deployment bot name"
}

variable "naka_discord_token" {
  type        = string
  description = "Discord token"
}

variable "naka_discord_prefix" {
  type        = string
  description = "Discord command prefix"
}

variable "naka_cache_dialect" {
  type        = string
  description = "Cache dialect"
}

variable "naka_cache_address" {
  type        = string
  description = "Cache address"
}

variable "naka_cache_password" {
  type        = string
  description = "Cache password"
}

variable "naka_cache_time" {
  type        = string
  description = "Cache time"
}

variable "naka_youtube_key" {
  type        = string
  description = "Youtube API key"
}

variable "naka_log_level" {
  type        = number
  description = "Log level"
}

variable "naka_log_json" {
  type        = bool
  description = "Log json"
}

variable "naka_newrelic_license_key" {
  type        = string
  description = "Newrelic license key"
}
