DROP TABLE IF EXISTS `publish_rules`;

CREATE TABLE `publish_rules` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `rule_name` varchar(32) NOT NULL COMMENT '策略名称',
  `rule_type` int(32) NOT NULL DEFAULT '1' COMMENT '匹配方式:1:百分比，2：组织架构与成员，3：动态',
  `rule_status` int(2) NOT NULL DEFAULT '1' COMMENT '规则状态, 0：删除，1：正常，2：关闭状态',
  `percent` int(5) DEFAULT '1' COMMENT '放量百分比',
  `entry` varchar(500) NOT NULL DEFAULT '' COMMENT '策略对应的入口地址',
  `creater` varchar(32) NOT NULL DEFAULT '',
  `updater` varchar(32) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `publish_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `rule_name` (`rule_name`),
  KEY `publish_id` (`publish_id`),
  CONSTRAINT `publish_rules_ibfk_1` FOREIGN KEY (`publish_id`) REFERENCES `publishes` (`id`) ON DELETE
  SET
    NULL ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

/* 百分比 */
INSERT INTO
  `publish_rules` (
    `rule_name`,
    `rule_type`,
    `rule_status`,
    `percent`,
    `entry`,
    `creater`,
    `updater`,
    `created_at`,
    `updated_at`,
    `publish_id`
  )
VALUES
  (
    '灰度10%',
    1,
    1,
    10,
    'http://localhost:8080/html/c.html',
    'jht',
    'jht',
    NOW(),
    NOW(),
    1
  );

/* 指定人员 */
INSERT INTO
  `publish_rules` (
    `rule_name`,
    `rule_type`,
    `rule_status`,
    `percent`,
    `entry`,
    `creater`,
    `updater`,
    `created_at`,
    `updated_at`,
    `publish_id`
  )
VALUES
  (
    '指定人员',
    2,
    1,
    0,
    'http://localhost:8080/html/somebody.html',
    'jht',
    'jht',
    NOW(),
    NOW(),
    1
  );