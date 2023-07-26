CREATE DATABASE IF NOT EXISTS chat_reports;
USE chat_reports;

DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS companies;

CREATE TABLE reports (
    id INT PRIMARY KEY AUTO_INCREMENT,
    operator_name VARCHAR(255),
    operator_id INT,
    dialog_id INT,
    tmr_in_seconds INT,
    opened_dialogs INT,
    client VARCHAR(255),
    status_tag VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT uc_dialog_id UNIQUE (dialog_id)
) ENGINE=INNODB;

CREATE TABLE companies (
    id INT PRIMARY KEY AUTO_INCREMENT,
    company_id INT,
    company_name VARCHAR(255),
    company_token VARCHAR(255),
    email_admin VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=INNODB;