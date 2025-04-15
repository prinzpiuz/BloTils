CREATE TABLE IF NOT EXISTS User (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    user_role NUMERIC check(
        user_role = 1
        or user_role = 2
    ),
    is_active NUMERIC check(
        is_active = 0
        or is_active = 1
    ),
    user_status NUMERIC check(
        user_status = 1
        or user_status = 2
        or user_status = 3
    ),
    created_time DATETIME DEFAULT CURRENT_TIMESTAMP
);
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
    user_id INTEGER,
    domain VARCHAR(255) UNIQUE,
    created_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (settings_id) REFERENCES DomainSettings(id),
    FOREIGN KEY (user_id) REFERENCES User(id) ON DELETE CASCADE
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
