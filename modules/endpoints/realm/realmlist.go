package realm

import (
  "encoding/json"
  "net/http"

  cmangos_realm "metagit.org/blizzlike/cmangos-api/cmangos/realmd/realm"
)

func DoRealmlist(w http.ResponseWriter, r *http.Request) {
  realmlist := cmangos_realm.GetRealms()

  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(realmlist)
  return
}
