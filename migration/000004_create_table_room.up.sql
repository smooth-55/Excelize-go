CREATE TABLE IF NOT EXISTS `rooms` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `name` varchar(50) NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NULL,
    `deleted_at` DATETIME NULL,
    PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS `messages` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `room_id` INT NOT NULL,
    `message_by` INT NOT NULL,
    `message` LONGTEXT NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NULL,
    `deleted_at` DATETIME NULL,
    PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS `room_users` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `room_id` INT NOT NULL,
    `user_id` INT NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NULL,
    `deleted_at` DATETIME NULL,
    PRIMARY KEY (id),
    CONSTRAINT `UQ_room_id_user_id` UNIQUE (`room_id`, `user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
