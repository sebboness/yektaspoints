# # Health check (/v1/health)
# resource "aws_api_gateway_resource" "health" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   parent_id   = aws_api_gateway_resource.root.id
#   path_part   = "health"
# }

# resource "aws_api_gateway_method" "proxy" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   resource_id = aws_api_gateway_resource.root.id
#   http_method = "POST"
#   authorization = "NONE"
# #   authorization = "COGNITO_USER_POOLS"
# #   authorizer_id = aws_api_gateway_authorizer.demo.id
# }

# resource "aws_api_gateway_integration" "lambda_integration" {
#   rest_api_id             = aws_api_gateway_rest_api.api.id
#   resource_id             = aws_api_gateway_resource.root.id
#   http_method             = aws_api_gateway_method.proxy.http_method
#   integration_http_method = "POST"
#   type                    = "AWS"
#   uri                     = aws_lambda_function.html_lambda.invoke_arn
# }

# resource "aws_api_gateway_method_response" "proxy" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   resource_id = aws_api_gateway_resource.root.id
#   http_method = aws_api_gateway_method.proxy.http_method
#   status_code = "200"

#   //cors section
#   response_parameters = {
#     "method.response.header.Access-Control-Allow-Headers" = true,
#     "method.response.header.Access-Control-Allow-Methods" = true,
#     "method.response.header.Access-Control-Allow-Origin"  = true
#   }

# }

# resource "aws_api_gateway_integration_response" "proxy" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   resource_id = aws_api_gateway_resource.root.id
#   http_method = aws_api_gateway_method.proxy.http_method
#   status_code = aws_api_gateway_method_response.proxy.status_code


#   //cors
#   response_parameters = {
#     "method.response.header.Access-Control-Allow-Headers" = "'${local.corsAllowHeaders}'",
#     "method.response.header.Access-Control-Allow-Methods" = "'${local.corsAllowMethods}'",
#     "method.response.header.Access-Control-Allow-Origin"  = "'${local.corsAllowOrigins}'"
#   }

#   depends_on = [
#     aws_api_gateway_method.proxy,
#     aws_api_gateway_integration.lambda_integration
#   ]
# }

# //options
# resource "aws_api_gateway_method" "options" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   resource_id = aws_api_gateway_resource.root.id
#   http_method = "OPTIONS"
#   authorization = "NONE"
# #   authorization = "COGNITO_USER_POOLS"
# #   authorizer_id = aws_api_gateway_authorizer.demo.id
# }

# resource "aws_api_gateway_integration" "options_integration" {
#   rest_api_id             = aws_api_gateway_rest_api.api.id
#   resource_id             = aws_api_gateway_resource.root.id
#   http_method             = aws_api_gateway_method.options.http_method
#   integration_http_method = "OPTIONS"
#   type                    = "MOCK"
#   request_templates = {
#     "application/json" = "{\"statusCode\": 200}"
#   }
# }

# resource "aws_api_gateway_method_response" "options_response" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   resource_id = aws_api_gateway_resource.root.id
#   http_method = aws_api_gateway_method.options.http_method
#   status_code = "200"

#   response_parameters = {
#     "method.response.header.Access-Control-Allow-Headers" = true,
#     "method.response.header.Access-Control-Allow-Methods" = true,
#     "method.response.header.Access-Control-Allow-Origin"  = true
#   }
# }

# resource "aws_api_gateway_integration_response" "options_integration_response" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   resource_id = aws_api_gateway_resource.root.id
#   http_method = aws_api_gateway_method.options.http_method
#   status_code = aws_api_gateway_method_response.options_response.status_code

#   response_parameters = {
#     "method.response.header.Access-Control-Allow-Headers" = "'${local.corsAllowHeaders}'",
#     "method.response.header.Access-Control-Allow-Methods" = "'${local.corsAllowMethods}'",
#     "method.response.header.Access-Control-Allow-Origin"  = "'${local.corsAllowOrigins}'"
#   }

#   depends_on = [
#     aws_api_gateway_method.options,
#     aws_api_gateway_integration.options_integration,
#   ]
# }
