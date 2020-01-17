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

# this prints the base URL after terraform apply to simplify testing
output "base_url" {
  value = aws_api_gateway_deployment.castle_gdpr_webhook.invoke_url
}
