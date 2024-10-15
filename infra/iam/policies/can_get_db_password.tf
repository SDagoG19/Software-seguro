resource "aws_iam_policy" "can_get_db_password" {
  name        = "can-get-db-password"
  path        = "/"
  description = "Allow access to retrieve secrets from Secrets Manager"

  policy = jsonencode({
    Version : "2012-10-17",
    Statement : [
      {
        Effect : "Allow",
        Action : [
          "secretsmanager:GetSecretValue"
        ],
        Resource : [
          "arn:aws:secretsmanager:us-east-2:390403852799:secret:rds!db-c9cdbbdf-9d83-4d56-819b-4ac77c4d5b47-sUabkb"
        ]
      }
    ]
  })
}

output "can_get_db_password_arn" {
  value = aws_iam_policy.can_get_db_password.arn
}