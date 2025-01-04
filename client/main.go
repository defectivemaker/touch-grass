package main

import (
	"bufio"
	"bytes"
	"client-indicum/common"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

// frame structure
// FRAMESTART | frameLength | frameType | data

// Read device UUID, public and private key from
// /etc/indicum/pub_key.pem and /etc/indicum/priv_key.pem and /etc/indicum/uuid.txt

func main() {

	// needed since it is self signed key
	config := &tls.Config{InsecureSkipVerify: true}

	address := "touchgrass.au:8888"
	deviceUUIDPath := "/etc/indicum/uuid.txt"
	deviceRSAPrivPath := "/etc/indicum/priv_key.pem"

	conn, err := tls.Dial("tcp", address, config)

	defer conn.Close()
	if err != nil {
		log.Println("Can't dial", err)
		return
	}

	deviceFileContent, err := os.ReadFile(deviceUUIDPath)
	if err != nil {
		log.Fatalf("Can't read deviceUUID file %v\n", err)
	}
	deviceUUID := string(deviceFileContent)

	deviceRSAPrivBytes, err := os.ReadFile(deviceRSAPrivPath)
	if err != nil {
		log.Fatalf("Can't read deviceRSAPriv file %v\n", err)
	}

	deviceRSAPrivBlock, _ := pem.Decode(deviceRSAPrivBytes)
	parsedPrivateKey, err := x509.ParsePKCS8PrivateKey(deviceRSAPrivBlock.Bytes)
	if err != nil {
		log.Fatalf("Can't parse from bytes to rsa.PrivKey: %v\n", err)
	}

	deviceRSAPriv := parsedPrivateKey.(*rsa.PrivateKey)

	fmt.Printf("Client Private Key Type: %T\n", deviceRSAPriv)
	fmt.Printf("Client Private Key Size: %d bits\n", deviceRSAPriv.Size()*8)
	fmt.Printf("Client Private Key Public Exponent: %d\n", deviceRSAPriv.PublicKey.E)
	fmt.Printf("Client Private Key Public Modulus: %x\n", deviceRSAPriv.PublicKey.N.Bytes())

	if len(os.Args) < 3 {
		log.Fatalf("Not enough args")
	}

	payphoneMAC := os.Args[1]
	payphoneID := os.Args[2]
	payphoneTime, _ := strconv.Atoi(os.Args[3])
	payphoneTime64 := int64(payphoneTime)

	err = sendDeviceData(conn, deviceUUID, deviceRSAPriv, payphoneMAC, payphoneID, payphoneTime64)

	if err != nil {
		log.Println("Can't send device data: ", err)
		return
	}

	err = readResponse(conn)

	if err != nil {
		log.Println("Can't read response", err)
	}

}

// function to send data about device to server
// error check by returning error to main
func sendDeviceData(conn net.Conn, deviceUUID string, deviceRSAPriv *rsa.PrivateKey, payphoneMAC, payphoneID string, payphoneTime int64) error {
	// placeholder data
	// DeviceID and PayphoneID will be 40 chars
	// will also send geoData (probably through IP geo API)
	// geoLocation := common.Coord{ Lat: -25.36364, Long: 134.21173}
	data := common.Payload{
		PayphoneMAC:  payphoneMAC,
		PayphoneID:   payphoneID,
		PayphoneTime: payphoneTime,
		Time:         time.Now().Unix(),
	}
	if len(data.PayphoneMAC) != 17 {
		return fmt.Errorf("Incorrect format for MAC\n")
	}
	if len(data.PayphoneID) != 40 {
		return fmt.Errorf("Incorrect format for PayphoneID\n")
	}
	if len(data.PayphoneID) != 40 {
		return fmt.Errorf("Incorrect format for PayphoneID\n")
	}
	hash := sha256.New()
	// add forgeResistance string to allow only trusted data (from this program)
	// to send requests to server
	forgeString := data.PayphoneID[3:len(data.PayphoneID)-3] + "_forge_resistance"
	hash.Write([]byte(forgeString))
	data.ForgeResistance = hex.EncodeToString(hash.Sum(nil))

	key1, err := hex.DecodeString(common.KeyOne)
	if err != nil {
		return fmt.Errorf("Can't decode key %s", err.Error())
	}

	dataBytes, err := json.Marshal(data)
	// encrypt data using key
	ciphertext, nonce, err := common.Encrypt(dataBytes, key1)
	if err != nil {
		return fmt.Errorf("Can't encrypt %v\n", err)
	}

	// sign data using public key
	hashedCipher := sha256.Sum256(ciphertext)

	signature, err := rsa.SignPKCS1v15(nil, deviceRSAPriv, crypto.SHA256, hashedCipher[:])

	fmt.Printf("Client Signature Length: %d bytes\n", len(signature))
	if err != nil {
		return fmt.Errorf("Failed to sign data with deviceRSA public key %v\n", err)
	}

	// sends the UUID in clear (so server knows how to decrypt)
	// sends the signature in clear (so server can verify)
	// sends the nonce in the clear (so the server can decrypt the symmetric encryption)
	var combinedData bytes.Buffer
	combinedData.Write([]byte(deviceUUID))
	combinedData.Write(signature)
	combinedData.Write(nonce)
	combinedData.Write(ciphertext)

	dataLength := combinedData.Len()

	if dataLength >= 65536 {
		return fmt.Errorf("Payload too large")
	}

	binaryDataLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(binaryDataLen, uint16(dataLength))

	var buf bytes.Buffer

	buf.Write(common.FRAMESTART[:])
	buf.WriteByte(common.FrameTypeSendDeviceData)
	buf.Write(binaryDataLen)
	buf.Write(combinedData.Bytes())
	fmt.Println("length", binaryDataLen)
	fmt.Println("deviceUUID", deviceUUID)
	fmt.Println("signature", signature)
	fmt.Println("nonce", nonce)
	fmt.Println("ciphertext", ciphertext)

	if _, err := conn.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("Failed to write payload %s\n", err)
	}

	return nil
}

func readResponse(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	response, err := io.ReadAll(reader)

	if err != nil {
		return fmt.Errorf("Can't read response: %v", err)
	}
	// Convert response to a string and print it
	fmt.Println("Client said: ", string(response))
	return nil
}
