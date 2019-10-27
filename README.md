# cmangos-api

a rest api for the cmangos core database.

## api

### examples

#### create account

Based on the config the `Authorization` header has to be set to a `invite` token.
Set `needInvite` in `contrib/config.ini.dist`.

    curl \
      -X POST \
      --header "Authorization: Token <invite-token>" \
      --header "Content-Type: application/json" \
      --data '{"username":"test","password":"test","repeat":"test","email":"test@example.org"}' \
      http://127.0.0.1:5556/account

To create an invite token you have to authenticate yourself first.

    curl \
      -X POST \
      --user "<username>:<password>" \
      http://127.0.0.1:5556/account/auth

You will receive a `X-Auth-Token` header containing a token for further requests.

    curl \
      -X POST \
      --header "Authorization: Token <token>" \
      http://127.0.0.1:5556/account/invite
# api-cmangos
