variable "project" {
  type        = string
  description = "Nome do projeto"
  default     = "code-craft"
}

variable "context" {
  type        = string
  description = "Contexto da aplicação"
}

variable "aws_region" {
  type        = string
  description = "Região da AWS"
  default     = "us-east-1"
}

variable "environment" {
  type        = string
  description = "Ambiente da aplicação"
  default     = "dev"
}

## API Gateway

variable "api_key_source" {
  type        = string
  description = "Contexto da aplicação"
  default     = null
}

variable "stage" {
  type        = string
  description = "Stage do deploy"
  default     = "dev"
}

## Route53

variable "domain_name" {
  type        = string
  description = "Nome do domínio"
  default     = "code-craft.club"
}


## Lambdas

variable "runtime" {
  type        = string
  description = "Runtime da aplicação"
}
