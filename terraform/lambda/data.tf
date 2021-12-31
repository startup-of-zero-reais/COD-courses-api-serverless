data "aws_iam_role" "for_lambda" {
  name = format("lambda-%s-role", var.context)
}

data "aws_api_gateway_rest_api" "this" {
  name = format("%s-%s", var.project, var.context)
}

data "aws_api_gateway_resource" "v1" {
  rest_api_id = data.aws_api_gateway_rest_api.this.id
  path        = "/v1"
}
