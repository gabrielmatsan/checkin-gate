# S3 for backend state
resource "aws_s3_bucket" "backend-state" {

  bucket = "checkin-gate-terraform-state"

  tags = {
    Name = "checkin-gate-terraform-state"
    Environment = var.environment
    IaC = "true"
  }
}


resource "aws_s3_bucket_versioning" "backend-state-versioning" {
  bucket = aws_s3_bucket.backend-state.id
  versioning_configuration {
    status = "Enabled"
  }
}


resource "aws_s3_bucket_server_side_encryption_configuration" "backend-state-encryption" {
  bucket = aws_s3_bucket.backend-state.id

  rule {
    
  }
}