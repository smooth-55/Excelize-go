ALTER TABLE `users`
    ADD COLUMN `uuid` VARCHAR(255) DEFAULT NULL,
    ADD CONSTRAINT `UQ_uuid` UNIQUE (`uuid`);


ALTER TABLE `Todo`
    ADD COLUMN `uuid` VARCHAR(255) DEFAULT NULL,
    ADD CONSTRAINT `UQ_uuid` UNIQUE (`uuid`);