variable "aws_region" {
  type        = string
  description = "The AWS region to deploy the infrastructure"
  default     = "us-east-1"
}

variable "environment" {
  type        = string
  description = "The environment to deploy the infrastructure"
  default     = "production"
}

variable "project_name" {
  type        = string
  description = "The project name to deploy the infrastructure"
  default     = "checkin-gate"
}

# variable "github_username" {
#   description = "GitHub username"
#   type        = string
# }

# variable "github_token" {
#   description = "GitHub personal access token"
#   type        = string
#   sensitive   = true
# }