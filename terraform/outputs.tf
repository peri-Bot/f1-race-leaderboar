output "api_gateway_url" {
  description = "Base URL of the API Gateway endpoint"
  value       = aws_apigatewayv2_stage.default.invoke_url
}

output "lambda_function_name" {
  description = "Name of the deployed Lambda function"
  value       = aws_lambda_function.f1_api.function_name
}

output "dynamodb_table_name" {
  description = "Name of the DynamoDB cache table"
  value       = aws_dynamodb_table.cache.name
}

output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}
