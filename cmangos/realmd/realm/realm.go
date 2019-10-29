package realm

import (
  "github.com/CedricThomas/api-cmangos/cmangos"
  "github.com/CedricThomas/api-cmangos/cmangos/mangosd/character"
)

type Realm struct {
  Id int `json:"id,omitempty"`
  Name string `json:"name,omitempty"`
  Icon int `json:"icon,omitempty"`
  Population float64 `json:"population,omitempty"`
  Host cmangos.DaemonInfo `json:"host,omitempty"`
  CharacterInstance character.CharacterInstanceInfo `json:"-"`
}
