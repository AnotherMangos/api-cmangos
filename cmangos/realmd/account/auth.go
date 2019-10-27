package account

import (
  "fmt"

  "metagit.org/blizzlike/cmangos-api/modules/database"
  "metagit.org/blizzlike/cmangos-api/modules/logger"
)

func Authenticate(username, password string) (AccountInfo, error) {
  var a AccountInfo
  stmt, err := database.Realmd.Prepare(
    `SELECT id, username FROM account
     WHERE UPPER(username) = UPPER(?) AND
     sha_pass_hash = SHA1(CONCAT(UPPER(?), ':', UPPER(?)));`)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot prepare query to authenticate %s", username))
    logger.Debug(fmt.Sprintf("%v", err))
    return a, err
  }
  defer stmt.Close()

  err = stmt.QueryRow(username, username, password).Scan(
    &a.Id, &a.Username)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot authenticate account %s", username))
    logger.Debug(fmt.Sprintf("%v", err))
    return a, err
  }

  return a, nil
}
