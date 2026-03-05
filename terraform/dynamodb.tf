# ──────────────────────────────────────────────
# DynamoDB — Cache table + VPC Gateway Endpoint
# ──────────────────────────────────────────────

resource "aws_dynamodb_table" "cache" {
  name         = "${var.project_name}-cache"
  billing_mode = "PAY_PER_REQUEST" # on-demand, scales to zero cost

  hash_key = "cache_key"

  attribute {
    name = "cache_key"
    type = "S"
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }

  tags = { Name = "${var.project_name}-cache" }
}

# VPC Gateway Endpoint for DynamoDB — free, keeps traffic off the internet
resource "aws_vpc_endpoint" "dynamodb" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${data.aws_region.current.name}.dynamodb"
  vpc_endpoint_type = "Gateway"

  route_table_ids = [
    aws_route_table.private.id,
  ]

  tags = { Name = "${var.project_name}-dynamodb-endpoint" }
}
