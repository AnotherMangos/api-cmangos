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

	fmt.Printf("h1 : %s\n", h1)
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
	// basic hash creation VALID
	h1 := createBasicHash(username, password)

	// salt creation VALID
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
	// TODO remove this one
	saltInt.SetString("9920F537AC88822BAB914007DBF79E26FEEAB7F80BA98CC6694012731F226A5D", 16)
	s = fmt.Sprintf("%X", saltInt)

	// digest VALID
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

	// sha1 calculator VALID
	sha1Hasher := sha1.New()
	sha1Hasher.Write(swapEndian(saltInt.Bytes()))
	sha1Hasher.Write(swapEndian(digest))
	paper := swapEndian(sha1Hasher.Sum(nil))
	x := new(big.Int)
	x.SetBytes(paper)

	// v calculator VALID
	g := big.NewInt(7)
	N := new(big.Int)
	N.SetString("894B645E89E1535BBDAD5B8B290650530801B18EBFBF5E8FAB3C82872A3E9BB7", 16)
	verifierInt := big.NewInt(0).Exp(g, x, N)
	v = fmt.Sprintf("%X", verifierInt)

	fmt.Printf("digest1 : %s %d\n", strings.ToUpper(hex.EncodeToString(digest)), len(digest))
	fmt.Printf("s       : %s %d\n", s, len(saltInt.Bytes()))
	fmt.Printf("v       : %s %d\n", v, len(verifierInt.Bytes()))
	fmt.Printf("x       : %X %d\n", x, len(x.Bytes()))

	return
}
