DROP TABLE IF EXISTS `t_rule`;

CREATE TABLE `t_rule` (
  `rule_id` int auto_increment primary key,
  `rule_name` varchar(32) not null COMMENT '策略名称',
  `rule_type` tinyint not null COMMENT '匹配方式: 1-百分比 2-指定成员',
  `rule_status` tinyint not null default 0 COMMENT '规则状态 0-启用 1-关闭 2-删除',
  `percent` tinyint DEFAULT 0 COMMENT '放量百分比',
  `entry` varchar(500) not null COMMENT '本策略的html入口',
  `creater` varchar(32) default null,
  `updater` varchar(32) default null,
  `created_at` datetime default null,
  `updated_at` datetime default null,
  `publish_id` int not null,
  PRIMARY KEY (`id`)
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