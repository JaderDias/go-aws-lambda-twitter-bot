provider "aws" {
  region = var.aws_region
}

resource "random_pet" "this" {
  length = 2
}

module "dynamodb_table" {
  source = "terraform-aws-modules/dynamodb-table/aws"

  name      = "my-table-${random_pet.this.id}"
  hash_key  = "Id"
  range_key = "Title"

  attributes = [
    {
      name = "Id"
      type = "N"
    },
    {
      name = "Title"
      type = "S"
    },
    {
      name = "Age"
      type = "N"
    }
  ]

  global_secondary_indexes = [
    {
      name               = "TitleIndex"
      hash_key           = "Title"
      range_key          = "Age"
      projection_type    = "INCLUDE"
      non_key_attributes = ["Id"]
    }
  ]

  tags = {
    Terraform   = "true"
    Environment = "staging"
  }
}

resource "aws_s3_bucket" "b" {
  bucket = "my-bucket-included-giraffe"
}

resource "aws_s3_bucket_object" "object" {
  bucket = aws_s3_bucket.b.id
  key    = "twitter_secrets"
  content = var.twitter_secrets
  etag = md5(var.twitter_secrets)
}

module "my_function" {
  source         = "./modules/function"

  function_name  = "tweet_function"
  lambda_handler = "tweet"
  source_file = "../bin/tweet"
  schedule_expression = "rate(60 minutes)"
  dynamodb_table_id = module.dynamodb_table.dynamodb_table_id
  dynamodb_table_arn = module.dynamodb_table.dynamodb_table_arn
}