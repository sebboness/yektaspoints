resource "aws_iam_role" "health_lambda_exec" {
  name = "${var.app}-${var.env}-health-lambda"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "health_lambda_policy" {
  role       = aws_iam_role.health_lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

variable "output_path" {
  default = "../../api/cmd/lambda/health/health.zip"
}

data "external" "output_hash" {
  program = ["/bin/sh", "${path.module}/compute_file_hash.sh", "${var.output_path}"]
}

resource "aws_lambda_function" "health" {
  function_name = "${var.app}-${var.env}-health"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_health.key

  package_type = "Zip"
  runtime = "provided.al2023"
  handler = "bootstrap.handler"
  architectures = ["amd64"]

  source_code_hash = data.external.output_hash.result.filebase64sha256

  role = aws_iam_role.health_lambda_exec.arn

  environment {
    variables = {
      APPNAME = var.app
      ENV     = var.env
      COGNITO_USER_POOL_ID  = local.ssm_secrets["COGNITO_USER_POOL_ID"]
      COGNITO_CLIENT_ID     = local.ssm_secrets["COGNITO_CLIENT_ID"]
      COGNITO_CLIENT_SECRET = local.ssm_secrets["COGNITO_CLIENT_SECRET"]
    }
  }
}

resource "aws_cloudwatch_log_group" "health" {
  name = "/aws/lambda/${aws_lambda_function.health.function_name}"

  retention_in_days = 14
}

resource "aws_s3_object" "lambda_health" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key    = "${var.env}_health.zip"
  source = var.output_path

  etag = filemd5(var.output_path)
}