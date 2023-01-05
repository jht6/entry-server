DROP TABLE IF EXISTS `publish_rule_orgs`;

CREATE TABLE `publish_rule_orgs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `org_id` int(11) NOT NULL COMMENT '组织Id',
  `name` varchar(100) NOT NULL DEFAULT '1' COMMENT '组织架构名称',
  `full_name` varchar(500) NOT NULL DEFAULT '1' COMMENT '组织架构全称',
  `creater` varchar(32) NOT NULL DEFAULT '',
  `updater` varchar(32) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `publish_rule_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB AUTO_INCREMENT = 2 DEFAULT CHARSET = utf8;
