
variable "v" {
  description = "in root module"
  default     = ""
}

module "child_a" {
  source = "./child_a"
}
