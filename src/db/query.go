// Package db provides functionality for interacting with the database.
package db

// getDomainQuery is a SQL query that selects all rows from the Domain and DomainSettings tables
// where the domain column in the Domain table matches the provided parameter.
// The query joins the two tables on the id column.
const getDomainQuery string = `SELECT * FROM Domain
							   JOIN DomainSettings
							   ON Domain.id = DomainSettings.id
							   WHERE domain = ?`

// getLikesQuery is a SQL query that selects all rows from the Likes table
// where the uri matches the provided value, and the domain matches the
// provided domain. It joins the Likes table with the Domain table to
// retrieve the domain information.
const getLikesQuery string = `SELECT * FROM Likes
							  JOIN Domain
							  ON Likes.domain_id = Domain.id
							  where Domain.domain = ?
							  AND Likes.uri = ?`
