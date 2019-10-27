package database

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"

  "metagit.org/blizzlike/cmangos-api/modules/logger"
  "metagit.org/blizzlike/cmangos-api/modules/config"
)

type MangosdDB struct {
  Id int
  Character *sql.DB
  World *sql.DB
}

var Api *sql.DB
var Realmd *sql.DB
var Mangosd []MangosdDB

func Close() {
  Api.Close()
  Realmd.Close()

  for _, v := range Mangosd {
    v.Character.Close()
    v.World.Close()
  }
}

func Open() error {
  var err error
  logger.Info("Initialize api database connection")
  Api, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
    config.Settings.Api.Db.Username,
    config.Settings.Api.Db.Password,
    config.Settings.Api.Db.Address,
    config.Settings.Api.Db.Port,
    config.Settings.Api.Db.Database))
  if err != nil {
    logger.Error("Cannot connect to api database")
    logger.Debug(fmt.Sprintf("%v", err))
    return err
  }

  logger.Info("Initialize realmd database connection")
  Realmd, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
    config.Settings.Realmd.Db.Username,
    config.Settings.Realmd.Db.Password,
    config.Settings.Realmd.Db.Address,
    config.Settings.Realmd.Db.Port,
    config.Settings.Realmd.Db.Database))
  if err != nil {
    logger.Error("Cannot connect to realmd database")
    logger.Debug(fmt.Sprintf("%v", err))
    Api.Close()
    return err
  }

  var db MangosdDB
  for _, v := range config.Settings.Mangosd {
    db.Id = v.Id
    logger.Info(fmt.Sprintf("Initialize mangosd[%d] character database connection", v.Id))
    db.Character, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
      v.CharacterDb.Username, v.CharacterDb.Password,
      v.CharacterDb.Address, v.CharacterDb.Port,
      v.CharacterDb.Database))
    if err != nil {
      logger.Error(fmt.Sprintf("Cannot connect to mangosd[%d] character database", v.Id))
      logger.Debug(fmt.Sprintf("%v", err))
      db.Character.Close()
      return err
    }
    logger.Info(fmt.Sprintf("Initialize mangosd[%d] world database connection", v.Id))
    db.World, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
      v.WorldDb.Username, v.WorldDb.Password,
      v.WorldDb.Address, v.WorldDb.Port,
      v.WorldDb.Database))
    if err != nil {
      logger.Error(fmt.Sprintf("Cannot connect to mangosd[%d] world database", v.Id))
      logger.Debug(fmt.Sprintf("%v", err))
      db.World.Close()
      return err
    }
    Mangosd = append(Mangosd, db)
  }

  return nil
}
