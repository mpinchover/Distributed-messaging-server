DROP DATABASE IF EXISTS messaging;
CREATE DATABASE messaging;
USE messaging;

CREATE TABLE chat_messages (
    uuid VARCHAR(36) NOT NULL UNIQUE PRIMARY KEY,
    message_text TEXT NOT NULL,
    from_uuid VARCHAR(36) NOT NULL,
    room_uuid VARCHAR(36) NOT NULL
);

CREATE TABLE chat_rooms (
    uuid VARCHAR(36) NOT NULL UNIQUE PRIMARY KEY
);

CREATE TABLE chat_participants (
    uuid VARCHAR(36) NOT NULL UNIQUE PRIMARY KEY,
    room_uuid VARCHAR(36) NOT NULL,
    user_uuid VARCHAR(36) NOT NULL
);
