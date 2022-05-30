provider "aws" {
  region = var.aws_region
}

resource "random_pet" "this" {
  length = 2
}

resource "aws_s3_bucket" "bucket" {
  bucket = "my-bucket-included-giraffe"
  acl    = "private"
}

resource "aws_s3_bucket_object" "object" {
  bucket  = aws_s3_bucket.bucket.id
  key     = "twitter_secrets"
  content = var.twitter_secrets
  etag    = md5(var.twitter_secrets)
}

module "my_function" {
  source = "./modules/function"

  function_name       = "tweet_function-${random_pet.this.id}"
  lambda_handler      = "tweet"
  source_file         = "../bin/tweet"
  schedule_expression = "rate(60 minutes)"

  bucket_arn = aws_s3_bucket.bucket.arn
}