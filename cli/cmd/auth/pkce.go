package auth

import (
	"crypto/rand"
	"math/big"
)

func GenerateCodeVerifier() (codeVerifier []byte, err error) {
	max := big.NewInt(int64(len(CodeVerifierChars)))
	var n *big.Int
	for i := 0; i != 128; i++ {
		n, err = rand.Int(rand.Reader, max)
		if err != nil {
			return
		}
		codeVerifier = append(codeVerifier, CodeVerifierChars[n.Int64()])
	}
	return
}

var (
	CodeVerifierChars []byte
)

func init() {
	for b := byte('a'); b <= byte('z'); b++ {
		CodeVerifierChars = append(CodeVerifierChars, b)
	}
	for b := byte('A'); b <= byte('Z'); b++ {
		CodeVerifierChars = append(CodeVerifierChars, b)
	}
	for b := byte('0'); b <= byte('9'); b++ {
		CodeVerifierChars = append(CodeVerifierChars, b)
	}
	CodeVerifierChars = append(CodeVerifierChars, byte('-'), byte('.'), byte('_'), byte('~'))
}
