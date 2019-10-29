package character

import (
  "encoding/json"
  "net/http"
  "strconv"
  "github.com/gorilla/mux"

  e_account "github.com/CedricThomas/api-cmangos/modules/endpoints/account"
  cmangos_realm "github.com/CedricThomas/api-cmangos/cmangos/realmd/realm"
  cmangos_character "github.com/CedricThomas/api-cmangos/cmangos/mangosd/character"
  "github.com/CedricThomas/api-cmangos/modules/database"
)

func DoCharacterGuild(w http.ResponseWriter, r *http.Request) {
  id, err := e_account.AuthenticateByToken(w, r)
  if err != nil {
    return
  }

  realmlist := cmangos_realm.GetRealms()
  vars := mux.Vars(r)
  realmid, _ := strconv.Atoi(vars["realm"])
  characterid, _ := strconv.Atoi(vars["character"])

  exists := false
  for _, v := range database.Mangosd {
    if realmid == v.Id {
      exists = true
    }
  }

  if !exists {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  var guild cmangos_character.GuildMemberInfo
  for _, v := range realmlist {
    if realmid == v.Id {
      guild, err = v.CharacterInstance.GetGuild(characterid)
      if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
      }
    }
  }

  if guild.Account != id {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(guild)
  return
}
