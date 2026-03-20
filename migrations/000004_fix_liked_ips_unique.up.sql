-- Drop the old table and recreate with correct unique constraint
DROP TABLE IF EXISTS Liked_IPs;
CREATE TABLE Liked_IPs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip VARCHAR(255) NOT NULL,
    count INTEGER DEFAULT 1 CHECK(count >= 0),
    domain VARCHAR(255) NOT NULL,
    path TEXT NOT NULL,
    created_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(ip, domain, path) -- Unique per IP + domain + path combination
);
CREATE INDEX idx_liked_ips_lookup ON Liked_IPs(ip, domain, path);
-- Drop and recreate Likes table with correct unique constraint
DROP TABLE IF EXISTS Likes;
CREATE TABLE Likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uri TEXT NOT NULL,
    count INTEGER DEFAULT 0 CHECK(count >= 0),
    domain_id INTEGER NOT NULL,
    created_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(uri, domain_id),
    FOREIGN KEY (domain_id) REFERENCES Domain(id) ON DELETE CASCADE
);
CREATE INDEX idx_likes_lookup ON Likes(uri, domain_id);
