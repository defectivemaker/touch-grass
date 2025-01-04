// This file will have all the functions that handle the TCP requests/responses that come from the RPi itself (using my own protocol on top)
package device

import (
    "os"
    "fmt"
    "log"
    "crypto/tls"
    "crypto"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "encoding/hex"
    "encoding/binary"
    "encoding/json"
    "crypto/sha256"
    "sync/atomic"
    "net"
    "io"

    "server-indicum/internal/common"
    "server-indicum/internal/server/db"

)



func InitDeviceServer() {

    tlsCert := os.Getenv("TLS_CERT_FILE")
    tlsPrivkey := os.Getenv("TLS_PRIV_KEY")
    cert, err := tls.LoadX509KeyPair( tlsCert, tlsPrivkey )
    if err != nil { log.Fatalf("Failed to load key pair: %v", err) }

    config := &tls.Config{Certificates: []tls.Certificate{cert},}
    tcpListen := os.Getenv("TCPLISTENADDRESS")
    ln, err := tls.Listen("tcp", tcpListen, config) 

    if err != nil { log.Fatalf("Failed to listen on %s: %v", tcpListen, err)}

    fmt.Println("TCP Server listening on address", tcpListen)

    var connectionCount uint64 = 0
    for {
        conn, err := ln.Accept()
        atomic.AddUint64(&connectionCount, 1)
        fmt.Println("Connection count:", connectionCount)
        if err != nil { 
            log.Println("Error accepting connection:", err)
            continue
        }
        go handleDeviceConnection(conn)
    }
}

// If there is an error, log.Printf() the error and then early return
// handleDeviceConnection will then just close the connection and move on
func handleDeviceConnection(conn net.Conn) {
    defer conn.Close()

    // First read the first N bytes to ensure that it is coming from one of my devices
    var frameCheck [2]byte
    binary.Read(conn, binary.LittleEndian, &frameCheck)

    if frameCheck != common.FRAMESTART {
        log.Println("FRAMESTART doesn't match"); return 
    }
    var frameType byte
    binary.Read(conn, binary.LittleEndian, &frameType)

    var err error
    switch frameType {
        case common.FrameTypeSendDeviceData:
            err = handleDeviceData(conn)
        case common.FrameTypeGetKey:
        case common.FrameTypeTest:
        default:
            log.Println("Frame type invalid")
            return
    }
    if err != nil { log.Printf("using frametype %x led to :%v\n", frameType, err); return }

    err = sendResponse(conn)

    if err != nil { log.Println("can't send response", err); return }
}

// function that handles when the device sends data about itself to server
// will include PayphoneID, payphoneID, geodata etc
func handleDeviceData(conn net.Conn) error {

    key1, err := hex.DecodeString(common.KeyOne)
    if err != nil { return fmt.Errorf("Can't decode key %v\n", err)}

    // reads 2 bytes (for frame length) and then reads length bytes
    var length uint16
    binary.Read(conn, binary.LittleEndian, &length)
    if length < 310 { return fmt.Errorf("Length is only %d. Make sure you are sending the correct data\n")}
    data := make([]byte, length)
    _, err = io.ReadFull(conn, data)
    if err != nil { return fmt.Errorf("Can't read data %v\n", err)}

    // first 36 bytes is deviceUUID
    deviceUUID := data[:36]
    // next 256 bytes is signature
    signature := data[36:292]
    // next 12 bytes is nonce
    nonce := data[292:304]
    // rest is ciphertext
    ciphertext := data[304:]

    deviceRSAPubBytes, err := db.DBFindDeviceRSAPub(string(deviceUUID))
    if err != nil { return fmt.Errorf("Failed to get deviceRSAPub %v\n", err)}
    fmt.Println("publickeybytes", deviceRSAPubBytes)

    deviceRSAPubBlock, _ := pem.Decode(deviceRSAPubBytes)
    parsedPublicKey, err := x509.ParsePKIXPublicKey(deviceRSAPubBlock.Bytes)
    if err != nil { log.Fatalf("Can't parse from bytes to rsa.PubKey: %v\n", err)}
    
    deviceRSAPub := parsedPublicKey.(*rsa.PublicKey)
    fmt.Println("length", length)
    fmt.Println("deviceUUID", deviceUUID)
    fmt.Println("signature", signature)
    fmt.Println("nonce", nonce)
    fmt.Println("ciphertext", ciphertext)

    hashedCipher := sha256.Sum256(ciphertext)

    fmt.Println("publickey", deviceRSAPub)
    fmt.Printf("Server Public Key Type: %T\n", deviceRSAPub)
    fmt.Printf("Server Public Key Size: %d bits\n", deviceRSAPub.Size()*8)
    fmt.Printf("Server Public Key Exponent: %d\n", deviceRSAPub.E)
    fmt.Printf("Server Public Key Modulus: %x\n", deviceRSAPub.N.Bytes())

    fmt.Printf("Server Received Signature Length: %d bytes\n", len(signature))
    fmt.Printf("Server Expected Signature Length: %d bytes\n", deviceRSAPub.Size())


    err = rsa.VerifyPKCS1v15(deviceRSAPub, crypto.SHA256, hashedCipher[:], signature)
    if err != nil { return fmt.Errorf("Failed to sign ciphertext, integrity compromised %v\n", err)}

    plaintext, err := common.Decrypt(ciphertext, key1, nonce)
    if err != nil { return fmt.Errorf("Can't decrypt %v\n", err)}

    // unmarshals data
    var dataPayload common.Payload 
    err = json.Unmarshal(plaintext, &dataPayload)
    if err != nil { return fmt.Errorf("Can't unmarshal %v\n", err)}

    // fmt.Println("data", dataPayload)

    if len(dataPayload.PayphoneID) <= 6 { return fmt.Errorf("PayphoneID too short\n") }
    forgeString := dataPayload.PayphoneID[3:len(dataPayload.PayphoneID)-3] + "_forge_resistance"
    forgeResistanceHash := sha256.New()
    forgeResistanceHash.Write([]byte(forgeString))
    forgeHashString := hex.EncodeToString(forgeResistanceHash.Sum(nil))
    
    // verifies that there was no tampering or forging
    if forgeHashString != dataPayload.ForgeResistance{
        return fmt.Errorf("Tampering/Forgery detected\n")
    }

    id, err := db.AddEntryToDB(dataPayload, string(deviceUUID))

    if err != nil { return fmt.Errorf("Failed to add to DB: %v\n", err)}
    
    fmt.Println("Added data to DB with id", id)
    

    return nil
}

func sendResponse(conn net.Conn) error {
    response := []byte("ty\n")
    _, err := conn.Write(response)
    if err != nil { return fmt.Errorf("Can't write response %s", err.Error()) }
    fmt.Println("message sent:", string(response))
    return nil
}
