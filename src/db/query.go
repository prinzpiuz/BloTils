// Package db provides functionality for interacting with the database.
package db

// getDomain is a SQL query that selects all rows from the Domain and DomainSettings tables
// where the domain column in the Domain table matches the provided parameter.
// The query joins the two tables on the id column.
const getDomain = `SELECT * FROM Domain
				   JOIN DomainSettings
				   ON Domain.id = DomainSettings.id
				   WHERE domain = ?`

// getLikes is a SQL query that selects all rows from the Likes table
// where the uri matches the provided value, and the domain matches the
// provided domain. It joins the Likes table with the Domain table to
// retrieve the domain information.
const getLikes = `SELECT
				  Likes.id,
				  Likes.uri,
				  Likes.domain_id,
				  Likes.count,
				  Domain.id,
				  Domain.settings_id,
				  Domain.domain,
				  Domain.created_time
				  FROM Likes
				  JOIN Domain
				  ON Likes.domain_id = Domain.id
				  WHERE Domain.domain = ?
				  AND Likes.uri = ?`

// getIPlikedOrNot is a SQL query that selects all rows from the Liked_IPs table
// where the ip, path, and domain columns match the provided parameters. This query
// is used to check if a specific IP address has already been liked for a given
// domain and path.
const getIPlikedOrNot = `SELECT * FROM Liked_IPs
						 WHERE Liked_IPs.ip = ?
						 AND Liked_IPs.domain = ?
						 AND Liked_IPs.path = ?`

// updateOrInsertLikedIP is a SQL query that inserts a new row into the Liked_IPs table
// with the provided IP address, an initial count of 1, and the current timestamp. If a
// row already exists for the provided IP address, the query updates the existing row by
// incrementing the Count column by 1.
const updateOrInsertLikedIP = `INSERT INTO Liked_IPs(domain, path, ip, count, created_time)
							   VALUES(?, ?, ?, 1, datetime())
							   ON CONFLICT(ip)
							   DO UPDATE
							   SET count = count + 1`

// updateLike is a SQL query that inserts a new row into the Likes table with the provided uri and domain_id,
// and an initial count of 1. If a row already exists for the provided uri and domain_id, the query updates
// the existing row by incrementing the count column by 1.
const updateLike = `INSERT INTO Likes(uri, count, domain_id)
					VALUES (?, 1, ?)
					ON CONFLICT(uri)
					DO UPDATE
					SET count = count + 1`

// addDomainAndSettings is a SQL transaction that atomically inserts a new domain and its associated settings.
// It first creates a record in the DomainSettings table with likes, comments, and a timestamp,
// then uses the last inserted row ID to create a corresponding Domain record with the settings ID,
// domain name, and timestamp. The transaction ensures that both insertions are completed successfully.
const addDomainAndSettings = `BEGIN TRANSACTION;
							  INSERT INTO DomainSettings (likes, comments, created_time)
							  VALUES (?, ?, datetime());
							  INSERT INTO Domain (settings_id, domain, created_time)
							  VALUES (last_insert_rowid(), ?, datetime());
							  COMMIT TRANSACTION;`

// addAdminUser is a SQL query that inserts a new admin user into the User table
// with predefined active status and admin role. The query uses placeholders for
// email and password hash, allowing dynamic user creation with admin privileges.
const addAdminUser = `INSERT INTO User(email, password_hash, user_role, is_active, user_status)
					  VALUES (?, ?, 1, 1, 1)`

// accountCreationRequest is a SQL query that inserts a new user into the User table
// with a standard user role, inactive status, and a pending user status. The query
// uses placeholders for email and password hash, allowing dynamic user creation
// with default non-admin settings.
const accountCreationRequest = `INSERT INTO User(email, password_hash, user_role, is_active, user_status)
					  VALUES (?, ?, 2, 0, 2)`

// getUser is a SQL query that retrieves all columns for a User record matching the specified email address.
// The query uses a parameterized input to safely select a user by their unique email identifier.
const getUser = `SELECT * FROM User WHERE email = ?`

// addResetToken is a SQL query that inserts a new password reset token into the PasswordResetToken table
// with the provided token, user ID, and expiration time. The query uses placeholders to allow
// dynamic token creation for password reset functionality.
const addResetToken = `INSERT INTO PasswordResetToken (token, user_id, expiry_time)
					   VALUES (?, ?, ?)`

// resetTokenAndUpdatePassword is a SQL transaction that atomically deletes a password reset token
// and updates the user's password. The transaction ensures that the token is removed and the
// password is updated in a single, consistent operation, preventing partial updates in case
// of an error during the process.
const resetTokenAndUpdatePassword = `BEGIN TRANSACTION;
									 DELETE FROM PasswordResetToken
									 WHERE token = ?;
									 UPDATE User
									 SET password_hash = ?
									 WHERE id = ?;
									 COMMIT TRANSACTION;`

const getTokenUser = `SELECT prt.token,
							prt.user_id,
							prt.expiry_time,
							prt.used,
							u.id,
							u.email
					 FROM  PasswordResetToken prt
       				 JOIN USER u
         			 ON u.id = prt.user_id
					 WHERE  prt.token = ?
       				 AND prt.used = 0 `
