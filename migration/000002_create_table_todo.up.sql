CREATE TABLE IF NOT EXISTS `Todo` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `task` VARCHAR(255) NULL,
  `is_completed` tinyint(4) DEFAULT 0,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NULL,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;