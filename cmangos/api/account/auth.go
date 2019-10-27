package account

import (
  "fmt"
  "time"
  "database/sql"

  "metagit.org/blizzlike/cmangos-api/modules/config"
  "metagit.org/blizzlike/cmangos-api/modules/database"
  "metagit.org/blizzlike/cmangos-api/modules/logger"
)

func InviteTokenAuth(token string) bool {
  var friend int
  stmt, err := database.Api.Prepare(
    "SELECT friend FROM invitetoken WHERE token = ? AND account IS NULL;")
  if err != nil {
    logger.Error("Cannot query invite token")
    logger.Debug(fmt.Sprintf("%v", err))
    return false
  }
  defer stmt.Close()

  err = stmt.QueryRow(token).Scan(&friend)
  if err != nil {
    logger.Error(fmt.Sprintf("Invalid/Missing invite token %s", token))
    logger.Debug(fmt.Sprintf("%v", err))
    return false
  }

  return true
}

func AuthenticateByToken(token string) (int, error) {
  var owner, expiry int
  stmt, err := database.Api.Prepare(
    "SELECT owner, expiry FROM authtoken WHERE token = ?;")
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot prepare query for token %s", token))
    logger.Debug(fmt.Sprintf("%v", err))
    return 0, err
  }
  defer stmt.Close()

  err = stmt.QueryRow(token).Scan(&owner, &expiry)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot query for auth token %s", token))
    logger.Debug(fmt.Sprintf("%v", err))
    return 0, err
  }

  var stmtUpdate *sql.Stmt
  now := time.Now().Unix()
  if int(now) <= expiry {
    stmtUpdate, err = database.Api.Prepare(
      "UPDATE authtoken SET expiry = ? WHERE token = ?;")
    if err != nil {
      logger.Error(fmt.Sprintf("Cannot prepare update query for auth token %s", token))
      logger.Debug(fmt.Sprintf("%v", err))
      return 0, err
    }
    defer stmtUpdate.Close()

    _, err = stmtUpdate.Exec(now + int64(config.Settings.Api.AuthTokenExpiry), token)
    if err != nil {
      logger.Error(fmt.Sprintf("Cannot update auth token %s", token))
      logger.Debug(fmt.Sprintf("%v", err))
      return 0, err
    }
  } else {
    stmtDelete, err := database.Api.Prepare(
      "DELETE FROM authtoken WHERE token = ?;")
    if err != nil {
      logger.Error(fmt.Sprintf("Cannot prepare delete query for auth token %s", token))
      logger.Debug(fmt.Sprintf("%v", err))
      return 0, err
    }
    defer stmtDelete.Close()

    if err != nil {
      _, err = stmtDelete.Exec(token)
      if err != nil {
        logger.Error(fmt.Sprintf("Cannot delete auth token %s", token))
        logger.Debug(fmt.Sprintf("%v", err))
        return 0, err
      }
    }
  }

  return owner, nil
}

func WriteAuthToken(token string, id int64) error {
  expiry := time.Now().Unix() + int64(config.Settings.Api.AuthTokenExpiry)
  stmt, err := database.Api.Prepare(
    "INSERT INTO authtoken (token, owner, expiry) VALUES (?, ?, ?);")
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot prepare insert query for auth token %s", token))
    logger.Debug(fmt.Sprintf("%v", err))
    return err
  }
  defer stmt.Close()

  _, err = stmt.Exec(token, id, expiry)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot insert auth token %s", token))
    logger.Debug(fmt.Sprintf("%v", err))
    return err
  }

  return nil
}
