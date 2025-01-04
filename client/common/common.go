package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
)

// key used to encrypt data. AES-256 GCM. Server has same key
const KeyOne = "ENTER_KEY_HERE"

type Coord struct {
	Lat  float64
	Long float64
}

// DeviceUUID can't be encrypted because we need this to find the public key mapping

// payload that is sent
type Payload struct {
	// DeviceUUID string
	PayphoneMAC  string
	PayphoneID   string
	PayphoneTime int64
	Time         int64
	// ApproxLocation Coord
	ForgeResistance string
}

type Entry struct {
	ID           int
	DeviceUUID   string
	PayphoneMAC  string
	PayphoneID   string
	PayphoneTime int64
	RecordedTime int64
	MapUUID      string
	MapLocation  string
}

type DataPoint struct {
	Point   Coord
	UUID    string
	Address string
}

type LeaderboardVal struct {
	// user
	// total unique points
	UUID        string
	TotalPoints int
	// for later:
	// total points
	// max distance between points
	//

}

// bytes to represent the start of a frame
var FRAMESTART = [2]byte{0xAA, 0x55}

// to denote different types of frames
const (
	FrameTypeSendDeviceData = 0x01
	FrameTypeGetKey         = 0x02
	FrameTypeTest           = 0x03
)

// GenerateSecureRandomString creates a cryptographically secure random string of length x.
func GenerateRandomString(x int) (string, error) {
	// Define a set of characters to use.
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, x)

	// Loop x times to generate a string of length x.
	for i := 0; i < x; i++ {
		// Generate a random index to select a character from the charset.
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err // Return the error if it occurs
		}

		// Use the random index to select a character from the charset and assign it to the result slice.
		result[i] = charset[num.Int64()]
	}

	// Convert the byte slice to a string and return it.
	return string(result), nil
}

func Encrypt(plaintext, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf(err.Error())
	}

	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, nil, fmt.Errorf(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf(err.Error())
	}
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

func Decrypt(ciphertext, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return plaintext, nil
}

func GeneratePubPrivKey() ([]byte, []byte, error) {
	// Generate a new private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Marshal the private key to PKCS#1 ASN.1 PEM.
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}
	privatePEM := pem.EncodeToMemory(&privBlock)

	// Marshal the public key to PKIX ASN.1 DER and then to PEM.
	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	pubBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	}
	publicPEM := pem.EncodeToMemory(&pubBlock)

	return publicPEM, privatePEM, nil
}
