CREATE TABLE urls (
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id      BIGINT NOT NULL,
    short_code   VARCHAR(10) NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    hits         BIGINT DEFAULT 0,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_short_code (short_code)
) ENGINE=InnoDB;