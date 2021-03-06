package realm

import (
  "encoding/json"
  "net/http"

  cmangos_realm "github.com/CedricThomas/api-cmangos/cmangos/realmd/realm"
)

func DoRealmlist(w http.ResponseWriter, r *http.Request) {
  realmlist := cmangos_realm.GetRealms()

  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(realmlist)
  return
}
