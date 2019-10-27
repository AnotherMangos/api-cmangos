package realmd

import (
  "time"

  "metagit.org/blizzlike/cmangos-api/cmangos"
  "metagit.org/blizzlike/cmangos-api/modules/config"
)

func GetRealmd() cmangos.DaemonInfo {
  var d cmangos.DaemonInfo
  d.Address = config.Settings.Realmd.Address
  d.Port = config.Settings.Realmd.Port
  cmangos.CheckDaemon(
    &d, time.Duration(config.Settings.Api.CheckTimeout))
  return d
}
