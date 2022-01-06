# GOChat
GOChat is a private secure messaging app that allows for users to create their on individual parties as well as servers to send information all the while having an emphasis on anonymity, allowing users to generate their very own AES-256—this being the primier of the `Advanced Encryption Standard` 

## Usage
The usage of GOChat is divided up based on the server and user. Users use listener.go file which supports client-side encryption and decryption—this meaning that the server host is unable to read these messages without the encryption key.

### Client
To startup the listener, use:
```go
go run listener.go
```
This will run the software in the command-line in port 8080 (can be changed in sourcecode). After starting up the program, we are given a variety of features that we can accessed these being: 

    '/create <name>'     create a room and join it
    '/join <name>'       join a room with associated room name
    '/nick <name>'       set name, otherwise name will be "stranger"
    '/exit'              disconnects from chat server 
    '/rooms'             show list of available rooms to join
    '/msg <msg>'         broadcasts message to everyone in a room
    '/encrypt'           generate AES-256 encryption key for messages
    '/encrypt <key>'     assign AES-256 encryption key for messages

Note: When using `/encrypt` the encryption key will be set and displayed for use. 

### Server
The server can be activated by compiling the classes followed by running it. This can be done in one command in GO by using:
```go
go run .
```
Running this command in the relative directory will compile and run the program in port 8080 and allows users to connect via the listener. 
