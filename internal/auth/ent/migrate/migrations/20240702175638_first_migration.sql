-- Create "roles" table
CREATE TABLE `roles` (`id` text NOT NULL, `name` text NOT NULL, `permissions` json NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, PRIMARY KEY (`id`));
-- Create index "roles_name_key" to table: "roles"
CREATE UNIQUE INDEX `roles_name_key` ON `roles` (`name`);
-- Create "tokens" table
CREATE TABLE `tokens` (`id` text NOT NULL, `token` text NOT NULL, `type` text NOT NULL, `expires_at` datetime NOT NULL, `revoked` bool NOT NULL DEFAULT (false), `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `user_tokens` text NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `tokens_users_tokens` FOREIGN KEY (`user_tokens`) REFERENCES `users` (`id`) ON DELETE NO ACTION);
-- Create index "tokens_token_key" to table: "tokens"
CREATE UNIQUE INDEX `tokens_token_key` ON `tokens` (`token`);
-- Create "users" table
CREATE TABLE `users` (`id` text NOT NULL, `username` text NOT NULL, `email` text NOT NULL, `password` text NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, PRIMARY KEY (`id`));
-- Create index "users_username_key" to table: "users"
CREATE UNIQUE INDEX `users_username_key` ON `users` (`username`);
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX `users_email_key` ON `users` (`email`);
-- Create "user_roles" table
CREATE TABLE `user_roles` (`user_id` text NOT NULL, `role_id` text NOT NULL, PRIMARY KEY (`user_id`, `role_id`), CONSTRAINT `user_roles_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE, CONSTRAINT `user_roles_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE);
