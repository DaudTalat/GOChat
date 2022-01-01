package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var encryption_key string

func main() {
	// connect to server
	cn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		log.Fatalf("Error: Unable to open TCP Connection: %s", err)
	}

	defer cn.Close()

	time.Sleep(2)
	print("Listener Started. Enter Commands:" + "\n")

	for {
		go listen(cn)
		send_cmd(cn)
	}
}

func listen(cn net.Conn) {

	message, err := bufio.NewReader(cn).ReadString('\n')

	if err != nil {
		log.Fatalf("Error from TCP Session: %s", err)
	}

	if len(encryption_key) != 0 && strings.Contains(message, "encrypted:") {
		fmt.Print("->: " + decrypt((strings.Split(message, " ")[2])[10:], encryption_key))
	} else {
		fmt.Print("->: " + message)
	}
}

func send_cmd(cn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	args := strings.Split(text, " ")
	cmd := strings.TrimSpace(args[0])

	// encryption
	if cmd == "/encrypt" && len(args) == 1 {
		encryption_key = generateKey()
		print("ENCRYPTION KEY: " + encryption_key + "\n")
		return
	} else if cmd == "/encrypt" {
		encryption_key = args[1]
		return
	}

	// plaintext
	if len(encryption_key) == 0 || cmd != "/msg" {
		fmt.Fprintf(cn, text+"\n")
	} else {
		encrypted_text := args[0] + " encrypted:" + encrypt(text[4:], encryption_key) + "\n"
		fmt.Fprintf(cn, encrypted_text)
	}
}

//------Hashing Functions-----//
func generateKey() string {
	bytes := make([]byte, 32) // random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes) //encode key in bytes to string for saving
	return key
}

func encrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
