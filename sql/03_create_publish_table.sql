/* 主配置表 */
create table if not exists `publishes` (
  `id` int(5) auto_increment primary key,
  `domain` varchar(500) not null comment '域名',
  `name` varchar(32) not null comment '项目名称',
  `status` int(11) default 0 comment '状态',
  `entry` varchar(500) not null comment 'index.html文件地址',
  `creater` varchar(32) not null,
  `updater` varchar(32) not null,
  `created_at` datetime not null,
  `updated_at` datetime not null
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