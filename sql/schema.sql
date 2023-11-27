DELETE DATABASE IF EXISTS pubsub; 
CREATE DATABASE pubsub; 

USE pubsub; 

DELETE TABLE IF EXISTS users; 

CREATE TABLE users(
    id INTEGER PRIMARY KEY AUTO_INCREMENT, 
    first_name VARCHAR(50) NOT NULL, 
    last_name VARCHAR(50) NOT NULL,
    email_address VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    deleted_at TIMESTAMP NULL, 
    merged_at TIMESTAMP NULL, 
    parent_user_id INT, 
    FOREIGN KEY (parent_user_id) REFERENCES users(id)
)