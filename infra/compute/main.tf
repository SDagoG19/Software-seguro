module "lambda" {
  source                  = "./lambda"
  lambda_bucket           = var.lambda_bucket
  repo_collector_role_arn = var.repo_collector_role_arn

  security_groups_ids = ["sg-0b8feefd748b59858"]
  subnet_ids          = var.subnet_ids

}