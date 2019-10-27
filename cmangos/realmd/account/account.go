package account

import (
  "fmt"
  "regexp"
  "database/sql"

  "metagit.org/blizzlike/cmangos-api/modules/database"
  "metagit.org/blizzlike/cmangos-api/modules/logger"
)

const _CMANGOS_MAX_INPUT = 16
const _EMAIL_REGEX = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

type AccountError struct {
  Username bool `json:"username"`
  Password bool `json:"password"`
  Repeat bool `json:"repeat"`
  Email bool `json:"email"`
}

type AccountInfo struct {
  Id int64 `json:"id,omitempty"`
  Username string `json:"username,omitempty"`
  Password string `json:"password,omitempty"`
  Repeat string `json:"repeat,omitempty"`
  Gmlevel int `json:"gmlevel,omitempty"`
  Email string `json:"email,omitempty"`
  Joindate string `json:"joindate,omitempty"`
  Last_ip string `json:"last_ip,omitempty"`
  Failed_logins int `json:"failed_logins,omitempty"`
  Locked int `json:"locked,omitempty"`
  Last_login string `json:"last_login,omitempty"`
  Active_realm_id int `json:"active_realm_id,omitempty"`
  Expansion int `json:"expansion,omitempty"`
  Mutetime int `json:"mutetime,omitempty"`
  Locale int `json:"locale,omitempty"`
}

type RealmCharacters struct {
  Realmid int `json:"realmid,omitempty"`
  Acctid int `json:"acctid,omitempty"`
  Numchars int `json:"numchars,omitempty"`
}

func CreateAccount(account *AccountInfo) (AccountError, error) {
  a := AccountError{true, true, true, true}

  ue, _ := UsernameExists(account.Username)
  if account.Username == "" || len(account.Username) > _CMANGOS_MAX_INPUT || ue {
    a.Username = false
  }
  if account.Password == "" || len(account.Password) > _CMANGOS_MAX_INPUT {
    a.Password = false
  }
  if account.Repeat == "" || account.Password != account.Repeat {
    a.Repeat = false
  }

  ee, _ := EmailExists(account.Email)
  re := regexp.MustCompile(_EMAIL_REGEX)
  if account.Email == "" || !re.MatchString(account.Email) ||
     ee {
    a.Email = false
  }

  if !a.Username || !a.Password || !a.Repeat || !a.Email {
    return a, fmt.Errorf("Cannot create account")
  }

  v, s := CreateSaltAndVerifier(account.Username, account.Password)
  stmt, err := database.Realmd.Prepare(
    `INSERT INTO account
     (username, v, s, email, joindate)
     VALUES (UPPER(?), ?, ?, LOWER(?), NOW());`)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot prepare query to create new account %s", account.Username))
    logger.Debug(fmt.Sprintf("%v", err))
    return a, err
  }
  defer stmt.Close()

  var res sql.Result
  res, err = stmt.Exec(account.Username, v, s, account.Email)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot create new account %s", account.Username))
    logger.Debug(fmt.Sprintf("%v", err))
    return a, err
  }

  account.Id, _ = res.LastInsertId()

  _, err = database.Realmd.Exec(`INSERT INTO realmcharacters (realmid, acctid, numchars) 
    SELECT realmlist.id, account.id, 0 
    FROM realmlist,account LEFT JOIN 
    realmcharacters ON acctid=account.id WHERE acctid IS NULL`)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot create new account %s", account.Username))
    logger.Debug(fmt.Sprintf("%v", err))
    return a, err
  }

  return a, nil
}

func EmailExists(email string) (bool, error) {
  var id int
  stmt, err := database.Realmd.Prepare(
    "SELECT id FROM account WHERE email = ?;")
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot prepare query to get email %s", email))
    logger.Debug(fmt.Sprintf("%v", err))
    return false, err
  }
  defer stmt.Close()

  err = stmt.QueryRow(email).Scan(&id)
  if err != nil {
    logger.Debug(fmt.Sprintf("%v", err))
    return false, nil
  }

  return true, nil
}

func GetAccountInfo(id int) (AccountInfo, error) {
  var ai AccountInfo
  stmt, err := database.Realmd.Prepare(
    `SELECT id, username, gmlevel, email, joindate, last_ip,
     failed_logins, locked, last_login, active_realm_id,
     expansion, mutetime, locale
     FROM account
     WHERE id = ?;`)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot prepare query to get acount info %d", id))
    logger.Debug(fmt.Sprintf("%v", err))
    return ai, err
  }
  defer stmt.Close()

  err = stmt.QueryRow(id).Scan(
    &ai.Id, &ai.Username, &ai.Gmlevel, &ai.Email, &ai.Joindate,
    &ai.Last_ip, &ai.Failed_logins, &ai.Locked, &ai.Last_login,
    &ai.Active_realm_id, &ai.Expansion, &ai.Mutetime, &ai.Locale)
  if err != nil {
    logger.Debug(fmt.Sprintf("%v", err))
    return ai, err
  }

  return ai, nil
}

func GetRealmCharacters(id int) ([]RealmCharacters, error) {
  var realmcharacters []RealmCharacters
  stmt, err := database.Realmd.Prepare(
    "SELECT realmid, acctid, numchars FROM realmcharacters WHERE acctid = ?;")
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot prepare query to fetch characters of account %d", id))
    logger.Debug(fmt.Sprintf("%v", err))
    return realmcharacters, err
  }
  defer stmt.Close()

  rows, _ := stmt.Query(id)
  for rows.Next() {
    var rc RealmCharacters
    err = rows.Scan(&rc.Realmid, &rc.Acctid, &rc.Numchars)
    if err != nil {
      logger.Error(fmt.Sprintf("Cannot query characters of account %d", id))
      logger.Debug(fmt.Sprintf("%v", err))
      return realmcharacters, err
    }

    realmcharacters = append(realmcharacters, rc)
  }

  return realmcharacters, nil
}

func UsernameExists(username string) (bool, error) {
  var id int
  stmt, err := database.Realmd.Prepare(
    "SELECT id FROM account WHERE username = ?;")
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot preapre query to find account %s", username))
    logger.Debug(fmt.Sprintf("%v", err))
    return false, err
  }
  defer stmt.Close()

  err = stmt.QueryRow(username).Scan(&id)
  if err != nil {
    logger.Debug(fmt.Sprintf("%v", err))
    return false, nil
  }

  return true, nil
}
