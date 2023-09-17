resource "kubernetes_deployment" "bot" {
  metadata {
    name = var.gke_deployment_name
    labels = {
      app = var.gke_deployment_name
    }
  }

  spec {
    replicas = 1
    selector {
      match_labels = {
        app = var.gke_deployment_name
      }
    }
    template {
      metadata {
        labels = {
          app = var.gke_deployment_name
        }
      }
      spec {
        container {
          name    = var.gke_deployment_name
          image   = var.gcr_image_name
          command = ["./naka"]
          args    = ["bot"]
          env {
            name  = "NAKA_DISCORD_TOKEN"
            value = var.naka_discord_token
          }
          env {
            name  = "NAKA_DISCORD_PREFIX"
            value = var.naka_discord_prefix
          }
          env {
            name  = "NAKA_CACHE_DIALECT"
            value = var.naka_cache_dialect
          }
          env {
            name  = "NAKA_CACHE_ADDRESS"
            value = var.naka_cache_address
          }
          env {
            name  = "NAKA_CACHE_PASSWORD"
            value = var.naka_cache_password
          }
          env {
            name  = "NAKA_CACHE_TIME"
            value = var.naka_cache_time
          }
          env {
            name  = "NAKA_YOUTUBE_KEY"
            value = var.naka_youtube_key
          }
          env {
            name  = "NAKA_LOG_LEVEL"
            value = var.naka_log_level
          }
          env {
            name  = "NAKA_LOG_JSON"
            value = var.naka_log_json
          }
          env {
            name  = "NAKA_NEWRELIC_LICENSE_KEY"
            value = var.naka_newrelic_license_key
          }
        }
      }
    }
  }
}
