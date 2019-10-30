package main

import (
  "fmt"
  "log"
  "os"
  "time"
  "net/http"
  "github.com/gorilla/handlers"
  "github.com/gorilla/mux"

  "github.com/CedricThomas/api-cmangos/modules/logger"

  cmangos_realm "github.com/CedricThomas/api-cmangos/cmangos/realmd/realm"

  "github.com/CedricThomas/api-cmangos/modules/config"
  "github.com/CedricThomas/api-cmangos/modules/database"
  e_account "github.com/CedricThomas/api-cmangos/modules/endpoints/account"
  e_config "github.com/CedricThomas/api-cmangos/modules/endpoints/config"
  e_realm "github.com/CedricThomas/api-cmangos/modules/endpoints/realm"
  e_character "github.com/CedricThomas/api-cmangos/modules/endpoints/realm/character"
)

func main() {
  if len(os.Args) != 2 {
    logger.Error(fmt.Sprintf("USAGE: %s <config>", os.Args[0]))
    os.Exit(1)
  }
  _, err := config.Read(os.Args[1])
  if err != nil {
    logger.Error(fmt.Sprintf("Failed to read file (%v)", err))
    os.Exit(2)
  }

  logger.Verbosity = config.Settings.Api.Verbosity

  err = database.Open()
  if err != nil {
    logger.Error(
      fmt.Sprintf("Cannot initialize database connections (%v)", err))
    os.Exit(3)
  }
  defer database.Close()

  logger.Info("Initialize RealmStates poller")
  go cmangos_realm.PollRealmStates(
    time.Duration(config.Settings.Api.CheckInterval))

  logger.Info("Initialize url multiplexer")
  router := mux.NewRouter()
  router.HandleFunc("/account", e_account.DoGetAccount).Methods("GET")
  router.HandleFunc("/account", e_account.DoCreateAccount).Methods("POST")
  router.HandleFunc("/account/auth", e_account.DoAuthVerify).Methods("GET")
  router.HandleFunc("/account/auth", e_account.DoAuth).Methods("POST")
  router.HandleFunc("/account/invite", e_account.DoGetInvites).Methods("GET") // not working
  router.HandleFunc("/account/invite", e_account.DoInvite).Methods("POST")

  router.HandleFunc("/config", e_config.DoConfig).Methods("GET")

  router.HandleFunc("/realm", e_realm.DoRealmlist).Methods("GET")
  router.HandleFunc("/realm/{realm}/characters/{account}",
    e_character.DoCharacterlistByAccount).Methods("GET")

  router.HandleFunc("/realm/{realm}/character/{character}/guild",
    e_character.DoCharacterGuild).Methods("GET")

  originsOk := handlers.AllowedOrigins([]string{"*"})

  logger.Info("Start serving http requests")
  log.Fatal(http.ListenAndServe(
    fmt.Sprintf("%s:%d", config.Settings.Api.Listen, config.Settings.Api.Port), handlers.CORS(originsOk)(router)))
}
