package realmd

import (
  "time"

  "github.com/CedricThomas/api-cmangos/cmangos"
  "github.com/CedricThomas/api-cmangos/modules/config"
)

func GetRealmd() cmangos.DaemonInfo {
  var d cmangos.DaemonInfo
  d.Address = config.Settings.Realmd.Address
  d.Port = config.Settings.Realmd.Port
  cmangos.CheckDaemon(
    &d, time.Duration(config.Settings.Api.CheckTimeout))
  return d
}
