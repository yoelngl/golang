CREATE TABLE IF NOT EXISTS users (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    image VARCHAR(255) NULL,
    created_time TIMESTAMP NOT NULL,
    updated_time TIMESTAMP NULL,
    UNIQUE KEY unique_email (email)
)ENGINE=innoDB;