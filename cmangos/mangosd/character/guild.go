package character

import (
  "fmt"

  "github.com/CedricThomas/api-cmangos/modules/logger"
)

type GuildInfo struct {
  Guildid int `json:"guildid"`
  Name string `json:"name"`
  Leaderguid int `json:"leaderguid"`
  EmblemStyle int64 `json:"EmblemStyle"`
  EmblemColor int64 `json:"EmblemColor"`
  BorderStyle int64 `json:"BorderStyle"`
  BorderColor int64 `json:"BorderColor"`
  BackgroundColor int64 `json:"BackgroundColor"`
  Info string `json:"info"`
  Motd string `json:"motd"`
  Createdate int64 `json:"createdate"`
}

type GuildRankInfo struct {
  Guildid int `json:"guildid"`
  Rid int `json:"rid"`
  Rname string `json:"rname"`
  Rights int `json:"rights"`
}

type GuildMemberInfo struct {
  Guild GuildInfo `json:"guild"`
  Guid int `json:"guid"`
  Rank GuildRankInfo `json:"rank"`
  Pnote string `json:"pnote"`
  Offnote string `json:"offnote"`
  Account int `json:"account"`
}

func (c *CharacterInstanceInfo) GetGuild(id int) (GuildMemberInfo, error) {
  var gmi GuildMemberInfo
  stmt, err := c.Db.Prepare(
    `SELECT
       g.guildid, g.name, g.leaderguid, g.EmblemStyle, g.EmblemColor,
       g.BorderStyle, g.BorderColor, g.BackgroundColor, g.info, g.motd, g.createdate,
       gm.guid, gr.guildid, gr.rid, gr.rname, gr.rights, gm.pnote, gm.offnote, c.account
     FROM guild_member AS gm
     INNER JOIN guild AS g ON g.guildid = gm.guildid
     INNER JOIN guild_rank AS gr ON gr.rid = gm.rank
     INNER JOIN characters AS c ON gm.guid = c.guid
     WHERE gm.guid = ?;`)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot prepare query to fetch guild info %d", id))
    logger.Debug(fmt.Sprintf("%v", err))
    return gmi, err
  }
  defer stmt.Close()

  err = stmt.QueryRow(id).Scan(
    &gmi.Guild.Guildid, &gmi.Guild.Name, &gmi.Guild.Leaderguid,
    &gmi.Guild.EmblemStyle, &gmi.Guild.EmblemColor, &gmi.Guild.BorderStyle,
    &gmi.Guild.BorderColor, &gmi.Guild.BackgroundColor, &gmi.Guild.Info,
    &gmi.Guild.Motd, &gmi.Guild.Createdate, &gmi.Guid, &gmi.Rank.Guildid,
    &gmi.Rank.Rid, &gmi.Rank.Rname, &gmi.Rank.Rights, &gmi.Pnote, &gmi.Offnote,
    &gmi.Account)
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot query guild info %d", id))
    logger.Debug(fmt.Sprintf("%v", err))
    return gmi, err
  }

  return gmi, nil
}
