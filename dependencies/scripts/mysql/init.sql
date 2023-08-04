DROP DATABASE IF EXISTS messaging;
CREATE DATABASE messaging;
USE messaging;

CREATE TABLE rooms (
    id INT NOT NULL AUTO_INCREMENT,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,

    uuid VARCHAR(50) NOT NULL,
    created_at_nano DOUBLE NOT NULL,
    PRIMARY KEY (id)
) ENGINE=INNODB;

CREATE TABLE messages (
    id INT NOT NULL AUTO_INCREMENT,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    
    uuid VARCHAR(50) NOT NULL,
    message_text TEXT NOT NULL,
    from_uuid VARCHAR(50) NOT NULL,
    room_uuid VARCHAR(50) NOT NULL,
    room_id INT,
    message_status TINYTEXT NOT NULL,
    created_at_nano DOUBLE NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (room_id)
        REFERENCES rooms(id)
) ENGINE=INNODB;

CREATE TABLE seen_by (
    id INT NOT NULL AUTO_INCREMENT,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    
    message_uuid VARCHAR(50) NOT NULL,
    user_uuid VARCHAR(50) NOT NULL,
    message_id INT,
    PRIMARY KEY (id),
    UNIQUE KEY `user_message_uuids` (`user_uuid`,`message_uuid`),
    FOREIGN KEY (message_id)
        REFERENCES messages(id)
) ENGINE=INNODB;

CREATE TABLE members (
    id int NOT NULL UNIQUE AUTO_INCREMENT,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,

    uuid VARCHAR(50) NOT NULL,
    room_uuid VARCHAR(50) NOT NULL,
    room_id INT,
    user_uuid VARCHAR(50) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (room_id)
        REFERENCES rooms(id)
) ENGINE=INNODB;

CREATE TABLE auth_profiles (
    id int NOT NULL UNIQUE AUTO_INCREMENT,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,

    uuid VARCHAR(50) NOT NULL,
    hashed_password VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL
) ENGINE=INNODB;
