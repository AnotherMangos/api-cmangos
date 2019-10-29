package config

import (
  ini "gopkg.in/ini.v1"
)

type ApiConfig struct {
  Listen string
  Port int
  Db DBConfig
  RequireInvite bool
  AuthTokenExpiry int
  CheckTimeout int
  CheckInterval int
  Verbosity int
}

type DBConfig struct {
  Username string
  Password string
  Address string
  Port int
  Database string
}

type RealmdConfig struct {
  Address string
  Port int
  Db DBConfig
}

type MangosdConfig struct{
  Id int
  CharacterDb DBConfig
  WorldDb DBConfig
}

type Config struct {
  Api ApiConfig
  Realmd RealmdConfig
  Mangosd []MangosdConfig
}

var Settings Config

func Read(file string) (Config, error) {
  c, err := ini.Load(file)
  if err != nil {
    return Settings, err
  }

  // [api]
  Settings.Api.Listen = c.Section("api").Key("listen").MustString("127.0.0.1")
  Settings.Api.Port = c.Section("api").Key("port").MustInt(5556)
  Settings.Api.RequireInvite = c.Section("api").Key("requireInvite").MustBool(false)
  Settings.Api.AuthTokenExpiry = c.Section("api").Key("expiry").MustInt(3600)
  Settings.Api.CheckTimeout = c.Section("api").Key("timeout").MustInt(10)
  Settings.Api.CheckInterval = c.Section("api").Key("interval").MustInt(300)
  Settings.Api.Verbosity = c.Section("api").Key("loglevel").MustInt(1)

  // [api.mysql]
  Settings.Api.Db.Address = c.Section("api.mysql").Key("hostname").MustString("127.0.0.1")
  Settings.Api.Db.Port = c.Section("api.mysql").Key("port").MustInt(3306)
  Settings.Api.Db.Username = c.Section("api.mysql").Key("username").MustString("api-cmangos")
  Settings.Api.Db.Password = c.Section("api.mysql").Key("password").MustString("api-cmangos")
  Settings.Api.Db.Database = c.Section("api.mysql").Key("database").MustString("api-cmangos")

  // [realmd]
  Settings.Realmd.Address = c.Section("realmd").Key("hostname").MustString("logon.example.org")
  Settings.Realmd.Port = c.Section("realmd").Key("port").MustInt(3724)

  // [realmd.mysql]
  Settings.Realmd.Db.Address = c.Section("realmd.mysql").Key("hostname").MustString("127.0.0.1")
  Settings.Realmd.Db.Port = c.Section("realmd.mysql").Key("port").MustInt(3306)
  Settings.Realmd.Db.Username = c.Section("realmd.mysql").Key("username").MustString("mangos")
  Settings.Realmd.Db.Password = c.Section("realmd.mysql").Key("password").MustString("mangos")
  Settings.Realmd.Db.Database = c.Section("realmd.mysql").Key("database").MustString("realmd")

  // [mangosd]
  realms := c.Section("mangosd").Key("realms").Strings(",")

  // [mangosd.<realm>.character.mysql]
  // [mangosd.<realm>.world.mysql]
  if len(realms) != 0 {
    for _, v := range realms {
      realm := MangosdConfig{
        Id: c.Section("mangosd." + v).Key("id").MustInt(1),
	CharacterDb: DBConfig{
          Address: c.Section("mangosd." + v + ".character.mysql").Key("hostname").MustString("127.0.0.1"),
          Port: c.Section("mangosd." + v + ".character.mysql").Key("port").MustInt(3306),
          Username: c.Section("mangosd." + v + ".character.mysql").Key("username").MustString("mangos"),
          Password: c.Section("mangosd." + v + ".character.mysql").Key("password").MustString("mangos"),
          Database: c.Section("mangosd." + v + ".character.mysql").Key("database").MustString("character"),
	},
	WorldDb: DBConfig{
          Address: c.Section("mangosd." + v + ".world.mysql").Key("hostname").MustString("127.0.0.1"),
          Port: c.Section("mangosd." + v + ".world.mysql").Key("port").MustInt(3306),
          Username: c.Section("mangosd." + v + ".world.mysql").Key("username").MustString("mangos"),
          Password: c.Section("mangosd." + v + ".world.mysql").Key("password").MustString("mangos"),
          Database: c.Section("mangosd." + v + ".world.mysql").Key("database").MustString("world"),
	},
      }
      Settings.Mangosd = append(Settings.Mangosd, realm)
    }
  }

  return Settings, nil
}
