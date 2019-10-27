package realm

import (
  "fmt"
  "time"
  "database/sql"

  "metagit.org/blizzlike/cmangos-api/cmangos"
  "metagit.org/blizzlike/cmangos-api/modules/database"
  "metagit.org/blizzlike/cmangos-api/modules/config"
  "metagit.org/blizzlike/cmangos-api/modules/logger"
)

var realmlist []Realm

func FetchRealms() ([]Realm, error) {
  var rl []Realm
  stmt, err := database.Realmd.Prepare(
    `SELECT id, name, address, port, icon, population
     FROM realmlist
     ORDER BY id ASC;`)
  if err != nil {
    logger.Error("Cannot prepare query to fetch realms")
    logger.Debug(fmt.Sprintf("%v", err))
    return rl, err
  }
  defer stmt.Close()

  var rows *sql.Rows
  rows, err = stmt.Query()
  for rows.Next() {
    var realm Realm
    err = rows.Scan(&realm.Id, &realm.Name, &realm.Host.Address,
      &realm.Host.Port, &realm.Icon, &realm.Population)
    if err != nil {
      logger.Error("Cannot query realms")
      logger.Debug(fmt.Sprintf("%v", err))
      return rl, err
    }

    for k, v := range database.Mangosd {
      if v.Id == realm.Id {
        realm.CharacterInstance.Db = database.Mangosd[k].Character
      }
    }

    cmangos.CheckDaemon(&realm.Host, time.Duration(config.Settings.Api.CheckTimeout))
    rl = append(rl, realm)
  }

  return rl, nil
}

func GetRealms() []Realm {
  return realmlist
}

func PollRealmStates(interval time.Duration) {
  realmlist, _ = FetchRealms()

  t := time.Duration(time.Duration(interval) * time.Second)
  for range time.Tick(t) {
    rl, err := FetchRealms()
    if err != nil {
      logger.Error("Cannot fetch realmlist")
      logger.Debug(fmt.Sprintf("%v", err))
      continue
    }
    logger.Info("Fetched realmlist")

    realmlist = rl
  }
}
