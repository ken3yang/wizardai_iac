// Output variables for the S3 Bucket module

output "bucket_name" {
  description = "The name of the bucket"
  value       = aws_s3_bucket.this.bucket
}

output "bucket_arn" {
  description = "The ARN of the bucket"
  value       = aws_s3_bucket.this.arn
}
