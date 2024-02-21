data "aws_route53_zone" "my_domain" {
  name         = "hexonite.net"
  private_zone = false
}

resource "aws_route53_record" "custom_domain_record" {
  name = local.apiSubdomain
  type = "CNAME"
  ttl  = "300" # TTL in seconds

  records = ["${aws_api_gateway_rest_api.api.id}.execute-api.us-west-2.amazonaws.com"]

  zone_id = data.aws_route53_zone.my_domain.zone_id
}