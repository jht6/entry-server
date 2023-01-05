/* 发布项目主配置表 */
DROP TABLE IF EXISTS `t_publish`;

create table if not exists `t_publish` (
  `publish_id` int auto_increment primary key,
  `name` varchar(100) not null comment '项目名称',
  `domain` varchar(500) not null comment '域名',
  `entry` varchar(500) not null comment 'index.html文件地址',
  `status` tinyint default 0 comment '状态：0-正常 1-禁用 2-删除',
  `create_user` varchar(32) default null,
  `update_user` varchar(32) default null,
  `created_at` datetime default null,
  `updated_at` datetime default null
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 添加一条配置 */
insert into publishes (
  `domain`,
  `name`,
  `entry`,
  `creater`,
  `updater`,
  `created_at`,
  `updated_at`
) values (
  'jht1.woa.com',
  'jht1',
  'http://localhost:8080/html/a.html',
  'jht',
  'jht',
  NOW(),
  NOW()
);

insert into publishes (
  `domain`,
  `name`,
  `entry`,
  `creater`,
  `updater`,
  `created_at`,
  `updated_at`
) values (
  'jht2.woa.com',
  'jht2',
  'http://localhost:8080/html/b.html',
  'jht',
  'jht',
  NOW(),
  NOW()
);