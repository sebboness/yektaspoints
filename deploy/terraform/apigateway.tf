resource "aws_api_gateway_rest_api" "api" {
  name = "${local.app}-${local.env}-api"
  description = "${local.app} ${local.env} api"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

data "aws_cognito_user_pools" "pools" {
  name = "mypoints"
}

resource "aws_api_gateway_authorizer" "cognito" {
  name          = "${local.app}-authorizer"
  rest_api_id   = aws_api_gateway_rest_api.api.id
  type          = "COGNITO_USER_POOLS"
  provider_arns = tolist(data.aws_cognito_user_pools.pools.arns)
}

# Root (/)
resource "aws_api_gateway_resource" "root" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "v1"
}

resource "aws_api_gateway_resource" "root_v1" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_resource.root.id
  path_part   = "{proxy+}"
  depends_on  = [
    aws_api_gateway_resource.root
  ]
}

resource "aws_api_gateway_method" "root_get" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.root_v1.id
  http_method = "ANY"
  authorization = "COGNITO_USER_POOLS"
  authorizer_id = aws_api_gateway_authorizer.cognito.id
}

resource "aws_api_gateway_integration" "root_integration" {
  rest_api_id             = aws_api_gateway_rest_api.api.id
  resource_id             = aws_api_gateway_resource.root_v1.id
  http_method             = aws_api_gateway_method.root_get.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.main.invoke_arn
}

resource "aws_api_gateway_method_response" "root_get" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.root_v1.id
  http_method = aws_api_gateway_method.root_get.http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true,
    "method.response.header.Access-Control-Allow-Methods" = true,
    "method.response.header.Access-Control-Allow-Origin"  = true
  }

}

resource "aws_api_gateway_integration_response" "root_int_resp" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.root_v1.id
  http_method = aws_api_gateway_method.root_get.http_method
  status_code = aws_api_gateway_method_response.root_get.status_code

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'${local.corsAllowHeaders}'",
    "method.response.header.Access-Control-Allow-Methods" = "'${local.corsAllowMethods}'",
    "method.response.header.Access-Control-Allow-Origin"  = "'${local.corsAllowOrigins}'"
  }

  depends_on = [
    aws_api_gateway_method.root_get,
    aws_api_gateway_integration.root_integration
  ]
}

# options for /v1
module "apigw_root_options" {
  source = "./apigw-options"
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.root_v1.id
  corsAllowHeaders = local.corsAllowHeaders
  corsAllowMethods = local.corsAllowMethods
  corsAllowOrigins = local.corsAllowOrigins
}

# /auth
resource "aws_api_gateway_resource" "auth" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "auth"
  depends_on = [
    aws_api_gateway_rest_api.api
  ]
}

# /auth/token
resource "aws_api_gateway_resource" "auth_token" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_resource.auth.id
  path_part   = "token"
  depends_on = [
    aws_api_gateway_resource.auth
  ]
}

resource "aws_api_gateway_method" "auth_token_post" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.auth_token.id
  http_method = "POST"
  authorization = "NONE"
  depends_on = [
    aws_api_gateway_resource.auth_token
  ]
}

resource "aws_api_gateway_integration" "auth_token_integration" {
  rest_api_id             = aws_api_gateway_rest_api.api.id
  resource_id             = aws_api_gateway_resource.auth_token.id
  http_method             = aws_api_gateway_method.auth_token_post.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.main.invoke_arn
  depends_on = [
    aws_api_gateway_method.auth_token_post
  ]

}

resource "aws_api_gateway_method_response" "auth_token_post" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.auth_token.id
  http_method = aws_api_gateway_method.auth_token_post.http_method
  status_code = "200"
  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true,
    "method.response.header.Access-Control-Allow-Methods" = true,
    "method.response.header.Access-Control-Allow-Origin"  = true
  }
  depends_on = [
    aws_api_gateway_method.auth_token_post
  ]
}

resource "aws_api_gateway_integration_response" "auth_token_post" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.auth_token.id
  http_method = aws_api_gateway_method.auth_token_post.http_method
  status_code = aws_api_gateway_method_response.auth_token_post.status_code
  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'${local.corsAllowHeaders}'",
    "method.response.header.Access-Control-Allow-Methods" = "'${local.corsAllowMethods}'",
    "method.response.header.Access-Control-Allow-Origin"  = "'${local.corsAllowOrigins}'"
  }
  depends_on = [
    aws_api_gateway_method.auth_token_post,
    aws_api_gateway_method_response.auth_token_post
  ]
}

# /health
resource "aws_api_gateway_resource" "health" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "health"
  depends_on = [
    aws_api_gateway_rest_api.api
  ]
}

resource "aws_api_gateway_method" "health_get" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.health.id
  http_method = "GET"
  authorization = "NONE"
  depends_on = [
    aws_api_gateway_resource.health
  ]
}

resource "aws_api_gateway_integration" "health_integration" {
  rest_api_id             = aws_api_gateway_rest_api.api.id
  resource_id             = aws_api_gateway_resource.health.id
  http_method             = aws_api_gateway_method.health_get.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.main.invoke_arn
  depends_on = [
    aws_api_gateway_method.health_get
  ]

}

resource "aws_api_gateway_method_response" "health_get" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.health.id
  http_method = aws_api_gateway_method.health_get.http_method
  status_code = "200"
  depends_on = [
    aws_api_gateway_method.health_get
  ]
}

resource "aws_api_gateway_integration_response" "health_get" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.health.id
  http_method = aws_api_gateway_method.health_get.http_method
  status_code = aws_api_gateway_method_response.health_get.status_code
  depends_on = [
    aws_api_gateway_method.health_get,
    aws_api_gateway_method_response.health_get
  ]
}

# Deployment and domain
resource "aws_api_gateway_deployment" "deployment" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  triggers = {
    redeployment = filesha1("${path.module}/apigateway.tf")
  }
  lifecycle {
    create_before_destroy = true
  }
  depends_on = [
    aws_api_gateway_integration.root_integration,
    aws_api_gateway_integration.auth_token_integration,
    aws_api_gateway_integration.health_integration,
    module.apigw_root_options,
  ]
}

resource "aws_api_gateway_stage" "stage" {
  deployment_id = aws_api_gateway_deployment.deployment.id
  rest_api_id   = aws_api_gateway_rest_api.api.id
  stage_name    = "deploy"
  depends_on = [
    aws_api_gateway_deployment.deployment
  ]
}
