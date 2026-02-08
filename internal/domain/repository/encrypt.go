package repository

type EncryptRepository interface {
	Encrypt(plainText string) string
	Decrypt(cipherText string) (string, error)
}
