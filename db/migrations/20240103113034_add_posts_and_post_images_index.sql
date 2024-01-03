-- +goose Up
CREATE INDEX index_posts_on_user_id ON `posts` (user_id);
CREATE INDEX index_post_image_on_post_id ON `post_images` (post_id);

-- +goose Down
ALTER TABLE `posts` DROP FOREIGN KEY `posts_ibfk_1`;
DROP INDEX index_posts_on_user_id ON `posts`;
ALTER TABLE `posts` ADD CONSTRAINT `posts_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `post_images` DROP FOREIGN KEY `post_images_ibfk_1`;
DROP INDEX index_post_image_on_post_id ON post_images;
ALTER TABLE `post_images` ADD CONSTRAINT post_images_ibfk_1 FOREIGN KEY (`post_id`) REFERENCES `posts` (id);
