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
}

variable "env" {
    default = "local"
}

variable "app" {
    default = "mypoints"
}

provider "aws" {
  region = "us-west-2"
  default_tags {
    tags = {
        app = "mypoints"
        env = var.env
    }
  }
}