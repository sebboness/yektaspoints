resource "aws_dynamodb_table" "points" {
    name = "${local.app}-${local.env}-points"
    billing_mode = "PROVISIONED"
    read_capacity= "10"
    write_capacity= "5"

    attribute {
        name = "id"
        type = "S"
    }

    attribute {
        name = "user_id"
        type = "S"
    }

    hash_key = "id"
    range_key = "user_id"
}

resource "aws_dynamodb_table" "user" {
    name = "${local.app}-${local.env}-user"
    billing_mode = "PROVISIONED"
    read_capacity= "10"
    write_capacity= "5"

    attribute {
        name = "user_id"
        type = "S"
    }

    hash_key = "user_id"
}

resource "aws_dynamodb_table" "family-user" {
    name = "${local.app}-${local.env}-family-user"
    billing_mode = "PROVISIONED"
    read_capacity= "10"
    write_capacity= "5"

    attribute {
        name = "family_id"
        type = "S"
    }

    attribute {
        name = "user_id"
        type = "S"
    }

    hash_key = "family_id"
    range_key = "user_id"
}