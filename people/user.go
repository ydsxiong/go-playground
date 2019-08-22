package people

type User struct {
	Uid      int    `gorm:"primary_key"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func CreateUser(uid int, username, password, email string) *User {
	return &User{uid, username, password, email}
}
