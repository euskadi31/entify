
table "users" {
  schema = schema.demo

  column "id" {
    null = false
    type = char(36)
  }

  column "email" {
    null = false
    type = varchar(90)
  }

  column "firstname" {
    null = true
    type = varchar(45)
  }

  column "lastname" {
    null = true
    type = varchar(45)
  }

  column "password" {
    null = false
    type = varchar(128)
  }

  column "salt" {
    null = false
    type = varchar(40)
  }

  column "enabled" {
    null    = false
    type    = bool
    default = 0
  }

  column "expired" {
    null    = false
    type    = bool
    default = 0
  }

  column "locked" {
    null    = false
    type    = bool
    default = 0
  }

  column "timezone" {
    null = true
    type = varchar(40)
  }

  column "locale" {
    null = true
    type = varchar(5)
  }

  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null = true
    type = timestamp
  }

  column "deleted_at" {
    null = true
    type = timestamp
  }

  primary_key {
    columns = [column.id]
  }

  index "users_email_index" {
    columns = [column.email]
    type    = BTREE
  }
}

table "users_activations" {
  schema = schema.demo

  column "user_id" {
    null = false
    type = char(36)
  }

  column "code" {
    null = false
    type = char(15)
  }

  column "status" {
    null     = false
    type     = tinyint
    default  = 0
    unsigned = true
    comment  = "is used"
  }

  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null = true
    type = timestamp
  }

  primary_key {
    columns = [column.user_id, column.code]
  }

  foreign_key "users_activations_ibfk_1" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = "RESTRICT"
    on_delete   = "RESTRICT"
  }
}

schema "demo" {
  charset   = "utf8mb4"
  collation = "utf8mb4_general_ci"
}
