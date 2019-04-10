package user

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var users = make(map[uint64]User)

func AddUser(u User) {
	log.Printf("AddUser %v", u)

	u.Registrationtime = time.Now().UTC().UnixNano()
	users[u.Token] = u

	PrintUsers()
}

func RemoveUser(token uint64) {
	log.Printf("RemoveUser %v", token)
	delete(users, token)

	PrintUsers()
}

func LookupUser(token uint64) bool {
	//log.Printf("Search for user token %v.", token)
	ue := users[token]
	if ue == (User{}) {
		return false
	}
	return true
}

func GetUser(token uint64) (u User) {
	//log.Printf("GetUser ue with token %v", token)
	user := users[token]

	return user
}

func PrintUsers() {
	log.Printf("#### Print user directory")
	for token, user := range users {
		log.Printf("\tFound user with token: [%v] \t and IP: [%v] -> user: [%v]", token, user.IP, user)
	}
	log.Printf("####")
}
