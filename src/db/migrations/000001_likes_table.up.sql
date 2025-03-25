CREATE TABLE IF NOT EXISTS DomainSettings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    likes NUMERIC check(
        likes = 0
        or likes = 1
    ),
    comments NUMERIC check(
        comments = 0
        or comments = 1
    ),
    created_time DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS Domain (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    settings_id INTEGER,
    domain VARCHAR(255) UNIQUE,
    created_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (settings_id) REFERENCES DomainSettings(id)
);
CREATE TABLE IF NOT EXISTS Liked_IPs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip VARCHAR(255) UNIQUE,
    count INTEGER check(count >= 0),
    domain VARCHAR(255),
    path TEXT,
    created_time DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS Likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uri TEXT UNIQUE,
    count INTEGER check(count >= 0),
    domain_id INTEGER,
    FOREIGN KEY (domain_id) REFERENCES Domain(id)
);
