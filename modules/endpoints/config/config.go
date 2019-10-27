package config

import (
  "encoding/json"
  "net/http"

  api_config "metagit.org/blizzlike/cmangos-api/cmangos/api/config"
  "metagit.org/blizzlike/cmangos-api/modules/config"
)

func DoConfig(w http.ResponseWriter, r *http.Request) {
  var resp api_config.ApiConfig
  resp.RequireInvite = config.Settings.Api.RequireInvite
  resp.RealmdAddress = config.Settings.Realmd.Address
  resp.RealmdPort = config.Settings.Realmd.Port

  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(resp)
  return
}
