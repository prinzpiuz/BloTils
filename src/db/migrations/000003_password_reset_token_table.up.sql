CREATE TABLE IF NOT EXISTS PasswordResetToken (
    token TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expiry_time DATETIME NOT NULL,
    used INTEGER DEFAULT 0 check(
        used = 0
        or used = 1
    ),
    created_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User(id)
);
