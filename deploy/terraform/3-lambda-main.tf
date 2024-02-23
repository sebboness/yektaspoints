resource "aws_iam_role" "lambda_exec" {
  name = "hexonite-${local.app}-${local.env}-main-lambda"

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

variable "lambda_version" {
  default = "../../api/VERSION"
}

data "external" "output_hash" {
  program = ["/bin/sh", "${path.module}/compute_file_hash.sh", "${var.output_path}"]
}

resource "aws_lambda_function" "main" {
  function_name = "${local.app}-${local.env}-main"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_main.key

  package_type = "Zip"
  runtime = "provided.al2023"
  handler = "bootstrap"
  architectures = ["x86_64"]

  source_code_hash = data.external.output_hash.result.filebase64sha256

  role = aws_iam_role.lambda_exec.arn

  environment {
    variables = {
      APPNAME  = local.app
      BUILT_AT = timestamp()
      COGNITO_USER_POOL_ID  = local.ssm_secrets["COGNITO_USER_POOL_ID"]
      COGNITO_CLIENT_ID     = local.ssm_secrets["COGNITO_CLIENT_ID"]
      COGNITO_CLIENT_SECRET = local.ssm_secrets["COGNITO_CLIENT_SECRET"]
      ENV      = local.env
      GIN_MODE = local.env == "prod" ? "release" : "debug" 
      VERSION  = file(var.lambda_version)
    }
  }
}

resource "aws_cloudwatch_log_group" "main" {
  name = "/aws/lambda/${aws_lambda_function.main.function_name}"

  retention_in_days = 14
}

resource "aws_s3_object" "lambda_main" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key    = "${local.app}-${local.env}-main.zip"
  source = var.output_path

  etag = filemd5(var.output_path)
}

resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.main.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.api.execution_arn}/*/*/*"
}