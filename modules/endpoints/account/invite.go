package account

import (
  "encoding/json"
  "net/http"
  "fmt"
  "strings"

  api_account "metagit.org/blizzlike/cmangos-api/cmangos/api/account"
  "metagit.org/blizzlike/cmangos-api/modules/logger"
)

func AuthenticateByInviteToken(w http.ResponseWriter, r *http.Request) (string, error) {
  header := r.Header.Get("Authorization")
  auth := strings.Split(header, " ")

  if len(auth) != 2 {
    errmsg := "Invalid/Missing Authorization header"
    logger.Error(errmsg)
    w.WriteHeader(http.StatusBadRequest)
    return "", fmt.Errorf(errmsg)
  }

  if !strings.EqualFold(auth[0], "token") {
    errmsg := "Authentication method not supported"
    logger.Error(errmsg)
    w.WriteHeader(http.StatusBadRequest)
    return "", fmt.Errorf(errmsg)
  }

  if !api_account.InviteTokenAuth(auth[1]) {
    w.WriteHeader(http.StatusUnauthorized)
    return "", fmt.Errorf("Cannot authenticate invite %s", auth[1])
  }

  return auth[1], nil
}

func DoGetInvites(w http.ResponseWriter, r *http.Request) {
  id, err := AuthenticateByToken(w, r)
  if err != nil {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  tokens, err := api_account.GetInviteTokens(id)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(tokens)
  return
}

func DoInvite(w http.ResponseWriter, r *http.Request) {
  id, err := AuthenticateByToken(w, r)
  if err != nil {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  logger.Info(fmt.Sprintf("Authenticated id %d", id))
  token, err := api_account.WriteInviteToken(id)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  var inv = api_account.InviteInfo{token}
  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(inv)
  return
}
