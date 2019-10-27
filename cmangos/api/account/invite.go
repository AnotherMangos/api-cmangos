package account

import (
  "fmt"
  "database/sql"

  "github.com/google/uuid"

  "metagit.org/blizzlike/cmangos-api/modules/database"
  "metagit.org/blizzlike/cmangos-api/modules/logger"
)

type InviteInfo struct {
  Token string `json:"token,omitempty"`
}

func AddAccountToInviteToken(token string, id int64) error {
  stmt, err := database.Api.Prepare(
    "UPDATE invitetoken SET account = ? WHERE token = ?;")
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot prepare query to update invite token owner %s", token))
    logger.Debug(fmt.Sprintf("%v", err))
    return err
  }
  defer stmt.Close()

  _, err = stmt.Exec(id, token)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot update invite token owner %s", token))
    logger.Debug(fmt.Sprintf("%v", err))
    return err
  }

  return nil
}

func GetInviteTokens(id int) ([]InviteInfo, error) {
  var ii []InviteInfo = []InviteInfo{}
  stmt, err := database.Api.Prepare(
    "SELECT token FROM invitetoken WHERE account IS NULL AND friend = ?;")
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot prepare query to fetch invite tokens for account %d", id))
    logger.Debug(fmt.Sprintf("%v", err))
    return ii, err
  }
  defer stmt.Close()

  var rows *sql.Rows
  rows, err = stmt.Query(id)
  for rows.Next() {
    var t InviteInfo
    err = rows.Scan(&t.Token)
    if err != nil {
      logger.Error(fmt.Sprintf("Cannot query invite tokens for account %d", id))
      logger.Debug(fmt.Sprintf("%v", err))
      return ii, err
    }

    ii = append(ii, t)
  }

  return ii, nil
}

func WriteInviteToken(id int) (string, error) {
  t, _ := uuid.NewRandom()
  token := t.String()
  stmt, err := database.Api.Prepare(
    "INSERT INTO invitetoken (token, friend) VALUES (?, ?);")
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot prepare query to insert invite token for account %d", id))
    logger.Debug(fmt.Sprintf("%v", err))
    return token, err
  }
  defer stmt.Close()

  _, err = stmt.Exec(token, id)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot insert invite token for account %d", id))
    logger.Debug(fmt.Sprintf("%v", err))
    return token, err
  }

  return token, nil
}
