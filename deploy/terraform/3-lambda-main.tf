resource "aws_iam_role" "lambda_exec" {
  name = "hexonite-${var.app}-${var.env}-main-lambda"

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

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

variable "output_path" {
  default = "../../api/cmd/lambda/bootstrap.zip"
}

data "external" "output_hash" {
  program = ["/bin/sh", "${path.module}/compute_file_hash.sh", "${var.output_path}"]
}

resource "aws_lambda_function" "main" {
  function_name = "${var.app}-${var.env}-main"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_main.key

  package_type = "Zip"
  runtime = "provided.al2"
  handler = "bootstrap"
  architectures = ["x86_64"]

  source_code_hash = data.external.output_hash.result.filebase64sha256

  role = aws_iam_role.lambda_exec.arn

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

resource "aws_cloudwatch_log_group" "main" {
  name = "/aws/lambda/${aws_lambda_function.main.function_name}"

  retention_in_days = 14
}

resource "aws_s3_object" "lambda_main" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key    = "${var.app}-${var.env}-main.zip"
  source = var.output_path

  etag = filemd5(var.output_path)
}