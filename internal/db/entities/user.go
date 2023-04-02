package entities

type User struct {
	ID        int
	Login     string
	Password  string
	Balance   float32
	Withdrawn float32
}
