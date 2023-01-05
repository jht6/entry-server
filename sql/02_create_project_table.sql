create table if not exists `t_project` (
  `project_id` int(5) auto_increment primary key,
  `project_name` varchar(50) not null comment '项目名称',
  `host` varchar(200) not null comment '主机名',
  `html_url` varchar(200) not null comment 'html源文件地址',
  `is_delete` tinyint(1) default 0 comment '是否已删除 0=未删除 1=已删除',
  `create_user` int(10) not null,
  `created_at` datetime not null,
  `update_user` int(10) not null,
  `updated_at` datetime not null
) ENGINE=InnoDB DEFAULT CHARSET=utf8;