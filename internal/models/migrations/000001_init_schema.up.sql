-- schema.up.sql

-- Create users table
CREATE TABLE `users` (
 `id` int NOT NULL AUTO_INCREMENT,
 `name` varchar(255) NOT NULL,
 `email` varchar(255) NOT NULL,
 `hashed_password` char(60) NOT NULL,
 `created` datetime NOT NULL,
 PRIMARY KEY (`id`),
 UNIQUE KEY `users_uc_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create habits table
CREATE TABLE `habits` (
 `id` int NOT NULL AUTO_INCREMENT,
 `title` varchar(100) NOT NULL,
 `created` datetime NOT NULL,
 `user_id` int NOT NULL,
 PRIMARY KEY (`id`),
 KEY `idx_habits_created` (`created`),
 KEY `fk_habits_users` (`user_id`),
 CONSTRAINT `fk_habits_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create sessions table
CREATE TABLE `sessions` (
 `token` char(43) NOT NULL,
 `data` blob NOT NULL,
 `expiry` timestamp(6) NOT NULL,
 PRIMARY KEY (`token`),
 KEY `sessions_expiry_idx` (`expiry`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create habit_logs table
CREATE TABLE `habit_logs` (
 `id` int NOT NULL AUTO_INCREMENT,
 `habit_id` int NOT NULL,
 `date` date NOT NULL,
 `is_completed` tinyint(1) NOT NULL DEFAULT '0',
 `user_id` int NOT NULL,
 PRIMARY KEY (`id`),
 UNIQUE KEY `uniq_habit_date` (`habit_id`,`date`),
 KEY `idx_habit_logs_habit_id` (`habit_id`),
 KEY `user_id` (`user_id`),
 CONSTRAINT `fk_habit_logs_habits` FOREIGN KEY (`habit_id`) REFERENCES `habits` (`id`),
 CONSTRAINT `habit_logs_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT ON habittracker.* TO 'web'@'localhost';
ALTER USER 'web'@'localhost' IDENTIFIED BY 'bacon';
