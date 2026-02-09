package ports

type EncryptPort interface {
	Encrypt(plainText string) string
	Decrypt(cipherText string) (string, error)
}
