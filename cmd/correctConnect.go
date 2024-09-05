package cmd

import (
	"TCPChat/cmd/models"
	"net"
	"strings"
)

func correctConnect(conn net.Conn, userCh chan *User) (*User, error) {
	// Спраишиваем имя
	// Проверяем имя на символы
	// Добавляем в юзеры и оздаем сообщение с 1ым байтом значения 2
	//
	adressStr := conn.RemoteAddr().String()

	welcomMes := strings.Join([]string{
		"Welcome to TCP-Chat!",
		"         _nnnn_",
		"        dGGGGMMb",
		"       @p~qp~~qMb",
		"       M|@||@) M|",
		"       @,----.JM|",
		"      JS^\\__/  qKL",
		"     dZP        qKRb",
		"    dZP          qKKb",
		"   fZP            SMMb",
		"   HZM            MMMM",
		"   FqM            MMMM",
		" __| \".        |\\dS\"qML",
		" |    `.       | `' \\Zq",
		"_)      \\.___.,|     .'",
		"\\____   )MMMMMP|   .'",
		"     `-'       `--'",
		"[ENTER YOUR NAME] (min-1 max-1024):",
	}, "\n")

	_, err := conn.Write([]byte(welcomMes))
	if err != nil {
		return &User{}, err
	}

	for {

		buf := make([]byte, 1024)
		namelen, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				return &User{}, models.ErrUserBrokeConn
			}
			return &User{}, err
		}

		chosenName := string(buf[:namelen-1])

		err = userNameChecker(chosenName)

		if err != nil {
			if err == models.ErrNameIsEmpety {

				_, err = conn.Write([]byte("Name is empety. Chose another name: "))
				if err != nil {
					return &User{}, err
				}

			} else if err == models.ErrNameExistsInServer {

				_, err = conn.Write([]byte("Name is exist in Chat. Chose another name: "))
				if err != nil {
					return &User{}, err
				}
			} else if err == models.ErrNameHasIllegalSims {

				_, err = conn.Write([]byte("Name has illegal simbols. Chose another name: "))
				if err != nil {
					return &User{}, err
				}
			}

			continue

		}

		if connQuantControl() {

			thisUser := NewUser(adressStr, chosenName)
			userCh <- thisUser

			return thisUser, nil
		} else {

			_, err = conn.Write([]byte("The server is full. Try again later"))
			if err != nil {
				return &User{}, err
			}

			return &User{}, models.ErrServerIsFull

		}

	}
}

func connQuantControl() bool {
	Mutex.Lock()
	users := ThisServer.Users
	Mutex.Unlock()

	connQuant := 0
	for i := 0; i < len(users); i++ {

		if users[i].ConnTime.After(users[i].LeftTime) {
			connQuant++
		}
	}

	return connQuant < 10
}

func userNameChecker(chosenName string) error {

	empetyStr := true
	illegalSim := false

	for i := 0; i < len(chosenName); i++ {

		if empetyStr && (chosenName[i] > ' ' && chosenName[i] <= '~') {
			empetyStr = false
		}

		if !illegalSim && (chosenName[i] < ' ' || chosenName[i] > '~') {
			illegalSim = true
		}
	}

	if illegalSim {
		return models.ErrNameHasIllegalSims
	} else if empetyStr {
		return models.ErrNameIsEmpety
	}

	Mutex.Lock()
	defer Mutex.Unlock()
	for i := 0; i < len(ThisServer.Users); i++ {

		user := ThisServer.Users[i]
		if user.LeftTime.Before(user.ConnTime) && user.Name == chosenName {
			return models.ErrNameExistsInServer
		}
	}

	return nil
}
