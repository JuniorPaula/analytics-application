CREATE DATABASE IF NOT EXISTS chat_reports;
USE chat_reports;

DROP TABLE IF EXISTS reports;

CREATE TABLE reports (
    id INT PRIMARY KEY AUTO_INCREMENT,
    operator_name VARCHAR(255),
    operator_id INT,
    dialog_id INT,
    tmr_in_seconds INT,
    opened_dialogs INT,
    client VARCHAR(255),
    status VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT uc_dialog_id UNIQUE (dialog_id),
    
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;