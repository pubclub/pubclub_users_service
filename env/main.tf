terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.48.0"
    }
  }
  backend "s3" {
    bucket = "pubclub-tf-state"
    key    = "users/terraform.tfstate"
    region = "eu-west-2"
  }
}

provider "aws" {
  region                   = "eu-west-2"
  shared_credentials_files = ["$HOME/.aws/credentials"]
}

module "confirmation_function" {
  source            = "../terraform/lambda_function"
  function_filename = "../builds/confirmation.zip"
  function_name     = "dynamo-confirm-user"
  environment_variables = {
    "TABLE_NAME" : "users-table"
  }
}
