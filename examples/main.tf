module "consul" {
  source = "hashicorp/consul/aws"
  version = "> 0.1.0"
}

module "consul_github_https" {
  source = "github.com/hashicorp/terraform-aws-consul"
  version = "0.7.3"
}

module "consul_github_ssh" {
  source = "git@github.com:hashicorp/terraform-aws-consul?ref=0.1.0"
  version = "0.1.0"
}

module "example_git_ssh" {
  source = "git::ssh://git@github.com/sgreben/terraform-module-versions?ref=0.10.0"
  version = "~> 0.10"
}
