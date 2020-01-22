resource "aws_api_gateway_rest_api" "castle_gdpr_webhook" {
  name        = "Castle GDPR Webhook"
  description = "Webhook that is invoked by Castle.io to download GDPR data"
}

resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = aws_api_gateway_rest_api.castle_gdpr_webhook.id
  parent_id   = aws_api_gateway_rest_api.castle_gdpr_webhook.root_resource_id
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_method" "proxy" {
  rest_api_id   = aws_api_gateway_rest_api.castle_gdpr_webhook.id
  resource_id   = aws_api_gateway_resource.proxy.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda" {
  rest_api_id = aws_api_gateway_rest_api.castle_gdpr_webhook.id
  resource_id = aws_api_gateway_method.proxy.resource_id
  http_method = aws_api_gateway_method.proxy.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.castle_webhook.invoke_arn
}

resource "aws_api_gateway_method" "proxy_root" {
  rest_api_id   = aws_api_gateway_rest_api.castle_gdpr_webhook.id
  resource_id   = aws_api_gateway_rest_api.castle_gdpr_webhook.root_resource_id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda_root" {
  rest_api_id = aws_api_gateway_rest_api.castle_gdpr_webhook.id
  resource_id = aws_api_gateway_method.proxy_root.resource_id
  http_method = aws_api_gateway_method.proxy_root.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.castle_webhook.invoke_arn
}

resource "aws_api_gateway_deployment" "castle_gdpr_webhook" {
  depends_on = [
    aws_api_gateway_integration.lambda,
    aws_api_gateway_integration.lambda_root,
  ]

  rest_api_id = aws_api_gateway_rest_api.castle_gdpr_webhook.id
  stage_name  = "v1"
}

# ensure the path mapping is setup, e.g. "v1"
resource "aws_api_gateway_base_path_mapping" "test" {
  api_id      = "${aws_api_gateway_rest_api.castle_gdpr_webhook.id}"
  stage_name  = "${aws_api_gateway_deployment.castle_gdpr_webhook.stage_name}"
  domain_name = "${aws_api_gateway_domain_name.castlewebhook_test.domain_name}"
  base_path   = "v1"
}

# setup the custom domain name
resource "aws_api_gateway_domain_name" "castlewebhook_test" {
  domain_name              = "castlewebhook-test.optimizely.com"
  security_policy          = "TLS_1_2"
  regional_certificate_arn = "${var.cert_arn}"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

# this prints the base URL after terraform apply to simplify testing
output "base_url" {
  value = aws_api_gateway_deployment.castle_gdpr_webhook.invoke_url
}
