package character

import (
  "fmt"
  "database/sql"
  "time"

  "metagit.org/blizzlike/cmangos-api/modules/logger"
)

type CharacterInfo struct {
  Guid int64 `json:"guid"`
  Account int64 `json:"account"`
  Name string `json:"name"`
  Race int `json:"race"`
  Class int `json:"class"`
  Gender int `json:"gender"`
  Level int `json:"level"`
  Xp int64 `json:"xp"`
  Money int64 `json:"money"`
  PlayerBytes int64 `json:"playerBytes"`
  PlayerBytes2 int64 `json:"playerBytes2"`
  PlayerFlags int64 `json:"playerFlags"`
  Position_x float64 `json:"position_x"`
  Position_y float64 `json:"position_y"`
  Position_z float64 `json:"position_z"`
  Map int64 `json:"map"`
  Orientation float64 `json:"orientation"`
  Taximask string `json:"taximask"`
  Online int `json:"online"`
  Cinematic int `json:"cinematic"`
  Totaltime int64 `json:"totaltime"`
  Leveltime int64 `json:"leveltime"`
  Logout_time int64 `json:"logout_time"`
  Is_logout_resting int `json:"is_logout_resting"`
  Rest_bonus float64 `json:"rest_bonus"`
  Resettalents_cost int64 `json:"resettalents_cost"`
  Resettalents_time int64 `json:"resettalents_time"`
  Trans_x float64 `json:"trans_x"`
  Trans_y float64 `json:"trans_y"`
  Trans_z float64 `json:"trans_z"`
  Trans_o float64 `json:"trans_o"`
  Transguid int64 `json:"transguid"`
  Extra_flags int64 `json:"extra_flags"`
  Stable_slots int `json:"stable_slots"`
  At_login int64 `json:"at_login"`
  Zone int64 `json:"zone"`
  Death_expire_time int64 `json:"death_expire_time"`
  Taxi_path string `json:"taxi_path"`
  Honor_highest_rank int64 `json:"honor_highest_rank"`
  Honor_standing int64 `json:"honor_standing"`
  Stored_honor_rating float64 `json:"stored_honor_rating"`
  Stored_dishonorable_kills int64 `json:"stored_dishonorable_kills"`
  Stored_honorable_kills int64 `json:"stored_honorable_kills"`
  WatchedFaction int64 `json:"watchedFaction"`
  Drunk int `json:"drunk"`
  Health int64 `json:"health"`
  Power1 int64 `json:"power1"`
  Power2 int64 `json:"power2"`
  Power3 int64 `json:"power3"`
  Power4 int64 `json:"power4"`
  Power5 int64 `json:"power5"`
  ExploredZones string `json:"exploredZones"`
  EquipmentCache string `json:"equipmentCache"`
  AmmoId int64 `json:"ammoId"`
  ActionBars int `json:"actionBars"`
  DeleteInfos_Account sql.NullInt64 `json:"deleteInfos_Account"`
  DeleteInfos_Name sql.NullString `json:"deleteInfos_Name"`
  DeleteDate sql.NullInt64 `json:"deleteDate"`
  GuildMember GuildMemberInfo `json:"guildmember"`
}

func (c *CharacterInfo) GetFaction() string {
  if c.Race == 1 || c.Race == 3 || c.Race == 4 || c.Race == 7 {
    return "Alliance"
  }
  if c.Race == 2 || c.Race == 5 || c.Race == 6 || c.Race == 8 {
    return "Horde"
  }
  return ""
}

func (c *CharacterInfo) GetRace() string {
  if c.Race == 1 { return "Human" }
  if c.Race == 2 { return "Orc" }
  if c.Race == 3 { return "Dwarf" }
  if c.Race == 4 { return "Nightelf" }
  if c.Race == 5 { return "Undead" }
  if c.Race == 6 { return "Tauren" }
  if c.Race == 7 { return "Gnome" }
  if c.Race == 8 { return "Troll" }
  return ""
}

