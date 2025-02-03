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
