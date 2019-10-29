package account

import (
  "errors"
  "fmt"

	"github.com/CedricThomas/api-cmangos/modules/database"
	"github.com/CedricThomas/api-cmangos/modules/logger"
)

func Authenticate(username, password string) (AccountInfo, error) {
	var a AccountInfo
	stmt, err := database.Realmd.Prepare(
		`SELECT id, username, v, s FROM account
                WHERE UPPER(username) = UPPER(?);`)
	if err != nil {
		logger.Error(fmt.Sprintf("Cannot prepare query to authenticate %s", username))
		logger.Debug(fmt.Sprintf("%v", err))
		return a, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(
		&a.Id, &a.Username, &a.V, &a.S)
	if err != nil {
		logger.Error(fmt.Sprintf("Cannot authenticate account %s", username))
		logger.Debug(fmt.Sprintf("%v", err))
		return a, err
	}

	v := CreateVerifier(username, password, a.S)
	if v != a.V {
      logger.Error(fmt.Sprintf("Cannot authenticate account %s", username))
	  return a, errors.New(fmt.Sprintf("Cannot authenticate account %s", username))
    }

	return a, nil
}
