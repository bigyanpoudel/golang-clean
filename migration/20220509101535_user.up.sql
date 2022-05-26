CREATE TABLE IF NOT EXISTS `user`(
    `id` BINARY(16) PRIMARY KEY NOT NULL,
    `uuid` VARCHAR(256) NOT NULL,
    `name` VARCHAR(60) NOT NULL,
    `address` VARCHAR(60) NOT NULL,
    `email` VARCHAR(100) UNIQUE NOT NULL,
    `verified` BOOLEAN,
    `role` VARCHAR(30) Not Null,
     `created_at`               DATETIME               NOT NULL,
    `updated_at`               DATETIME               NOT NULL,
    `deleted_at`               DATETIME
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;