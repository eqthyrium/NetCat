package cmd

import "time"

type User struct {
	Adress   string
	Name     string
	ConnTime time.Time
	LeftTime time.Time
}

type Message struct {
	User               *User
	Message            string
	DateAndTimeMassage time.Time
}

type TCPServer struct {
	Users    []*User
	Messages []Message
}

func NewUser(adress, name string) *User {
	return &User{Adress: adress, Name: name, ConnTime: time.Now()}
}

func NewMessage(user User, massage string) Message {
	return Message{User: &user, Message: massage, DateAndTimeMassage: time.Now()}
}