func (c *CharacterInfo) GetGender() string {
  if c.Gender == 0 { return "male"}
  if c.Gender == 1 { return "female"}
  return ""
}

func (c *CharacterInfo) GetClass() string {
  if c.Class == 1 { return "Warrior" }
  if c.Class == 2 { return "Paladin" }
  if c.Class == 3 { return "Hunter" }
  if c.Class == 4 { return "Rogue" }
  if c.Class == 5 { return "Priest" }
  if c.Class == 7 { return "Shaman" }
  if c.Class == 8 { return "Mage" }
  if c.Class == 9 { return "Warlock" }
  if c.Class == 11 { return "Druid" }
  return ""
}

func (c *CharacterInfo) GetGold() int64 {
  return c.Money / 10000
}

func (c *CharacterInfo) GetSilver() int64 {
  return (c.Money % 10000) / 100
}

func (c *CharacterInfo) GetCopper() int64 {
  return c.Money % 10000 % 100
}

func (c *CharacterInfo) LoggedOutSince() string {
  var w, d, h, m, s int
  msg := "-"

  t := time.Now().Unix() - c.Logout_time
  if t > 0 {
    s = int(t % 60)
    m = int(t / 60 % 60)
    h = int(t / 3600 % 24)
    d = int(t / 3600 / 24)
    w = int(t / 86400 / 7)
  } else {
    return msg
  }

  if w > 4 {
    msg = time.Unix(c.Logout_time, 0).Format("2006-01-02")
  } else if d > 0 {
    msg = fmt.Sprintf("%dd %dh", d, h)
  } else if h > 0 {
    msg = fmt.Sprintf("%dh %dm", h, m)
  } else if m > 0 {
    msg = fmt.Sprintf("%dm %ds", m, s)
  } else {
    msg = fmt.Sprintf("%ds", s)
  }

  return msg
}

func (c *CharacterInfo) PlayedTime() string {
  var d, h, m, s int
  msg := "-"

  t := c.Totaltime
  if t > 0 {
    s = int(t % 60)
    m = int(t / 60 % 60)
    h = int(t / 3600 % 24)
    d = int(t / 3600 / 24)
  } else {
    return msg
  }

  if d > 0 {
    msg = fmt.Sprintf("%dd %dh %dm", d, h, m)
  } else if h > 0 {
    msg = fmt.Sprintf("%dh %dm %ds", h, m, s)
  } else if m > 0 {
    msg = fmt.Sprintf("%dm %ds", m, s)
  } else {
    msg = fmt.Sprintf("%ds", s)
  }

  return msg
}

func (c *CharacterInstanceInfo) GetCharacterByAccountId(id int) ([]CharacterInfo, error) {
  var ci []CharacterInfo
  stmt, err := c.Db.Prepare(
    `SELECT
       guid, account, name, race, class, gender,
       level, xp, money, online, totaltime, leveltime,
       logout_time, is_logout_resting, rest_bonus, drunk, health,
       deleteInfos_Account, deleteInfos_Name, deleteDate
     FROM characters
     WHERE account = ?;`)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot prepare query to fetch all characters of account %d", id))
    logger.Debug(fmt.Sprintf("%v", err))
    return ci, err
  }
  defer stmt.Close()

  var rows *sql.Rows
  rows, err = stmt.Query(id)
  for rows.Next() {
    var c CharacterInfo
    err = rows.Scan(
      &c.Guid, &c.Account, &c.Name, &c.Race, &c.Class,
      &c.Gender, &c.Level, &c.Xp, &c.Money, &c.Online, &c.Totaltime,
      &c.Leveltime, &c.Logout_time, &c.Is_logout_resting, &c.Rest_bonus,
      &c.Drunk, &c.Health, &c.DeleteInfos_Account, &c.DeleteInfos_Name, &c.DeleteDate)
    if err != nil {
      logger.Error(fmt.Sprintf("Cannot query for character of account %d", id))
      logger.Debug(fmt.Sprintf("%v", err))
      return ci, err
    }

    ci = append(ci, c)
  }

  return ci, nil
}
