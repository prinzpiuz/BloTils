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
