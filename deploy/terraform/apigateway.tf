resource "aws_api_gateway_rest_api" "api" {
  name = "${local.app}-${local.env}-api"
  description = "${local.app} ${local.env} api"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_gateway_response" "responses" {
  for_each = {
    BAD_REQUEST_BODY               = 400
    BAD_REQUEST_PARAMETERS         = 400
    DEFAULT_4XX                    = 400
    MISSING_AUTHENTICATION_TOKEN   = 400
    UNAUTHORIZED                   = 401
    ACCESS_DENIED                  = 403
    EXPIRED_TOKEN                  = 403
    INVALID_API_KEY                = 403
    INVALID_SIGNATURE              = 403
    WAF_FILTERED                   = 403
    RESOURCE_NOT_FOUND             = 404
    REQUEST_TOO_LARGE              = 413
    UNSUPPORTED_MEDIA_TYPE         = 415
    QUOTA_EXCEEDED                 = 429
    THROTTLED                      = 429
    API_CONFIGURATION_ERROR        = 500
    AUTHORIZER_CONFIGURATION_ERROR = 500
    AUTHORIZER_FAILURE             = 500
    DEFAULT_5XX                    = 500
    INTEGRATION_FAILURE            = 504
    INTEGRATION_TIMEOUT            = 504
  }
  response_type = each.key
  status_code   = each.value
  rest_api_id   = aws_api_gateway_rest_api.api.id
  response_templates = {
    "application/json" = "{\"status\":\"FAILURE\",\"errors\":[$context.error.messageString],\"message\":$context.error.messageString}"
  }
  depends_on = [
    aws_api_gateway_rest_api.api
  ]
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

# resource /health
resource "aws_api_gateway_resource" "health" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "health"
  depends_on = [
    aws_api_gateway_rest_api.api
  ]
}

# method /health
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

# resource /v1/user
resource "aws_api_gateway_resource" "user" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_resource.root.id
  path_part   = "user"
  depends_on = [
    aws_api_gateway_resource.root
  ]
}

# resource /v1/user/register
resource "aws_api_gateway_resource" "user_register" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_resource.user.id
  path_part   = "register"
  depends_on = [
    aws_api_gateway_resource.user
  ]
}

# method POST /v1/user/register
resource "aws_api_gateway_method" "post_user_register" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.user_register.id
  http_method = "POST"
  authorization = "NONE"
  depends_on = [
    aws_api_gateway_resource.user_register
  ]
}

resource "aws_api_gateway_integration" "user_register_integration" {
  rest_api_id             = aws_api_gateway_rest_api.api.id
  resource_id             = aws_api_gateway_resource.user_register.id
  http_method             = aws_api_gateway_method.post_user_register.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.main.invoke_arn
  depends_on = [
    aws_api_gateway_method.post_user_register
  ]

}

resource "aws_api_gateway_method_response" "post_user_register" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.user_register.id
  http_method = aws_api_gateway_method.post_user_register.http_method
  status_code = "200"
  depends_on = [
    aws_api_gateway_method.post_user_register
  ]
}

resource "aws_api_gateway_integration_response" "post_user_register" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.user_register.id
  http_method = aws_api_gateway_method.post_user_register.http_method
  status_code = aws_api_gateway_method_response.post_user_register.status_code
  depends_on = [
    aws_api_gateway_method.post_user_register,
    aws_api_gateway_method_response.post_user_register
  ]
}

# resource /v1/user/register/confirm
resource "aws_api_gateway_resource" "user_register_confirm" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_resource.user_register.id
  path_part   = "confirm"
  depends_on = [
    aws_api_gateway_resource.user_register
  ]
}

# method POST /v1/user/register/confirm
resource "aws_api_gateway_method" "post_user_register_confirm" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.user_register_confirm.id
  http_method = "POST"
  authorization = "NONE"
  depends_on = [
    aws_api_gateway_resource.user_register_confirm
  ]
}

resource "aws_api_gateway_integration" "user_register_confirm_integration" {
  rest_api_id             = aws_api_gateway_rest_api.api.id
  resource_id             = aws_api_gateway_resource.user_register_confirm.id
  http_method             = aws_api_gateway_method.post_user_register_confirm.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.main.invoke_arn
  depends_on = [
    aws_api_gateway_method.post_user_register_confirm
  ]

}

resource "aws_api_gateway_method_response" "post_user_register_confirm" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.user_register_confirm.id
  http_method = aws_api_gateway_method.post_user_register_confirm.http_method
  status_code = "200"
  depends_on = [
    aws_api_gateway_method.post_user_register_confirm
  ]
}

resource "aws_api_gateway_integration_response" "post_user_register_confirm" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.user_register_confirm.id
  http_method = aws_api_gateway_method.post_user_register_confirm.http_method
  status_code = aws_api_gateway_method_response.post_user_register_confirm.status_code
  depends_on = [
    aws_api_gateway_method.post_user_register_confirm,
    aws_api_gateway_method_response.post_user_register_confirm
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
    aws_api_gateway_integration.user_register_integration,
    aws_api_gateway_integration.user_register_confirm_integration,
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
