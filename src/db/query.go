// Package db provides functionality for interacting with the database.
package db

// getDomain is a SQL query that selects all rows from the Domain and DomainSettings tables
// where the domain column in the Domain table matches the provided parameter.
// The query joins the two tables on the id column.
const getDomain = `SELECT
					 d.id,
					 d.settings_id,
					 d.domain,
					 d.created_time,
					 ds.id,
					 ds.likes,
					 ds.comments,
					 ds.created_time
				   FROM
					 Domain d
				   JOIN DomainSettings ds ON d.id = ds.id
				   WHERE
					 d.domain = ?`

// Check if domain exists
const domainExists = `SELECT COUNT(*) FROM Domain WHERE domain = ?`

// getDomainById is a SQL query that selects all rows from the Domain and DomainSettings tables
// where the id column in the Domain table matches the provided parameter.
// The query joins the two tables on the id column, retrieving domain details and settings.
const getDomainById = `SELECT
					 d.id,
					 d.settings_id,
					 d.domain,
					 d.created_time,
					 ds.id,
					 ds.likes,
					 ds.comments,
					 ds.created_time
				   FROM
					 Domain d
				   JOIN DomainSettings ds ON d.id = ds.id
				   WHERE
					 d.id = ?`

// getAllUserDomains is a SQL query that selects all domains for a specific user,
// retrieving the domain's ID, settings ID, domain name, and creation time.
// The results are ordered by creation time in descending order, showing the most
// recently created domains first.
const getAllUserDomains = `SELECT
						 d.id,
						 d.settings_id,
						 d.domain,
						 d.created_time,
						 ds.likes,
						 ds.comments,
						 ds.created_time
					   FROM
					     Domain d
					   JOIN DomainSettings ds ON d.id = ds.id
					   WHERE
					     d.user_id = ?
					   ORDER BY
					   	 d.created_time
					   DESC`

// Delete related likes first
const deleteLikesByDomainID = `DELETE FROM Likes WHERE domain_id = ?`

// Delete related liked IPs by domain name
const deleteLikedIPsByDomain = `DELETE FROM Liked_IPs WHERE domain = ?`

// Get domain name before deletion (for cleaning up Liked_IPs)
const getDomainNameByID = `SELECT domain FROM Domain WHERE id = ?`

// Get settings_id before deletion
const getSettingsIDByDomainID = `SELECT settings_id FROM Domain WHERE id = ?`

// Delete domain
const deleteDomainByID = `DELETE FROM Domain WHERE id = ?`

// Delete domain settings
const deleteDomainSettingsByID = `DELETE FROM DomainSettings WHERE id = ?`

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
				 FROM
				   Likes
				 JOIN Domain ON Likes.domain_id = Domain.id
				 WHERE
				  Domain.domain = ?
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

// Insert domain settings and return the ID
const insertDomainSettings = `
    INSERT INTO DomainSettings (likes, comments, created_time)
    VALUES (?, ?, datetime())
`

// Insert domain with settings_id reference
const insertDomain = `
    INSERT INTO Domain (settings_id, user_id, domain, created_time)
    VALUES (?, ?, ?, datetime())
`

const updateDomainDetails = `BEGIN TRANSACTION;
							 UPDATE Domain
							 SET domain = ?
							 WHERE id = ? AND user_id = ?;
							 UPDATE DomainSettings
							 SET likes =?, comments =?
                             WHERE id = (SELECT settings_id FROM Domain WHERE Domain.id =? AND Domain.user_id =?);
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

const getUserWithId = `SELECT * FROM User WHERE id = ?`

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

// getTokenUser is a SQL query that retrieves password reset token details along with associated user information.
// The query selects the token, user ID, expiry time, usage status, user ID, and email from the PasswordResetToken
// and User tables. It joins the tables on user ID and filters for an unused token, allowing verification
// of a valid, unused password reset token for a specific user.
//
//nolint:gosec
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
       				 AND prt.used = 0`

// getSessionUser is a SQL query that retrieves comprehensive user and session details
// by joining the Sessions and User tables. The query selects session ID, user information,
// and session metadata, filtered by a specific session ID. It allows retrieving full
// user context associated with an active session.
const getSessionUser = `SELECT
						 s.session_id,
						 u.id,
						 u.email,
						 u.user_role,
						 u.is_active,
						 u.user_status,
						 u.created_time,
						 s.user_id,
						 s.created_time,
						 s.expiry_time
						FROM
						 Sessions s
						JOIN User u ON s.user_id = u.id
						WHERE
						 s.session_id = ?`

// createSession is a SQL query that inserts a new session record into the Sessions table
// with the provided session ID, user ID, and expiration time. The query uses placeholders
// to allow dynamic session creation for user authentication and tracking.
const createSession = `INSERT INTO Sessions (session_id, user_id, expiry_time)
					   VALUES (?, ?, ?)`

// deleteUserSessions is a SQL query that removes all session records associated with a specific user ID
// from the Sessions table. This query is typically used when invalidating or cleaning up a user's
// active sessions, such as during logout or account management operations.
const deleteUserSessions = `DELETE FROM Sessions
					   WHERE user_id = ?`

// deleteSession is a SQL query that removes a specific session record from the Sessions table
// by matching the provided session ID. This query is used to invalidate or terminate
// an individual user session, typically during logout or session expiration processes.
const deleteSession = `DELETE FROM Sessions
					   WHERE session_id = ?`

// getAllusers is a SQL query that retrieves a comprehensive list of all users
// from the User table, sorted by their creation time in descending order.
// The query returns user details including ID, email, role, active status,
// user status, and creation timestamp, providing a full overview of users.
const getAllusers = `SELECT
					  u.id,
					  u.email,
					  u.user_role,
					  u.is_active,
					  u.user_status,
					  u.created_time
					FROM
					  User u
					ORDER BY
					  created_time DESC`

// approveUser is a SQL query that updates a user's status to approved by setting their user role,
// active status, and user status to 1 (typically indicating an active and approved account state)
// for a specific user identified by their unique ID. This query is used in user management processes
// to grant full access and confirm a user's account.
const approveUser = `UPDATE USER
					 SET user_role = 2,
						 is_active = 1,
						 user_status = 1
					WHERE  USER.id = ?
					AND USER.user_status = 2
					RETURNING email`

// deleteUser is a SQL query that removes a specific user record from the USER table
// by matching the provided user ID. This query is used to permanently delete
// a user from the system, typically during user management or account removal processes.
const deleteUser = `DELETE FROM USER WHERE id = ?;`

// getLikesByDomainId retrieves all likes for a specific domain by domain ID.
const getLikesByDomainId = `SELECT
                             Likes.id,
                             Likes.uri,
                             Likes.domain_id,
                             Likes.count
                           FROM Likes
                           WHERE Likes.domain_id = ?
                           ORDER BY Likes.count DESC`

// getLikesTimeline retrieves the timeline of likes for a specific URI and domain,
// grouped by date, to power the likes-over-time graph.
const getLikesTimeline = `SELECT
                           DATE(Liked_IPs.created_time) as like_date,
                           COUNT(*) as like_count
                          FROM Liked_IPs
                          JOIN Domain ON Liked_IPs.domain = Domain.domain
                          WHERE Domain.id = ?
                          AND Liked_IPs.path = ?
                          GROUP BY DATE(Liked_IPs.created_time)
                          ORDER BY like_date ASC`
