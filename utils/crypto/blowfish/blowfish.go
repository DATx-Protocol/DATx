package blowfish

import (
	"crypto/cipher"

	"golang.org/x/crypto/blowfish"
)

func blowfishChecksizeAndPad(pt []byte) []byte {
	modulus := len(pt) % blowfish.BlockSize
	if modulus != 0 {
		padlen := blowfish.BlockSize - modulus
		for i := 0; i < padlen; i++ {
			pt = append(pt, 0)
		}
	}
	return pt
}

func blowfishDecrypt(et, key []byte) []byte {
	dcipher, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err)
	}
	div := et[:blowfish.BlockSize]
	decrypted := et[blowfish.BlockSize:]
	if len(decrypted)%blowfish.BlockSize != 0 {
		panic("decrypted is not a multiple of blowfish.BlockSize")
	}
	dcbc := cipher.NewCBCDecrypter(dcipher, div)
	dcbc.CryptBlocks(decrypted, decrypted)
	return decrypted
}

func blowfishEncrypt(ppt, key []byte) []byte {
	ecipher, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, blowfish.BlockSize+len(ppt))
	eiv := ciphertext[:blowfish.BlockSize]
	ecbc := cipher.NewCBCEncrypter(ecipher, eiv)
	ecbc.CryptBlocks(ciphertext[blowfish.BlockSize:], ppt)
	return ciphertext
}

/*
func main() {
	var decryptedtext, encryptedtext, plaintext, paddedplaintext, secretkey []byte
	plaintext = []byte("this is the plaintext string")
	secretkey = []byte("1234567890abcdefghijklmnopqrstuvwxyz")
	paddedplaintext = blowfishChecksizeAndPad(plaintext)
	encryptedtext = blowfishEncrypt(paddedplaintext, secretkey)
	decryptedtext = blowfishDecrypt(encryptedtext, secretkey)
	fmt.Printf("      plaintext=%s\n", plaintext)
	fmt.Printf("  encryptedtext=%x\n", encryptedtext)
	fmt.Printf("  decryptedtext=%s\n", decryptedtext)
}
*/
