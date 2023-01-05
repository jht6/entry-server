DROP TABLE IF EXISTS `publish_rule_users`;

CREATE TABLE `publish_rule_users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '用户Id',
  `chinese_name` varchar(50) NOT NULL COMMENT '中文名称',
  `english_name` varchar(50) NOT NULL COMMENT '英文名称',
  `full_name` varchar(100) NOT NULL COMMENT '全称',
  `creater` varchar(32) NOT NULL DEFAULT '',
  `updater` varchar(32) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `publish_rule_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB AUTO_INCREMENT = 3 DEFAULT CHARSET = utf8;

INSERT INTO
  `publish_rule_users` (
    `user_id`,
    `chinese_name`,
    `english_name`,
    `full_name`,
    `creater`,
    `updater`,
    `created_at`,
    `updated_at`,
    `publish_rule_id`
  )
VALUES
  (
    285144,
    'jht',
    'jht',
    'jht',
    '',
    '',
    NOW(),
    NOW(),
    2
  );
