-- 用户测试数据
INSERT INTO user (user_id, user_name, password, age, gender, email, avatar) VALUES
(10001, 'testuser1', '123456', 20, '男', 'test1@example.com', '/static/avatar.png'),
(10002, 'testuser2', 'abcdef', 22, '女', 'test2@example.com', '/static/avatar.png');

-- 社区测试数据
INSERT INTO community (community_id, community_name, introduction) VALUES
(1, '编程', '编程技术交流'),
(2, '树洞', '树洞社区');

-- 帖子测试数据
INSERT INTO post (post_id, title, content, author_id, community_id) VALUES
(20001, 'Hello World', '这是我的第一篇帖子', 10001, 1),
(20002, '树洞发言', '匿名树洞内容', 10002, 2); 