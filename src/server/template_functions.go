package server

func notCurrentAdmin(userId int, sessionUserId int) bool {
	return userId != sessionUserId
}

func add(a, b int) int {
	return a + b
}
