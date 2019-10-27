package account

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"strings"
)

const (
	SALT_SIZE         = 32
	SHA_DIGEST_LENGTH = 20
)

func createBasicHash(username, password string) string {
	username = strings.ToUpper(username)
	password = strings.ToUpper(password)
	sha1Hasher := sha1.New()
	_, _ = io.WriteString(sha1Hasher, username)
	_, _ = io.WriteString(sha1Hasher, ":")
	_, _ = io.WriteString(sha1Hasher, password)
	sha1Result := sha1Hasher.Sum(nil)
	h1 := strings.ToUpper(hex.EncodeToString(sha1Result))

	return h1
}

func swapEndian(a []byte) []byte {
	size := len(a)
	b := make([]byte, size)
	for index := range a {
		b[index] = a[size - 1 - index]
	}
	return b
}

func CreateSaltAndVerifier(username, password string) (s, v string) {

	maxSalt := big.NewInt(SALT_SIZE * 8)
	minSalt := big.NewInt(SALT_SIZE * 8)
	minSalt.Exp(big.NewInt(2), minSalt.Sub(minSalt, big.NewInt(1)), nil)
	maxSalt.Exp(big.NewInt(2), maxSalt, nil).Sub(maxSalt, big.NewInt(1))
	plage := big.NewInt(0).Sub(maxSalt, minSalt)
	saltInt, _ := rand.Int(rand.Reader, plage)
	saltInt.Add(saltInt, minSalt)
	if big.NewInt(0).Mod(saltInt, big.NewInt(2)).Int64() == 0 {
		saltInt.Add(saltInt, big.NewInt(1))
	}
	s = fmt.Sprintf("%X", saltInt)
	v = CreateVerifier(username, password, s)
	return
}

func CreateVerifier(username, password, salt string) (v string) {
	h1 := createBasicHash(username, password)

	h1Int := new(big.Int)
	h1Int.SetString(h1, 16)
	h1IntRaw := h1Int.Bytes()
	digest := make([]byte, SHA_DIGEST_LENGTH)
	if len(h1IntRaw) <= SHA_DIGEST_LENGTH {
		for i, byteValue := range h1IntRaw {
			digest[i] = byteValue
		}
	}

	for i := len(digest)/2 - 1; i >= 0; i-- {
		opp := len(digest) - 1 - i
		digest[i], digest[opp] = digest[opp], digest[i]
	}

	saltInt := new(big.Int)
	saltInt.SetString(salt, 16)
	sha1Hasher := sha1.New()
	sha1Hasher.Write(swapEndian(saltInt.Bytes()))
	sha1Hasher.Write(swapEndian(digest))
	paper := swapEndian(sha1Hasher.Sum(nil))
	x := new(big.Int)
	x.SetBytes(paper)

	g := big.NewInt(7)
	N := new(big.Int)
	N.SetString("894B645E89E1535BBDAD5B8B290650530801B18EBFBF5E8FAB3C82872A3E9BB7", 16)
	verifierInt := big.NewInt(0).Exp(g, x, N)
	v = fmt.Sprintf("%X", verifierInt)

	return
}
