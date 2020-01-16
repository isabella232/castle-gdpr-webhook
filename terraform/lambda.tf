provider "aws" {
  profile    = "default"
  region     = "us-west-2"
}

variable "lambda_function_name" {
  description = "The name of the lambda function"
  default = "CastleHandler"
}

resource "aws_lambda_function" "example" {
	#function_name = "CastleHandler"
  function_name = "${var.lambda_function_name}"

	filename="../function.zip"

  handler = "castle-gdpr-webhook"
  runtime = "go1.x"

  role = aws_iam_role.iam_for_lambda.arn

	environment {
    variables = {
			S3BUCKET = "castle-gdpr-user-data"
			# these are for testing only
      HMACSECRET = "ssshhh..."
    }
  }
	## OLA TODO give permissions for S3
	depends_on = ["aws_iam_role_policy_attachment.lambda_logs", "aws_cloudwatch_log_group.example"]
	#depends_on = [aws_cloudwatch_log_group.example]
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.example.function_name
  principal     = "apigateway.amazonaws.com"

  # The "/*/*" portion grants access from any method on any resource
  # within the API Gateway REST API.
  source_arn = "${aws_api_gateway_rest_api.example.execution_arn}/*/*"
}

# IAM role which dictates what other AWS services the Lambda function
# may access.
resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

# This is to optionally manage the CloudWatch Log Group for the Lambda Function.
# If skipping this resource configuration, also add "logs:CreateLogGroup" to the IAM policy below.
resource "aws_cloudwatch_log_group" "example" {
  name              = "/aws/lambda/CastleHandler"
  retention_in_days = 14
}

# See also the following AWS managed policy: AWSLambdaBasicExecutionRole
resource "aws_iam_policy" "lambda_logging" {
  name = "lambda_logging"
  path = "/"
  description = "IAM policy for logging from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
	#role = "aws_iam_role.iam_for_lambda.name"
	role = "${aws_iam_role.iam_for_lambda.name}"
	policy_arn = "${aws_iam_policy.lambda_logging.arn}"

	#role = "serverless_example_lambda"
	#policy_arn = "aws_iam_policy.lambda_exec.arn"
}

