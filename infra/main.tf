terraform {
  backend "s3" {
    bucket = "sdg-terraform-bucket"
    key    = "terraform.tfstate"
    region = "us-east-2"
  }
}

provider "aws" {
  region = "us-east-2"
}

module "iam" {
  source = "./iam"
}

module "s3" {
  source = "./s3"
}

module "compute" {
  source                  = "./compute"
  lambda_bucket           = module.s3.lambda_bucket
  repo_collector_role_arn = module.iam.repo_collector_role_arn

  subnet_ids = ["subnet-02a7da13479223f25", "subnet-0648abf4e63f0be95", "subnet-049c7df8fad6419d4"]


}