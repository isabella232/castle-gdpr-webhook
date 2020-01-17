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
