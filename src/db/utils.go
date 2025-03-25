// Package db provides utility functions for interacting with the database.
// this page contains the utility functions for db package
package db

import (
	"database/sql"
	"log"
	"time"
)

// GetDomain retrieves a Domain from the database by the given domain_name.
// If the domain is not found, it returns an empty Domain.
func GetDomain(db *sql.DB, domain_name string) Domain {
	var domain Domain
	err := db.QueryRow(getDomain, domain_name).Scan(
		&domain.ID,
		&domain.settings.id,
		&domain.domain,
		&domain.timestamp,
		&domain.settings.id,
		&domain.settings.likes,
		&domain.settings.comments,
		&domain.settings.timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: Domain %s, Not Found", domain_name)
			return Domain{}
		}
		log.Printf("Error Getting Domain %s: %v", domain_name, err)
		return Domain{}
	}
	return domain
}

func GetLikes(db *sql.DB, domain_name string, page string) Likes {
	var likes Likes
	err := db.QueryRow(getLikes, domain_name, page).Scan(
		&likes.id,
		&likes.URI,
		&likes.domain_id,
		&likes.Count,
		&likes.Domain.ID,
		&likes.Domain.settings.id,
		&likes.Domain.domain,
		&likes.Domain.timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: Likes On %s, Not Found In %s", page, domain_name)
			return Likes{}
		}
		log.Printf("Error Getting Likes On %s For %s: %v", page, domain_name, err)
		return Likes{}
	}
	return likes
}

func GetLikedIP(db *sql.DB, domain string, page string, ip string) LikedIPs {
	var likedIP LikedIPs
	err := db.QueryRow(getIPlikedOrNot, ip, domain, page).Scan(
		&likedIP.id,
		&likedIP.IP,
		&likedIP.Count,
		&likedIP.Domain,
		&likedIP.Path,
		&likedIP.timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: IP %s, Not Found Liked in %s for %s", ip, domain, page)
			return LikedIPs{}
		}
		log.Printf("Error Getting LikedIPs On %s For %s: %v", ip, page, err)
		return LikedIPs{}
	}
	return likedIP
}

func UpdateIPLikeCount(db *sql.DB, domain string, path string, ip string) {
	_, err := db.Exec(updateOrInsertLikedIP, domain, path, ip)
	if err != nil {
		log.Printf("Error Updating IP Like Count: %v", err)
	}

}

func UpdateLikeCount(db *sql.DB, page string, doamin_id int) error {
	_, err := db.Exec(updateLike, page, doamin_id)
	return err
}

func AddDomainAndSettings(db *sql.DB, domain string, likes int, comments int) {
	_, err := db.Exec(addDomainAndSettings, likes, comments, domain)
	if err != nil {
		log.Printf("Error Adding Domain %s: %v", domain, err)
	}
}

func CreatSuperUser(db *sql.DB, email string, passwordHash string) error {
	_, err := db.Exec(addAdminUser, email, passwordHash)
	if err != nil {
		log.Printf("Error Adding Admin User %s: %v", email, err)
		return err
	}
	return nil
}

func AddUserRequest(db *sql.DB, email string, passwordHash string) error {
	_, err := db.Exec(accountCreationRequest, email, passwordHash)
	if err != nil {
		log.Printf("Error Adding User Request %s: %v", email, err)
		return err
	}
	return nil
}

func GetUser(db *sql.DB, email string) User {
	var user User
	err := db.QueryRow(getUser, email).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.userRole,
		&user.IsActive,
		&user.userStatus,
		&user.timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: User %s, Not Found", email)
			return User{}
		}
		log.Printf("Error Getting User %s: %v", email, err)
		return User{}
	}
	return user
}

func SetRestToken(db *sql.DB, token string, userId int) error {
	expiryTime := time.Now().Add(10 * time.Minute)
	_, err := db.Exec(addResetToken, token, userId, expiryTime)
	if err != nil {
		log.Printf("Error Saving Reset Token For %d: %v", userId, err)
		return err
	}
	return nil
}

func GetTokenUser(db *sql.DB, token string) PasswordResetToken {
	var tokenObj PasswordResetToken
	err := db.QueryRow(getTokenUser, token).Scan(&tokenObj.Token,
		&tokenObj.User.Id,
		&tokenObj.ExpiryTime,
		&tokenObj.Used,
		&tokenObj.User.Id,
		&tokenObj.User.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: Token %s, Not Found", token)
			return PasswordResetToken{}
		}
		log.Printf("Error Getting User By Token %s: %v", token, err)
		return PasswordResetToken{}
	}
	return tokenObj
}

func DeleteTokenAndSetPassword(db *sql.DB, token string, passwordHash string, userId int) error {
	_, err := db.Exec(resetTokenAndUpdatePassword, token, passwordHash, userId)
	if err != nil {
		log.Printf("Error Resetting Password For User %d err:%v", userId, err)
		return err
	}
	return nil
}

func GetSession(db *sql.DB, sessionId string) Session {
	var session Session
	err := db.QueryRow(getSessionUser, sessionId).Scan(
		&session.SessionId,
		&session.User.Id,
		&session.User.Email,
		&session.User.userRole,
		&session.User.IsActive,
		&session.User.userStatus,
		&session.User.timestamp,
		&session.User.Id,
		&session.timestamp,
		&session.Expiry,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: Session %s, Not Found", sessionId)
			return Session{}
		}
		log.Printf("Error Getting User By Session %s: %v", sessionId, err)
		return Session{}
	}
	return session
}

func CreateSession(db *sql.DB, sessionId string, userId int) error {
	expiryTime := time.Now().Add(24 * time.Hour) //todo: change to 30 days
	_, err := db.Exec(createSession, sessionId, userId, expiryTime)
	if err != nil {
		log.Printf("Error Adding Session %s For User %d: %v", sessionId, userId, err)
		return err
	}
	return nil
}

func DeleteSession(db *sql.DB, userId int) error {
	_, err := db.Exec(deleteUserSessions, userId)
	if err != nil {
		log.Printf("Error Deleting Session For User %d: %v", userId, err)
		return err
	}
	return nil
}

func DeleteSessionWithSessionId(db *sql.DB, sessionId string) error {
	_, err := db.Exec(deleteSession, sessionId)
	if err != nil {
		log.Printf("Error Deleting Session %s: %v", sessionId, err)
		return err
	}
	return nil
}

func GetAllUsers(db *sql.DB) []User {
	var users []User
	rows, err := db.Query(getAllusers)
	defer rows.Close()
	if err == sql.ErrNoRows {
		log.Print("No Users Found")
		return nil
	} else if err != nil {
		log.Printf("Error Querying Database: %v", err)
		return nil
	}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Email, &user.userRole, &user.IsActive, &user.userStatus, &user.timestamp)
		if err != nil {
			log.Printf("Scaning Rows Failed: %v", err)
			return nil
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Iteration Failed: %v", err)
		return nil
	}
	return users

}
