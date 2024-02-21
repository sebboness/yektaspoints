terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.30"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.3.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.2.0"
    }
  }

  required_version = "~> 1.7.3"

  backend "s3" {
    bucket = "hexonite-mypoints-infra"
    key    = "terraform"
    region = "us-west-2"
  }
}

locals {
  app = "mypoints"
  env = terraform.workspace == "default" ? "dev" : terraform.workspace

  corsAllowHeaders = "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token" 
  corsAllowMethods = "DELETE,GET,OPTIONS,PATCH,POST,PUT"
  corsAllowOrigins = "mypoints.hexonite.net,dev-mypoints.hexonite.net,staging-mypoints.hexonite.net"

  apiSubdomain = local.env == "prod" ? "mypoints-api" : "${local.env}-mypoints-api" # (i.e. mypoints-api.hexonite.net)
}

provider "aws" {
  region = "us-west-2"
  default_tags {
    tags = {
        app = local.app
        env = local.env
    }
  }
}