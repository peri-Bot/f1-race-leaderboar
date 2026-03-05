# ──────────────────────────────────────────────
# Terraform — F1 Race Leaderboard
# ──────────────────────────────────────────────

terraform {
  required_version = ">= 1.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.0"
    }
  }

  # Uncomment for remote state in production:
  # backend "s3" {
  #   bucket = "f1-leaderboard-tfstate"
  #   key    = "terraform.tfstate"
  #   region = "eu-west-1"
  # }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}

data "aws_region" "current" {}
data "aws_caller_identity" "current" {}
