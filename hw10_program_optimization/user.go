package hw10programoptimization

//go:generate easyjson -all user.go
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}
