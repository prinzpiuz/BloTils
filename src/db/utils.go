// Package db provides utility functions for interacting with the database.
// this page contains the utility functions for db package
package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

var (
	ErrDomainExists   = errors.New("domain already exists")
	ErrDomainNotFound = errors.New("domain not found")
)

// GetDomain retrieves a Domain from the database by the given domain_name.
// If the domain is not found, it returns an empty Domain.
func GetDomain(db *sql.DB, domainName string) Domain {
	var domain Domain
	err := db.QueryRow(getDomain, domainName).Scan(
		&domain.Id,
		&domain.Settings.Id,
		&domain.Domain,
		&domain.timestamp,
		&domain.Settings.Id,
		&domain.Settings.likes,
		&domain.Settings.comments,
		&domain.Settings.timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: Domain %s, Not Found", domainName)
			return Domain{}
		}
		log.Printf("Error Getting Domain %s: %v", domainName, err)
		return Domain{}
	}
	return domain
}

// GetDomainById retrieves a Domain from the database by the given domain ID.
// If the domain is not found, it returns an empty Domain.
func GetDomainById(db *sql.DB, domainId int) (*Domain, error) {
	var domain Domain
	err := db.QueryRow(getDomainById, domainId).Scan(
		&domain.Id,
		&domain.Settings.Id,
		&domain.Domain,
		&domain.timestamp,
		&domain.Settings.Id,
		&domain.Settings.likes,
		&domain.Settings.comments,
		&domain.Settings.timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: Domain %d, Not Found", domainId)
			return nil, ErrDomainNotFound
		}
		log.Printf("Error Getting Domain %d: %v", domainId, err)
		return nil, fmt.Errorf("failed to get domain by ID: %w", err)
	}
	return &domain, nil
}

// GetAllUserDomains retrieves all domains associated with a specific user from the database.
// It takes a database connection and a user ID as input, and returns a slice of Domain structs.
// If no domains are found or an error occurs during the database query, it returns nil.
func GetAllUserDomains(db *sql.DB, userID int) []Domain {
	var domains []Domain

	rows, err := db.Query(getAllUserDomains, userID)
	if err == sql.ErrNoRows {
		log.Print("No Domains Found")
		return nil
	} else if err != nil {
		log.Printf("Error Querying Database: %v", err)
		return nil
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		log.Printf("Error Closing Rows: %v", err)
		return nil
	}
	for rows.Next() {
		var domain Domain
		err := rows.Scan(&domain.Id,
			&domain.Settings.Id,
			&domain.Domain,
			&domain.timestamp,
			&domain.Settings.likes,
			&domain.Settings.comments,
			&domain.Settings.timestamp)
		if err != nil {
			log.Printf("Scaning Rows Failed: %v", err)
			return nil
		}
		domains = append(domains, domain)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Iteration Failed: %v", err)
		return nil
	}
	return domains
}

// DeleteDomain deletes a domain and all related data in a transaction
func DeleteDomain(db *sql.DB, domainID int, userID int) error {
	conn := GetDB()
	if conn == nil {
		return errors.New("database not connected")
	}

	// Verify domain exists and belongs to user
	domain, err := GetDomainById(db, domainID)
	if err != nil {
		if errors.Is(err, ErrDomainNotFound) {
			return ErrDomainNotFound
		}
		return fmt.Errorf("failed to verify domain ownership: %w", err)
	}

	// Start transaction
	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. Delete related Likes
	_, err = tx.Exec(deleteLikesByDomainID, domainID)
	if err != nil {
		return fmt.Errorf("failed to delete likes: %w", err)
	}

	// 2. Delete related Liked_IPs (by domain name)
	_, err = tx.Exec(deleteLikedIPsByDomain, domain.Domain)
	if err != nil {
		return fmt.Errorf("failed to delete liked IPs: %w", err)
	}

	// 3. Delete the domain
	_, err = tx.Exec(deleteDomainByID, domainID)
	if err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}

	// 4. Delete domain settings
	_, err = tx.Exec(deleteDomainSettingsByID, domain.Settings.Id)
	if err != nil {
		return fmt.Errorf("failed to delete domain settings: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Domain '%s' (ID: %d) deleted successfully", domain.Domain, domainID)
	return nil
}

