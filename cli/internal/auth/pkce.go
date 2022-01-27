package auth

import (
	"crypto/rand"
	"math/big"
)

var (
	CodeVerifierCharset []byte
)

// GenerateCodeVerifier creates a new CodeVerifier usable to prevent man-in-the-middle attacks
func GenerateCodeVerifier() (codeVerifier []byte, err error) {
	max := big.NewInt(int64(len(CodeVerifierCharset)))
	var n *big.Int
	for i := 0; i != 128; i++ {
		n, err = rand.Int(rand.Reader, max)
		if err != nil {
			return
		}
		codeVerifier = append(codeVerifier, CodeVerifierCharset[n.Int64()])
	}
	return
}

// init initializes the character set used for the verifier ( [A-Z] [a-z] [0-9] -  .  ~ )
func init() {
	for b := byte('a'); b <= byte('z'); b++ {
		CodeVerifierCharset = append(CodeVerifierCharset, b)
	}
	for b := byte('A'); b <= byte('Z'); b++ {
		CodeVerifierCharset = append(CodeVerifierCharset, b)
	}
	for b := byte('0'); b <= byte('9'); b++ {
		CodeVerifierCharset = append(CodeVerifierCharset, b)
	}
	CodeVerifierCharset = append(CodeVerifierCharset, byte('-'), byte('.'), byte('_'), byte('~'))
}
