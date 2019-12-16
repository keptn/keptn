package utils

import (
	"math/rand"
	"testing"
	"time"
)

// generateStringWithSpecialChars generates a string of the given length
// and containing at least one special character and digit.
func generateStringWithSpecialChars(length int) string {
	rand.Seed(time.Now().UnixNano())

	digits := "0123456789"
	specials := "-~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials

	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]

	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}

	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})

	str := string(buf)

	return str
}

// TestValidateK8sName tests whether a random string containing a special character or digit
// does not pass the name validation.
func TestValidateK8sName(t *testing.T) {
	invalidName := generateStringWithSpecialChars(8)
	if ValidateK8sName(invalidName) {
		t.Fatalf("%s starts with upper case letter(s) or contains special character(s), but passed the name validation", invalidName)
	}
}

var keptnVersions = []struct {
	in  string
	res bool
}{
	{"20191212.1033-latest", false},
	{"0.6.0.beta2", true},
	{"feature-443-20191213.1105", false},
	{"0.6.0.beta2-20191204.1329", false},
	{"0.6.0.beta2-201912044.1329", true},
}

func TestIsOfficialKeptnVersion(t *testing.T) {
	for _, tt := range keptnVersions {
		t.Run(tt.in, func(t *testing.T) {
			res := IsOfficialKeptnVersion(tt.in)
			if res != tt.res {
				t.Errorf("got %t, want %t", res, tt.res)
			}
		})
	}
}