// GetLikes retrieves the likes information for a specific page within a given domain from the database.
// If no likes are found for the specified page and domain, it returns an empty Likes struct.
// Returns a Likes struct containing details such as like count, URI, and domain information.
func GetLikes(db *sql.DB, domainName string, page string) Likes {
	var likes Likes
	err := db.QueryRow(getLikes, domainName, page).Scan(
		&likes.id,
		&likes.URI,
		&likes.domain_id,
		&likes.Count,
		&likes.Domain.Id,
		&likes.Domain.Settings.Id,
		&likes.Domain.Domain,
		&likes.Domain.timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: Likes On %s, Not Found In %s", page, domainName)
			return Likes{}
		}
		log.Printf("Error Getting Likes On %s For %s: %v", page, domainName, err)
		return Likes{}
	}
	return likes
}

// GetLikedIP retrieves the liked IP information for a specific IP address within a given domain and page from the database.
// If no liked IP is found for the specified IP, domain, and page, it returns an empty LikedIPs struct.
// Returns a LikedIPs struct containing details such as IP, like count, domain, and path.
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

// UpdateIPLikeCount updates or inserts the like count for a specific IP address within a given domain and page.
// It executes a database query to record the IP's interaction and logs any errors encountered during the process.
func UpdateIPLikeCount(db *sql.DB, domain string, path string, ip string) {
	_, err := db.Exec(updateOrInsertLikedIP, domain, path, ip)
	if err != nil {
		log.Printf("Error Updating IP Like Count: %v", err)
	}

}

// UpdateLikeCount updates the like count for a specific page within a given domain.
// It executes a database query to increment or update the like count and returns any error encountered during the process.
// The function takes a database connection, page identifier, and domain ID as parameters.
func UpdateLikeCount(db *sql.DB, page string, doamin_id int) error {
	_, err := db.Exec(updateLike, page, doamin_id)
	return err
}

// DomainExists checks if a domain already exists
func DomainExists(db *sql.DB, domainName string) (bool, error) {

	var count int
	err := db.QueryRow(domainExists, domainName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// AddDomainAndSettings inserts a domain and its settings in a transaction
func AddDomainAndSettings(db *sql.DB, domain *Domain) error {

	// Check if domain already exists
	exists, err := DomainExists(db, domain.Domain)
	if err != nil {
		return fmt.Errorf("failed to check domain existence: %w", err)
	}
	if exists {
		return ErrDomainExists
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert domain settings
	result, err := tx.Exec(insertDomainSettings, domain.Settings.likes, domain.Settings.comments)
	if err != nil {
		return fmt.Errorf("failed to insert domain settings: %w", err)
	}

	settingsID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get settings ID: %w", err)
	}

	// Insert domain
	_, err = tx.Exec(insertDomain, settingsID, domain.User.Id, domain.Domain)
	if err != nil {
		return fmt.Errorf("failed to insert domain: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Domain '%s' added successfully for user %d", domain.Domain, domain.User.Id)
	return nil
}

func UpdateDomainDetails(db *sql.DB, domain *Domain) {
	_, err := db.Exec(updateDomainDetails, domain.Domain, domain.Id, domain.User.Id, domain.Settings.likes, domain.Settings.comments, domain.Id, domain.User.Id)
	if err != nil {
		log.Printf("Error Updating Domain %s: %v", domain.Domain, err)
	}
}

// CreatSuperUser adds a new admin user to the database with the provided email and password hash.
// It executes a database query to insert the admin user and returns any error encountered during the process.
func CreatSuperUser(db *sql.DB, email string, passwordHash string) error {
	_, err := db.Exec(addAdminUser, email, passwordHash)
	if err != nil {
		log.Printf("Error Adding Admin User %s: %v", email, err)
		return err
	}
	return nil
}

// AddUserRequest adds a new user account creation request to the database with the provided email and password hash.
// It executes a database query to insert the user request and returns any error encountered during the process.
func AddUserRequest(db *sql.DB, email string, passwordHash string) error {
	_, err := db.Exec(accountCreationRequest, email, passwordHash)
	if err != nil {
		log.Printf("Error Adding User Request %s: %v", email, err)
		return err
	}
	return nil
}

// GetUser retrieves a user from the database by their email address.
// It queries the database for a user with the given email and returns the user's details.
// If no user is found or an error occurs, it returns an empty User struct and logs the error.
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

// GetUserWithId retrieves a user from the database by their user ID.
// It queries the database for a user with the given ID and returns the user's details.
// If no user is found or an error occurs, it returns an empty User struct and logs the error.
func GetUserWithId(db *sql.DB, id int) User {
	var user User
	err := db.QueryRow(getUserWithId, id).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.userRole,
		&user.IsActive,
		&user.userStatus,
		&user.timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: User With Id %d, Not Found", id)
			return User{}
		}
		log.Printf("Error Getting User  With Id %d: %v", id, err)
		return User{}
	}
	return user
}

// SetRestToken saves a password reset token for a specific user with a 10-minute expiration time.
// It inserts the token into the database and returns any error encountered during the process.
func SetRestToken(db *sql.DB, token string, userId int) error {
	expiryTime := time.Now().Add(10 * time.Minute)
	_, err := db.Exec(addResetToken, token, userId, expiryTime)
	if err != nil {
		log.Printf("Error Saving Reset Token For %d: %v", userId, err)
		return err
	}
	return nil
}

// GetTokenUser retrieves a password reset token from the database by its token string.
// It queries the database for a token with the given token and returns the associated PasswordResetToken.
// If no token is found or an error occurs, it returns an empty PasswordResetToken struct and logs the error.
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

// DeleteTokenAndSetPassword invalidates a password reset token and updates the user's password in a single database transaction.
// It takes the reset token, new password hash, and user ID as parameters and executes a database query to complete the password reset process.
// Returns an error if the database operation fails, otherwise returns nil.
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

// DeleteSessionWithSessionId removes a specific session from the database using its session ID
// It takes a database connection and a session ID as parameters
// Returns an error if the database operation fails
func DeleteSessionWithSessionId(db *sql.DB, sessionId string) error {
	_, err := db.Exec(deleteSession, sessionId)
	if err != nil {
		log.Printf("Error Deleting Session %s: %v", sessionId, err)
		return err
	}
	return nil
}

// GetAllUsers retrieves all users from the database
// Returns a slice of User structs containing user information
// Returns nil if no users are found or if a database error occurs
func GetAllUsers(db *sql.DB) []User {
	var users []User
	rows, err := db.Query(getAllusers)
	if err == sql.ErrNoRows {
		log.Print("No Users Found")
		return nil
	} else if err != nil {
		log.Printf("Error Querying Database: %v", err)
		return nil
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Println(err)
		}
	}()
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

// ApproveUser approves a pending user account and returns the user's email.
//
// It updates the user's record with:
// - user_role = 1 (admin)
// - is_active = 1 (active)
// - user_status = 1 (approved)
//
// The function only approves users with a current user_status of 2 (pending).
// If the user is not found or is already approved, it returns an error indicating
// the user was not found or already approved.
func ApproveUser(db *sql.DB, userID int) (string, error) {
	var email string
	err := db.QueryRow(approveUser, userID).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found or already approved")
		}
		return "", fmt.Errorf("failed to approve user: %w", err)
	}
	return email, nil
}

