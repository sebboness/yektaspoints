moved {
  from = aws_api_gateway_method.proxy
  to   = aws_api_gateway_method.root_get
}

moved {
  from = aws_api_gateway_integration.proxy
  to   = aws_api_gateway_integration.root_integration
}

moved {
  from = aws_api_gateway_method_response.proxy
  to   = aws_api_gateway_method_response.root_get
}

moved {
  from = aws_api_gateway_integration_response.integration_response
  to   = aws_api_gateway_integration_response.root_int_resp
}