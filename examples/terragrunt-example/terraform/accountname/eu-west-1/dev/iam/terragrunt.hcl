locals {
  # Automatically load environment-level variables
  global_vars = read_terragrunt_config(find_in_parent_folders("global.hcl"))
  environment_vars = read_terragrunt_config(find_in_parent_folders("env.hcl"))
  account_vars = read_terragrunt_config(find_in_parent_folders("account.hcl"))
}

terraform {
    source = "git@github.com:terraform-aws-modules/terraform-aws-iam.git//modules/iam-assumable-roles?ref=v2.21.0"
}

include {
  path = find_in_parent_folders()
}