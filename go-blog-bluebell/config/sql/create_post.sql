create table post
(
    id           bigint auto_increment
        primary key,
    post_id      bigint                              not null comment '帖子id',
    title        varchar(128)                        not null comment '标题',
    content      varchar(8192)                       not null comment '内容',
    author_id    bigint                              not null comment '作者的用户id',
    community_id bigint                              not null comment '所属社区',
    status       tinyint   default 1                 not null comment '帖子状态',
    create_time  timestamp default CURRENT_TIMESTAMP null comment '创建时间',
    update_time  timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间',
    constraint idx_post_id
        unique (post_id)
)
    collate = utf8mb4_general_ci;

create index idx_author_id
    on post (author_id);

create index idx_community_id
    on post (community_id);

INSERT INTO bluebell_v2.post
(created_at, updated_at, deleted_at, status, post_id, community_id, author_id, title, content) VALUES
('2020-08-09 09:58:39', '2020-08-09 09:58:39', NULL, 1, 14283784123846656, 1, 28018727488323585, '学习使我快乐', '只有学习才能变得更强'),
('2020-08-09 15:53:40', '2020-08-09 15:53:40', NULL, 1, 14373128436191232, 2, 28018727488323585, 'CSGO开箱子好上瘾', '花了钱不出金，我好气啊'),
('2020-08-09 15:54:08', '2020-08-09 15:54:08', NULL, 1, 14373246019309568, 3, 28018727488323585, 'IG牛逼', '打得好啊。。。'),
('2020-08-23 14:58:29', '2020-08-23 14:58:29', NULL, 1, 19432670719119360, 2, 28018727488323585, '投票功能真好玩', '12345'),
('2020-08-23 15:02:37', '2020-08-23 15:02:37', NULL, 1, 19433711036534784, 2, 28018727488323585, '投票功能真好玩2', '12345'),
('2020-08-23 15:04:26', '2020-08-23 15:04:26', NULL, 1, 19434165682311168, 2, 28018727488323585, '投票功能真好玩2', '12345'),
('2020-08-30 04:27:23', '2020-08-30 04:27:23', NULL, 1, 21810561880690688, 2, 28018727488323585, '看图说话', '4321'),
('2020-08-30 04:27:52', '2020-08-30 04:27:52', NULL, 1, 21810685746876416, 3, 28018727488323585, '永远不要高估自己', '做个普通人也挺难'),
('2020-08-30 04:28:35', '2020-08-30 04:28:35', NULL, 1, 21810865955147776, 1, 28018727488323585, '你知道泛型是什么吗？', '不知道泛型是什么却一直在问泛型什么时候出'),
('2020-08-30 04:28:52', '2020-08-30 04:28:52', NULL, 1, 21810938202034176, 1, 28018727488323585, '国庆假期哪里玩？', '走遍四海，还是威海。');