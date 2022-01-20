# GOChat
GOChat is an encrypted messaging app that allows users to create and host their own chat rooms or join existing chat rooms in a local network. Messages sent through the GOChat server can be encrypted using [AES-256](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) - the premier of `Advanced Encryption Standard`.

## Structure 

The usage of GOChat is divided up based on the server and user. Users have the option directly connect with the GOChat server directly using TCP. TCP connections are not recommended and the listener.go program was created to fix some of the usability issues related with raw TCP connections. The listener.go program supports client-side encryption and decryption—this meaning that the server host is unable to read these messages without the appropriate encryption key. 

Alternatively, the user can use to the server using a GUI interface instead of a command line. The GOChat GUI further abstracts the connection process and improves discoverability of features. The GUI sends HTTP requests to the listener.go file which relays the requests to the server. In future, technologies like Web Sockets can be used to improve the frontend preformance.  

![Alt text](https://media3.giphy.com/media/fSgLXRd90Jk0woZMnj/giphy.gif?cid=790b761177990e962ce38d3d612a2b8dc9f48bea343261e0&rid=giphy.gif&ct=g)



## Usage

The following instructions assumes the user has installed go, npm, and is running the program on Windows machine. Please refer to the [go documentation](https://go.dev/doc/install) and [Node Js documentation](https://nodejs.org/en/download/) for more information on installing the dependencies required to run the project.   

### Running GOChat Server 
To start the server:
1. Open up the `GoChat\Scripts` folder in a terminal
2. Start the server program using:
    ```Powershell
    .\startGoServer.ps1
    ```

### Using GOChat With Telnet
To connect with the GOChat server using telnet:
1. Open up a terminal and start a telnet connection using:
    ```ps1
    telnet localhost 8080
    ```
Press enter to start the connection. After starting up the program, we are given a variety of features that we can accessed these being: 

    '/create <name>'     create a room and join it
    '/join <name>'       join a room with associated room name
    '/nick <name>'       set name, otherwise name will be "stranger"
    '/exit'              disconnects from chat server 
    '/rooms'             show list of available rooms to join
    '/msg <msg>'         broadcasts message to everyone in a room
    '/encrypt'           generate AES-256 encryption key for messages
    '/encrypt <key>'     assign AES-256 encryption key for messages


### Using GOChat With Listener
To startup the listener:
1. Open up the `GoChat\Scripts` folder in a terminal
2. Start the listener program using:
    ```Powershell
    .\startGoListener.ps1
    ```
This will run the software in the command-line in port 8080 (can be changed in sourcecode). As before, we are given a variety of features that we can accessed these being: 

    '/create <name>'     create a room and join it
    '/join <name>'       join a room with associated room name
    '/nick <name>'       set name, otherwise name will be "stranger"
    '/exit'              disconnects from chat server 
    '/rooms'             show list of available rooms to join
    '/msg <msg>'         broadcasts message to everyone in a room
    '/encrypt'           generate AES-256 encryption key for messages
    '/encrypt <key>'     assign AES-256 encryption key for messages

Note: When using `/encrypt` the encryption key will be set and displayed for use.

### Using GOChat With GUI
To start Go GUI:
1. Open up the `GoChat\Scripts` folder in a terminal
2. Start the listener program using:
    ```Powershell
    .\startGoListener.ps1
    ```
3. In a new terminal window navigate to the `GoChat\Scripts` folder
4. Start the frontend development server using:
    ```Powershell
    .\startDevFrontend.ps1
    ```

The user 

