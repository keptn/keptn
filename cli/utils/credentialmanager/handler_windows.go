package credentialmanager

func SetCreds(secret string) error {
	return setCreds(wincred.Wincred{}, secret)
}

func GetCreds() (string, error) {
	return getCreds(wincred.Wincred{})
}
