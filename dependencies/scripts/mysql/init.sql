DROP DATABASE IF EXISTS messaging;
CREATE DATABASE messaging;
USE messaging;

CREATE TABLE rooms (
    id INT NOT NULL AUTO_INCREMENT,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,

    uuid VARCHAR(36) NOT NULL,
    PRIMARY KEY (id)
) ENGINE=INNODB;

CREATE TABLE messages (
    id INT NOT NULL AUTO_INCREMENT,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    
    uuid VARCHAR(36) NOT NULL,
    message_text TEXT NOT NULL,
    from_uuid VARCHAR(36) NOT NULL,
    room_uuid VARCHAR(36) NOT NULL,
    room_id INT,
    PRIMARY KEY (id),
    FOREIGN KEY (room_id)
        REFERENCES rooms(id)
) ENGINE=INNODB;

CREATE TABLE members (
    id int NOT NULL UNIQUE AUTO_INCREMENT,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,

    uuid VARCHAR(36) NOT NULL,
    room_uuid VARCHAR(36) NOT NULL,
    room_id INT,
    user_uuid VARCHAR(36) NOT NULL,
    user_role VARCHAR(36) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (room_id)
        REFERENCES rooms(id)
) ENGINE=INNODB;