// DeleteUser removes a user from the database by their user ID
// It takes a database connection and a user ID as parameters
// Returns an error if the database operation fails
func DeleteUser(db *sql.DB, userId int) error {
	_, err := db.Exec(deleteUser, userId)
	if err != nil {
		return err
	}
	return nil
}

// GetLikesByDomainId retrieves all likes associated with a specific domain ID.
// Returns a slice of Likes structs ordered by count descending.
func GetLikesByDomainId(db *sql.DB, domainId int) []Likes {
	var likesList []Likes
	rows, err := db.Query(getLikesByDomainId, domainId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: No likes found for domain ID %d", domainId)
			return nil
		}
		log.Printf("Error Getting Likes for Domain %d: %v", domainId, err)
		return nil
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Println(err)
		}
	}()
	for rows.Next() {
		var like Likes
		err := rows.Scan(&like.id, &like.URI, &like.domain_id, &like.Count)
		if err != nil {
			log.Printf("Scanning Likes Rows Failed: %v", err)
			return nil
		}
		like.Domain.Id = domainId
		likesList = append(likesList, like)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Likes Iteration Failed: %v", err)
		return nil
	}
	return likesList
}

// GetLikesTimeline retrieves the likes timeline for a specific URI within a domain,
// grouped by date, for rendering a likes-over-time graph.
func GetLikesTimeline(db *sql.DB, domainId int, uri string) []LikesTimeline {
	var timeline []LikesTimeline
	rows, err := db.Query(getLikesTimeline, domainId, uri)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: No timeline data for domain %d, uri %s", domainId, uri)
			return nil
		}
		log.Printf("Error Getting Likes Timeline: %v", err)
		return nil
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Println(err)
		}
	}()
	for rows.Next() {
		var entry LikesTimeline
		err := rows.Scan(&entry.Date, &entry.Count)
		if err != nil {
			log.Printf("Scanning Timeline Rows Failed: %v", err)
			return nil
		}
		timeline = append(timeline, entry)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Timeline Iteration Failed: %v", err)
		return nil
	}
	return timeline
}
