package account

import (
  "encoding/base64"
  "net/http"
  "fmt"
  "strings"

  "github.com/google/uuid"

  api_account "metagit.org/blizzlike/cmangos-api/cmangos/api/account"
  cmangos_account "metagit.org/blizzlike/cmangos-api/cmangos/realmd/account"
  "metagit.org/blizzlike/cmangos-api/modules/logger"
)

func Authenticate(w http.ResponseWriter, r *http.Request) (cmangos_account.AccountInfo, error) {
  header := r.Header.Get("Authorization")
  auth := strings.Split(header, " ")
  var a cmangos_account.AccountInfo

  if len(auth) != 2 {
    w.WriteHeader(http.StatusBadRequest)
    return a, fmt.Errorf("Invalid/Missing Authorization header")
  }

  if !strings.EqualFold(auth[0], "basic") {
    errmsg := fmt.Sprintf("Authentication method not supported (%s)", auth[0])
    logger.Error(errmsg)
    w.WriteHeader(http.StatusBadRequest)
    return a, fmt.Errorf(errmsg)
  }

  credentials, err := base64.StdEncoding.DecodeString(auth[1])
  c := strings.Split(string(credentials), ":")
  a, err = cmangos_account.Authenticate(c[0], c[1])
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot authenticate %s", c[0]))
    logger.Debug(fmt.Sprintf("%v", err))
    w.WriteHeader(http.StatusUnauthorized)
    return a, err
  }

  return a, nil
}

func DoAuth(w http.ResponseWriter, r *http.Request) {
  a, err := Authenticate(w, r)
  if err != nil {
    logger.Debug(fmt.Sprintf("%v", err))
    return
  }

  logger.Info(fmt.Sprintf("Authenticated %s", a.Username))

  t, err := uuid.NewRandom()
  token := t.String()
  err = api_account.WriteAuthToken(token, a.Id)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot write auth token %s", token))
    logger.Debug(fmt.Sprintf("%v", err))
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  w.Header().Add("X-Auth-Token", token)
  w.WriteHeader(http.StatusOK)
  return
}

func DoAuthVerify(w http.ResponseWriter, r *http.Request) {
  id, err := AuthenticateByToken(w, r)
  if err != nil {
    logger.Error("Cannot verify authentication")
    return
  }

  logger.Info(fmt.Sprintf("AUthenticated %d", id))
  w.WriteHeader(http.StatusOK)
  return
}
