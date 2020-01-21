provider "aws" {
  profile = "default"
  region  = "us-west-2"
}

variable "lambda_function_name" {
  description = "The name of the lambda function"
  default     = "CastleHandler"
}

variable "s3bucket" {
  description = "Where gdpr data files are kept"
  default     = "castle-gdpr-user-data"
}

variable "hmac_secret" {
  description = "The HMAC secret used to validate calls"
}

variable "cert_arn" {
	# For now the ACM is manually managed, the ARN is from it
	# The DNS is in Akamai so it cannot be updated from Terraform, the domain would need to move to Route53
	# Or another domain which is managed in Route53 should be used
	description = "This is the ACM certificate for the public domain castlewebhook-test.optimizely.com"
  default = "arn:aws:acm:us-west-2:987056895854:certificate/ba9c30ca-24b3-42e3-bbb5-52fee27df02e"
}
