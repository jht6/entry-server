DROP TABLE IF EXISTS `t_rule`;

CREATE TABLE `t_rule` (
  `rule_id` int auto_increment primary key,
  `name` varchar(32) not null COMMENT '策略名称',
  `type` tinyint not null COMMENT '匹配方式: 1-百分比 2-指定成员 3-指定header',
  `config` varchar(500) not null COMMENT '规则配置JSON',
  `status` tinyint not null default 0 COMMENT '规则状态 0-启用 1-关闭 2-删除',
  `entry` varchar(500) not null COMMENT '本策略的html入口',
  `create_user` varchar(32) default null,
  `update_user` varchar(32) default null,
  `created_at` datetime default null,
  `updated_at` datetime default null,
  `publish_domain` varchar(200) not null COMMENT '所属发布项目的域名'
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