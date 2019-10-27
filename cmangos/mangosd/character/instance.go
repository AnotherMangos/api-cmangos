package character

import (
  "database/sql"
)

type CharacterInstanceInfo struct {
  Db *sql.DB
}
