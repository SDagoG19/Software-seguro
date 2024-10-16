data "aws_s3_object" "handle_github_webhook" {
  bucket = var.lambda_bucket
  key    = "handle-github-webhook.zip"
}

resource "aws_lambda_function" "handle_github_webhook" {
  function_name    = "handle-github-webhook"
  handler          = "bootstrap"
  runtime          = "provided.al2"
  s3_bucket        = var.lambda_bucket
  timeout          = 300
  s3_key           = "handle-github-webhook.zip"
  role             = var.repo_collector_role_arn
  source_code_hash = data.aws_s3_object.handle_github_webhook.version_id

  vpc_config {
    security_group_ids = var.security_groups_ids
    subnet_ids         = var.subnet_ids
  }

  environment {
    variables = {
      DB_HOST     = "platzi-course.cnyukgcaesvz.us-east-2.rds.amazonaws.com" //endpoint se encuentra en AWS RDS en la bd
      DB_PORT     = "5432" //puerto genrico de postgres
      DB_USER     = "Sdg"
      DB_PASSWORD = "rds!db-c9cdbbdf-9d83-4d56-819b-4ac77c4d5b47"  // nombre del secrets name, se encuentra en secrets manage
      DB_NAME     = "postgres"

      GITHUB_SECRET = "github/secret" //secret name en el secret manager
    }
  }
}

output "handle_github_webhook_invoke_arn" {
  value = aws_lambda_function.handle_github_webhook.invoke_arn
}

output "handle_github_webhook_lambda_name" {
  value = aws_lambda_function.handle_github_webhook.function_name
}