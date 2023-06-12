DROP DATABASE IF EXISTS messaging;
CREATE DATABASE messaging;
USE messaging;

CREATE TABLE chat_messages (
    id int NOT NULL UNIQUE AUTO_INCREMENT PRIMARY KEY,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    
    uuid VARCHAR(36) NOT NULL UNIQUE,
    message_text TEXT NOT NULL,
    from_uuid VARCHAR(36) NOT NULL,
    room_uuid VARCHAR(36) NOT NULL
);

CREATE TABLE chat_rooms (
    id int NOT NULL UNIQUE AUTO_INCREMENT PRIMARY KEY,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,

    uuid VARCHAR(36) NOT NULL UNIQUE
);

CREATE TABLE chat_participants (
    id int NOT NULL UNIQUE AUTO_INCREMENT PRIMARY KEY,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,

    uuid VARCHAR(36) NOT NULL UNIQUE,
    room_uuid VARCHAR(36) NOT NULL,
    user_uuid VARCHAR(36) NOT NULL
);
