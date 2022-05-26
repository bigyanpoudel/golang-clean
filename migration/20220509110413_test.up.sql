CREATE TABLE IF NOT EXISTS `test`(
    `id` BINARY(16) PRIMARY KEY NOT NULL,
   `title` VARCHAR(256) NOT NULL,
     `created_at`               DATETIME               NOT NULL,
    `updated_at`               DATETIME               NOT NULL,
    `deleted_at`               DATETIME
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;