CREATE TABLE `messages` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `created_at` DATETIME NULL DEFAULT NOW(),
    `sender_id` INT NULL,
    `recipient_id` INT NULL,
    `message` VARCHAR(4096) NULL,
    PRIMARY KEY (`id`));
CREATE INDEX sender_id_idx ON messages (sender_id);
CREATE INDEX recipient_id_idx ON messages (recipient_id);