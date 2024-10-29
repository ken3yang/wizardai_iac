// Main Terraform configuration for the S3 Bucket module
provider "aws" {
  region = var.region
}

resource "aws_kms_key" "s3" {
  description             = "This key is used to encrypt bucket objects"
  deletion_window_in_days = 10
}

resource "aws_s3_bucket" "this" {
    bucket = "wizardai-${var.name}-${var.environment}"
    tags = merge(
        var.tags,
        {
        "Name" = "wizardai-${var.name}-${var.environment}"
        }
  )
  
}

// Enforce encryption at rest
resource "aws_s3_bucket_server_side_encryption_configuration" "encryption_at_rest" {
  bucket = aws_s3_bucket.this.id

  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = aws_kms_key.s3.arn
      sse_algorithm     = "aws:kms"
    }
  }
}

// Enforce encryption in transit for this s3 bucket
resource "aws_s3_bucket_policy" "encryption_in_transit_policy" {
  bucket = aws_s3_bucket.this.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action    = "s3:*",
        Effect    = "Deny",
        Principal = "*",
        Resource  = [
          "${aws_s3_bucket.this.arn}/*",
          aws_s3_bucket.this.arn
        ],
        Condition = {
          Bool: {
            "aws:SecureTransport" = "false"
          }
        }
      }
    ]
  })
}

// block all public access for security
resource "aws_s3_bucket_public_access_block" "this" {
  bucket = aws_s3_bucket.this.id

  block_public_acls       = true
  ignore_public_acls      = true
  block_public_policy     = true
  restrict_public_buckets = true
}

// enable versioning
resource "aws_s3_bucket_versioning" "this" {
  bucket = aws_s3_bucket.this.id

  versioning_configuration {
    status = "Enabled"
  }
}