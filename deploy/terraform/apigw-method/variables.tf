variable "rest_api_id" {
    type = string
}

variable "resource_id" {
    type = string
}

variable "http_method" {
    type = string
    default = "GET"
}

variable "authorization" {
    type = string
    default = "NONE"
}

variable "authorizer_id" {
    type = string
}

variable "corsAllowHeaders" {
    type = string
}

variable "corsAllowMethods" {
    type = string
}

variable "corsAllowOrigins" {
    type = string
}

variable "integration_http_method" {
    type = string
    default = "POST"
}

variable "integration_type" {
    type = string
    default = "AWS_PROXY"
}

variable "integration_uri" {
    type = string
}
