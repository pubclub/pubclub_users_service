variable "function_filename" {
  type        = string
  description = "Path to zip file containing lambda function"
}

variable "function_name" {
  type        = string
  description = "Name of lambda function"
}

variable "environment_variables" {
  type        = map(string)
  description = "Key-values for any environment variables to use in function"
}
