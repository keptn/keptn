package cmd

import (
	"testing"
)

func TestMigrateAlphaRequest(t *testing.T) {
	_, _ = migrateAlphaRequest("curl --data '{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}' -H \"Accept-Charset: utf-8\" -H 'Content-Type: application/json' https://httpbin.org/post --some-random-options -YYY -X POST")
}
