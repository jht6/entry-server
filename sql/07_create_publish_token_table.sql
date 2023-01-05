DROP TABLE IF EXISTS `publish_tokens`;

CREATE TABLE `publish_tokens` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增长id',
  `token` varchar(100) NOT NULL COMMENT '发布者的token',
  `landun_token` varchar(100) DEFAULT NULL COMMENT '发布者的蓝盾token',
  `creater` varchar(32) NOT NULL DEFAULT '',
  `updater` varchar(32) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `token` (`token`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;