# ──────────────────────────────────────────────
# Lambda Function
# ──────────────────────────────────────────────

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "${path.module}/../bootstrap"
  output_path = "${path.module}/../lambda.zip"
}

resource "aws_lambda_function" "f1_api" {
  function_name = "${var.project_name}-api"
  description   = "Fetches and caches live F1 race data"

  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  runtime       = "provided.al2023"
  architectures = ["arm64"]
  handler       = "bootstrap"
  role          = aws_iam_role.lambda_exec.arn

  memory_size = var.lambda_memory_size
  timeout     = var.lambda_timeout

  environment {
    variables = {
      DYNAMODB_TABLE    = aws_dynamodb_table.cache.name
      CACHE_TTL_SECONDS = tostring(var.cache_ttl_seconds)
    }
  }

  vpc_config {
    subnet_ids         = [aws_subnet.private_a.id, aws_subnet.private_b.id]
    security_group_ids = [aws_security_group.lambda.id]
  }

  depends_on = [
    aws_iam_role_policy_attachment.lambda_vpc,
    aws_iam_role_policy.lambda_dynamodb,
  ]

  tags = { Name = "${var.project_name}-api" }
}

# CloudWatch Log Group (explicit so Terraform manages retention)
resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${aws_lambda_function.f1_api.function_name}"
  retention_in_days = 14
}
