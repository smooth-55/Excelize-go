ALTER TABLE `users`
    ADD COLUMN `username` varchar(25) NOT NULL,
    ADD CONSTRAINT `UQ_username` UNIQUE (`username`);