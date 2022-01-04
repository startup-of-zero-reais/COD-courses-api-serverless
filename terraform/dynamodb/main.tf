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

  hash_key  = "PK"
  range_key = "SK"

  dynamic "attribute" {
    for_each = {
      PK           = "S"
      SK           = "S"
      Owner        = "S"
      Title        = "S"
      ParentCourse = "S"
      ParentModule = "S"
    }

    content {
      name = attribute.key
      type = attribute.value
    }
  }

  dynamic "global_secondary_index" {
    for_each = {
      CourseOwnerIndex = {
        hash_key        = "SK"
        range_key       = "Owner"
        projection_type = "ALL"
        read_capacity   = 5
        write_capacity  = 5
      }

      CourseTitleIndex = {
        hash_key        = "Title"
        range_key       = "SK"
        projection_type = "ALL"
        read_capacity   = 5
        write_capacity  = 5
      }

      ModuleLessonsIndex = {
        hash_key        = "ParentModule"
        range_key       = "SK"
        projection_type = "ALL"
        read_capacity   = 5
        write_capacity  = 5
      }

      CourseLessonsIndex = {
        hash_key        = "ParentCourse"
        range_key       = "SK"
        projection_type = "ALL"
        read_capacity   = 5
        write_capacity  = 5
      }
    }

    content {
      name            = global_secondary_index.key
      hash_key        = lookup(global_secondary_index.value, "hash_key")
      range_key       = lookup(global_secondary_index.value, "range_key")
      projection_type = lookup(global_secondary_index.value, "projection_type")
      read_capacity   = lookup(global_secondary_index.value, "read_capacity")
      write_capacity  = lookup(global_secondary_index.value, "write_capacity")
    }
  }

  tags = local.common_tags
}
