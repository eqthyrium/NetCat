package cmd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var (
	ThisServer TCPServer
	Mutex      sync.Mutex
)

func FromMain(args []string) error {
	var err error

	port := "8989"
	if len(args) == 2 {
		port = args[1]
	} else if len(args) > 2 {
		return errors.New("[USAGE]: ./TCPChat $port")
	}

	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return errors.New("[USAGE]: ./TCPChat $port")
	}
	log.Printf("Listening on the port :"+port+":\n\n	nc localhost %v\n\n", port)
	defer server.Close()

	userCh := make(chan *User)
	go userReceiver(userCh)

	messageCh := make(chan Message)
	go messageReceiver(messageCh)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go connHandler(conn, userCh, messageCh)
	}
}

// iota

func messageReceiver(messageCh chan Message) {
	for {

		message := <-messageCh

		Mutex.Lock()
		ThisServer.Messages = append(ThisServer.Messages, message)
		Mutex.Unlock()
	}
}

func userReceiver(userCh chan *User) {
	for {
		user := <-userCh

		Mutex.Lock()
		ThisServer.Users = append(ThisServer.Users, user)
		Mutex.Unlock()
	}
}

func connHandler(conn net.Conn, userCh chan *User, messageCh chan Message) {
	thisUser, err := correctConnect(conn, userCh)
	if err != nil {
		conn.Close()
		return
	}

	/*
		create notification for another
	*/
	connectMessage := NewMessage(*thisUser, string(byte(2))+"'"+thisUser.Name+"' Connected to this Chat")
	messageCh <- connectMessage

	/*
		go listener messages from others
	*/

	go listenMessageFromOtherUsers(conn)

	/*
		go listener messages form user
		if eof user diconected >> create message
	*/

	go listenMessageFromThisUser(conn, thisUser, messageCh)
}

func listenMessageFromOtherUsers(conn net.Conn) {
	massagesPrinted := make([]Message, 0)

	for {
		Mutex.Lock()
		nowMessages := ThisServer.Messages
		Mutex.Unlock()
		if len(nowMessages)-1 > len(massagesPrinted)-1 {

			forPrint := nowMessages[len(massagesPrinted):]
			massagesPrinted = nowMessages

			massagePrinter(forPrint, conn)

		}
	}
}

func listenMessageFromThisUser(conn net.Conn, thisUser *User, messageCh chan Message) {
	for {
		buf := make([]byte, 1024)
		messageLen, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {

				Mutex.Lock()
				(*thisUser).LeftTime = time.Now()
				Mutex.Unlock()

				disconnectMessage := NewMessage(*thisUser, string(byte(2))+"'"+thisUser.Name+"' Disconnected from this Chat")

				messageCh <- disconnectMessage

				conn.Close()

			}
			return
		}

		if messageLen > 0 {
			messageStr := string(buf[:messageLen-1])

			empetyMess := true
			illegalSim := false

			for i := 0; i < len(messageStr); i++ {

				if empetyMess && (messageStr[i] > ' ' && messageStr[i] <= '~') {
					empetyMess = false
				}

				if !illegalSim && (messageStr[i] < ' ' || messageStr[i] > '~') {
					illegalSim = true
				}
			}

			if !empetyMess && !illegalSim {

				conn.Write([]byte("\033[F"))
				for i := 0; i < len(messageStr); i++ {
					conn.Write([]byte(" "))
				}
				conn.Write([]byte("\r"))

				messageMess := NewMessage(*thisUser, messageStr)
				messageCh <- messageMess

				continue
			} else {
				conn.Write([]byte("Error: bad message\n"))
			}

		}

	}
}

func massagePrinter(forPrint []Message, conn net.Conn) {
	for i := 0; i < len(forPrint); i++ {
		message := forPrint[i]
		messageOwner := *message.User
		messageTime := message.DateAndTimeMassage
		messageText := message.Message
		formatedTime := messageTime.Format("2006-01-02 15:04:05")

		if message.Message[0] == 2 { // if the first character of the message text is the 2nd byte, this is a message about connecting or disconnecting a user
			messageString := fmt.Sprintf("[%v] %v\n", formatedTime, messageText[1:])
			_, err := conn.Write([]byte(messageString))
			if err != nil {
				conn.Write([]byte(messageString))
			}

		} else {
			messageString := fmt.Sprintf("[%v] %v: %v\n", formatedTime, messageOwner.Name, messageText)
			_, err := conn.Write([]byte(messageString))
			if err != nil {
				conn.Write([]byte(messageString))
			}
		}

	}
}
