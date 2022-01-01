terraform {
  backend "s3" {
    bucket = "code-craft-backend-supply"
    key    = "dynamodb/courses/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_dynamodb_table" "this" {
  name           = format("%s-%s-table", var.project, var.context)
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5

  hash_key  = "ID"
  range_key = "ParentID"

  dynamic "attribute" {
    for_each = {
      ID        = "S"
      ParentID  = "S"
      CreatedAt = "N"
    }

    content {
      name = attribute.key
      type = attribute.value
    }
  }

  global_secondary_index {
    name               = "ParentIDIndex"
    hash_key           = "ParentID"
    projection_type    = "INCLUDE"
    non_key_attributes = [
      "ID"
    ]
    range_key          = "CreatedAt"
    read_capacity      = 5
    write_capacity     = 5
  }

  tags = local.common_tags
}
