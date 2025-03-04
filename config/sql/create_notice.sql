CREATE TABLE `notice` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `text` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;