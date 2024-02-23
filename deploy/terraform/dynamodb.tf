resource "aws_dynamodb_table" "points" {
    name = "${local.env}-${local.env}-points"
    billing_mode = "PROVISIONED"
    read_capacity= "10"
    write_capacity= "10"
    attribute {
        name = "noteId"
        type = "S"
    }
    hash_key = "noteId"
}