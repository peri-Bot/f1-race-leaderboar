# ──────────────────────────────────────────────
# API Gateway (HTTP API v2)
# ──────────────────────────────────────────────

resource "aws_apigatewayv2_api" "http_api" {
  name          = "${var.project_name}-api"
  protocol_type = "HTTP"
  description   = "F1 Race Leaderboard API"

  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["GET", "OPTIONS"]
    allow_headers = ["Content-Type"]
    max_age       = 3600
  }
}

# Lambda integration
resource "aws_apigatewayv2_integration" "lambda" {
  api_id             = aws_apigatewayv2_api.http_api.id
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.f1_api.invoke_arn

  payload_format_version = "2.0"
}

# Routes
resource "aws_apigatewayv2_route" "root" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_route" "standings" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /standings"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_route" "standings_drivers" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /standings/drivers"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_route" "standings_constructors" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /standings/constructors"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_route" "results" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /results"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_route" "schedule" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /schedule"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

# Default stage (auto-deploy)
resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.http_api.id
  name        = "$default"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gw.arn
    format = jsonencode({
      requestId      = "$context.requestId"
      ip             = "$context.identity.sourceIp"
      requestTime    = "$context.requestTime"
      httpMethod     = "$context.httpMethod"
      routeKey       = "$context.routeKey"
      status         = "$context.status"
      protocol       = "$context.protocol"
      responseLength = "$context.responseLength"
      integrationLatency = "$context.integrationLatency"
    })
  }
}

resource "aws_cloudwatch_log_group" "api_gw" {
  name              = "/aws/apigateway/${var.project_name}"
  retention_in_days = 14
}

# Permission for API Gateway to invoke Lambda
resource "aws_lambda_permission" "api_gw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.f1_api.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}
