data "aws_secretsmanager_secret" "secrets" {
    name = "hexonite/${var.app}/secrets"
}

data "aws_secretsmanager_secret_version" "current" {
    secret_id = data.aws_secretsmanager_secret.secrets.id
}

locals {
    ssm_secrets = jsondecode(
        data.aws_secretsmanager_secret_version.current.secret_string
    )
}