variable "name" {
  description = "Name for the S3 Bucket"
  type        = string
}

variable "region" {
  description = "Region for the S3 Bucket"
  type        = string
}

variable "environment" {
  description = "Environment for the S3 Bucket (development, staging, production)"
  type        = string
}

variable "tags" {
  description = "Tags for the S3 Bucket"
  type        = map(string)
  default     = {}
}
