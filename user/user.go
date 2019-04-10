package user

type User struct {
	Token            uint64 `json:"token"`
	IP               string `json:"ip"`
	Registrationtime int64  `json:"registrationtime"`
}
