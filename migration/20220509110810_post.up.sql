CREATE TABLE IF NOT EXISTS `post`(
    `id` BINARY(16) PRIMARY KEY NOT NULL,
    `user_id` BINARY(16) NOT NULL,
    `title` VARCHAR(256) NOT NULL,
    `description` VARCHAR(256)  NOT NULL,
    `image` VARCHAR(256),
    `created_at`               DATETIME               NOT NULL,
    `updated_at`               DATETIME               NOT NULL,
    `deleted_at`               DATETIME,
    FOREIGN KEY (`user_id`)
      REFERENCES `user`(`id`)

) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;