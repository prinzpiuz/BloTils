package handlers

func notCurrentAdmin(userId int, sessionUserId int) bool {
	return userId != sessionUserId
}
