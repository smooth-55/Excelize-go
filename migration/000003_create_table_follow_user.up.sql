CREATE TABLE IF NOT EXISTS `follow_user` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `followed_by_id` INT NOT NULL,
    `followed_to_id` INT NOT NULL,
    `is_approved` tinyint DEFAULT 1,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NULL,
    `deleted_at` DATETIME NULL,
    PRIMARY KEY (id),
    CONSTRAINT `UQ_user_id_followed_by_to` UNIQUE (`followed_by_id`, `followed_to_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;