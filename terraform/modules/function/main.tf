data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = var.source_file
  output_path = "${var.source_file}.zip"
}

# Lambda function policy
resource "aws_iam_policy" "lambda_policy" {
  name        = "${var.function_name}-lambda-policy"
  description = "${var.function_name}-lambda-policy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Effect": "Allow",
      "Resource": [
        "${var.bucket_arn}",
        "${var.bucket_arn}/*"
        ]
    },
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

# Lambda function role
resource "aws_iam_role" "iam_for_terraform_lambda" {
  name               = "${var.function_name}-lambda-role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

# Role to Policy attachment
resource "aws_iam_role_policy_attachment" "terraform_lambda_iam_policy_basic_execution" {
  role       = aws_iam_role.iam_for_terraform_lambda.id
  policy_arn = aws_iam_policy.lambda_policy.arn
}

resource "aws_lambda_function" "func" {
  filename         = data.archive_file.lambda_zip.output_path
  function_name    = var.function_name
  role             = aws_iam_role.iam_for_terraform_lambda.arn
  handler          = var.lambda_handler
  source_code_hash = filebase64sha256(data.archive_file.lambda_zip.output_path)
  runtime          = "go1.x"
  memory_size      = 1024
  timeout          = 30
}