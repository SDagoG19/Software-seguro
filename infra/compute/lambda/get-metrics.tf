data "aws_s3_object" "get_metrics" {
  bucket = var.lambda_bucket
  key    = "get-metrics.zip"
}

resource "aws_lambda_function" "get_metrics" {
  function_name    = "get-metrics"
  handler          = "bootstrap"
  runtime          = "provided.al2"
  s3_bucket        = var.lambda_bucket
  timeout          = 300
  s3_key           = "get-metrics.zip"
  role             = var.repo_collector_role_arn
  source_code_hash = data.aws_s3_object.get_metrics.version_id

  vpc_config {
    security_group_ids = var.security_groups_ids
    subnet_ids         = var.subnet_ids
  }

  environment { //variables de entorno
     variables = {
      DB_HOST     = "platzi-course.cnyukgcaesvz.us-east-2.rds.amazonaws.com" //endpoint se encuentra en AWS RDS en la bd
      DB_PORT     = "5432" //puerto genrico de postgres
      DB_USER     = "Sdg"
      DB_PASSWORD = "rds!db-c9cdbbdf-9d83-4d56-819b-4ac77c4d5b47"  // nombre del secrets name, se encuentra en secrets manage
      DB_NAME     = "postgres"
    }
  }


}

output "get_metrics_invoke_arn" {
  value = aws_lambda_function.get_metrics.invoke_arn
}

output "get_metrics_lambda_name" {
  value = aws_lambda_function.get_metrics.function_name
}