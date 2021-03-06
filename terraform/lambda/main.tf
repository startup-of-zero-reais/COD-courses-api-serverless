terraform {
  backend "s3" {
    bucket = "code-craft-backend-supply"
    key    = "lambdas/courses-api/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_lambda_function" "this" {
  for_each = local.lambdas_resources

  filename      = "${path.module}/functions/${each.key}.zip"
  function_name = each.key
  role          = data.aws_iam_role.for_lambda.arn
  handler       = each.key
  runtime       = var.runtime
  timeout       = "10"
  memory_size   = "128"

  description = format("%s: %s", each.key, each.value.description)

  source_code_hash = filebase64sha256("${path.module}/functions/${each.key}.zip")

  tags = merge(local.common_tags, {
    description = format("%s: %s", each.key, each.value.description)
  })
}

resource "aws_lambda_permission" "this" {
  for_each = local.lambdas_resources

  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this[each.key].function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = format(
  "%s/*/%s%s",
  data.aws_api_gateway_rest_api.this.execution_arn,
  aws_api_gateway_method.this[each.key].http_method,
  aws_api_gateway_resource.this[each.value.path].path
  )

  depends_on = [
    aws_api_gateway_method.this,
    aws_api_gateway_resource.this,
  ]
}

# API Gateway resource

resource "aws_api_gateway_resource" "this" {
  for_each = local.api_resources

  rest_api_id = data.aws_api_gateway_rest_api.this.id
  parent_id   = data.aws_api_gateway_resource.v1.id
  path_part   = each.key
}

resource "aws_api_gateway_method" "this" {
  for_each = local.lambdas_resources

  rest_api_id      = data.aws_api_gateway_rest_api.this.id
  resource_id      = aws_api_gateway_resource.this[each.value.path].id
  http_method      = each.value.method
  authorization    = each.value.authorization
  api_key_required = true
}

resource "aws_api_gateway_method_response" "options_response" {
  for_each = local.lambdas_resources

  rest_api_id         = data.aws_api_gateway_rest_api.this.id
  resource_id         = aws_api_gateway_resource.this[each.value.path].id
  http_method         = aws_api_gateway_method.this[each.key].http_method
  status_code         = 200
  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = true
  }

  depends_on = [aws_api_gateway_method.this]
}

resource "aws_api_gateway_method" "options" {
  for_each = local.api_resources

  rest_api_id   = data.aws_api_gateway_rest_api.this.id
  resource_id   = aws_api_gateway_resource.this[each.key].id
  authorization = "NONE"
  http_method   = "OPTIONS"
}

resource "aws_api_gateway_method_response" "options_200" {
  for_each = local.api_resources

  rest_api_id = data.aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_resource.this[each.key].id
  http_method = aws_api_gateway_method.options[each.key].http_method
  status_code = 200

  response_models = {
    "application/json" = "Empty"
  }

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true,
    "method.response.header.Access-Control-Allow-Methods" = true,
    "method.response.header.Access-Control-Allow-Origin"  = true
    "method.response.header.user"                         = true
  }

  depends_on = [aws_api_gateway_method.options]
}

resource "aws_api_gateway_integration" "options_integration" {
  for_each = local.api_resources

  rest_api_id = data.aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_resource.this[each.key].id
  http_method = aws_api_gateway_method.options[each.key].http_method
  type        = "MOCK"
  depends_on  = [aws_api_gateway_method.options]

  request_templates = {
    "application/json" = jsonencode({
      statusCode = 200
    })
  }
}

resource "aws_api_gateway_integration_response" "options_integration_response" {
  for_each = local.api_resources

  rest_api_id = data.aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_resource.this[each.key].id
  http_method = aws_api_gateway_method.options[each.key].http_method
  status_code = aws_api_gateway_method_response.options_200[each.key].status_code

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,user'",
    "method.response.header.Access-Control-Allow-Methods" = "'GET,OPTIONS,POST,PUT'",
    "method.response.header.Access-Control-Allow-Origin"  = "'*'"
    "method.response.header.user"                         = "integration.response.header.user"
  }

  depends_on = [aws_api_gateway_method_response.options_200]
}

resource "aws_api_gateway_integration" "this" {
  for_each = local.lambdas_resources

  rest_api_id             = data.aws_api_gateway_rest_api.this.id
  resource_id             = aws_api_gateway_resource.this[each.value.path].id
  http_method             = lookup(local.lambdas_resources, each.key, { method = "GET" }).method
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
  uri                     = aws_lambda_function.this[each.key].invoke_arn
}

resource "aws_api_gateway_deployment" "this" {
  for_each = local.lambdas_resources

  rest_api_id = data.aws_api_gateway_rest_api.this.id
  stage_name  = var.stage

  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.this[each.value.path].id,
      aws_api_gateway_method.this[each.key].id,
      aws_api_gateway_integration.this[each.key].id
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}
