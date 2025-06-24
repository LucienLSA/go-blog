create table community
(
    id             bigint(20) auto_increment
        primary key,
    community_id   bigint(20) unsigned                        not null,
    community_name varchar(128)                        not null,
    introduction   varchar(256)                        not null,
    create_time    timestamp default CURRENT_TIMESTAMP not null,
    update_time    timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    constraint idx_community_id
        unique (community_id),
    constraint idx_community_name
        unique (community_name)
)
    collate = utf8mb4_general_ci;

INSERT INTO bluebell_v2.community 
(created_at, updated_at, deleted_at, community_id, community_name, introduction) 
VALUES 
('2016-11-01 08:10:10', '2016-11-01 08:10:10', NULL, 1, 'Go', 'Golang'),
('2020-01-01 08:00:00', '2020-01-01 08:00:00', NULL, 2, 'leetcode', '刷题刷题刷题'),
('2018-08-07 08:30:00', '2018-08-07 08:30:00', NULL, 3, 'CS:GO', 'Rush B。。。'),
('2016-01-01 08:00:00', '2016-01-01 08:00:00', NULL, 4, 'LOL', '欢迎来到英雄联盟!');