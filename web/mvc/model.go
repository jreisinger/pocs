package main

type User struct {
	ID   int
	Name string
}

var users = map[int]User{
	1: {ID: 1, Name: "Alice"},
	2: {ID: 2, Name: "Bob"},
}

func GetUserByID(id int) (User, bool) {
	user, exists := users[id]
	return user, exists
}
