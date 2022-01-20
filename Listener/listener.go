package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var encryption_key string
var backlog []string

func main() {
	// connect to server
	cn, err := net.Dial("tcp", "localhost:8080")

	//setting http client to
	httpMiddleware := func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, cn)
	}
	http.HandleFunc("/", httpMiddleware)

	//set backlog to nil
	backlog = nil

	if err != nil {
		log.Fatalf("Error: Unable to open TCP Connection: %s", err)
	}

	defer cn.Close()

	time.Sleep(2)
	print("Listener Started. Enter Commands:" + "\n")

	for {
		go tcpListen(cn)
		go http.ListenAndServe("127.0.0.1:8081", nil)
		send_cmd(cn)

	}
}

//------HTTP Handler Function-----//
func handler(w http.ResponseWriter, r *http.Request, cn net.Conn) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "GET" {
		b, err := json.Marshal(backlog)
		if err != nil {
			log.Fatalln("error with backlog")
		}
		go tcpListen(cn)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		w.Write(b)

		//clear backlog
		backlog = nil
	} else if r.Method == "POST" {
		w.WriteHeader(http.StatusOK)

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}
		send(string(b), cn)
		go tcpListen(cn)

	} else if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		json.NewEncoder(w).Encode("OKOK")
		w.WriteHeader(http.StatusOK)
	}

}

//-----------//
func tcpListen(cn net.Conn) {
	message, err := bufio.NewReader(cn).ReadString('\n')

	if err != nil {
		log.Fatalf("Error from TCP Session: %s", err)
	}

	if len(encryption_key) != 0 && strings.Contains(message, "encrypted:") {
		output := decrypt((strings.Split(message, " ")[2])[10:], encryption_key)
		fmt.Print("->: " + output)
		backlog = append(backlog, output)
	} else {
		fmt.Print("->: " + message)
		backlog = append(backlog, message)
	}
}

func send(text string, cn net.Conn) {

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
		encrypted_text := args[0] + " encrypted:" + encrypt(text, encryption_key) + "\n"
		fmt.Fprintf(cn, encrypted_text)
	}

}

//------Send Command to TCP server-----//
func send_cmd(cn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	send(text, cn)

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
