ALTER TABLE `users`
    ADD COLUMN `username` varchar(25) NOT NULL,
    ADD COLUMN `is_private` tinyint(4) DEFAULT 0,
    ADD CONSTRAINT `UQ_username` UNIQUE (`username`);