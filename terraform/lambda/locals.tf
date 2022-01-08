locals {
  common_tags = {
    Manager   = "Terraform",
    Terraform = "v1.1.2",
    Context   = var.context
    Project   = var.project
  }

  api_resources = {
    lessons  = true
    sections = true
    modules  = true
    courses  = true
  }

  lambdas_resources = {
    read_lesson = {
      path          = "lessons"
      method        = "GET",
      authorization = "NONE"
      description   = "Lambda de leitura de uma aula"
    }

    write_lesson = {
      path          = "lessons"
      method        = "POST",
      authorization = "NONE"
      description   = "Lambda para cadastro de uma aula"
    }

    read_section = {
      path          = "sections"
      method        = "GET",
      authorization = "NONE"
      description   = "Lambda para leitura de uma seção"
    }

    write_section = {
      path          = "sections"
      method        = "POST",
      authorization = "NONE"
      description   = "Lambda para cadastro de uma seção"
    }

    read_module = {
      path          = "modules"
      method        = "GET",
      authorization = "NONE"
      description   = "Lambda para leitura de um módulo"
    }

    write_module = {
      path          = "modules"
      method        = "POST",
      authorization = "NONE"
      description   = "Lambda para cadastro de um módulo"
    }

    read_course = {
      path          = "courses"
      method        = "GET",
      authorization = "NONE"
      description   = "Lambda para leitura de um cursos"
    }

    write_course = {
      path          = "courses"
      method        = "POST",
      authorization = "NONE"
      description   = "Lambda para cadastro de um cursos"
    }
  }
}
