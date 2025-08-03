/*
 Navicat Premium Data Transfer

 Source Server         : MySQL5.7
 Source Server Type    : MySQL
 Source Server Version : 50743 (5.7.43-log)
 Source Host           : localhost:3306
 Source Schema         : bluebell_v2

 Target Server Type    : MySQL
 Target Server Version : 50743 (5.7.43-log)
 File Encoding         : 65001

 Date: 03/08/2025 16:27:49
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for community
-- ----------------------------
DROP TABLE IF EXISTS `community`;
CREATE TABLE `community`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `community_id` bigint(20) NULL DEFAULT NULL,
  `community_name` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `introduction` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_community_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of community
-- ----------------------------
INSERT INTO `community` VALUES (1, '2016-11-01 08:10:10', '2016-11-01 08:10:10', NULL, 1, 'Go', 'Golang');
INSERT INTO `community` VALUES (2, '2020-01-01 08:00:00', '2020-01-01 08:00:00', NULL, 2, 'leetcode', '刷题刷题刷题');
INSERT INTO `community` VALUES (3, '2018-08-07 08:30:00', '2018-08-07 08:30:00', NULL, 3, 'CS:GO', 'Rush B。。。');
INSERT INTO `community` VALUES (4, '2016-01-01 08:00:00', '2016-01-01 08:00:00', NULL, 4, 'LOL', '欢迎来到英雄联盟!');
INSERT INTO `community` VALUES (5, '2025-06-24 06:57:29', '2025-06-24 06:57:29', NULL, 5, '最优控制', '最优控制理论');

-- ----------------------------
-- Table structure for notice
-- ----------------------------
DROP TABLE IF EXISTS `notice`;
CREATE TABLE `notice`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `text` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_notice_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of notice
-- ----------------------------
INSERT INTO `notice` VALUES (1, '2025-06-23 15:22:05', NULL, NULL, '您正在绑定邮箱Email');
INSERT INTO `notice` VALUES (2, '2025-06-23 15:22:19', NULL, NULL, '您正在解绑邮箱Email');
INSERT INTO `notice` VALUES (3, '2025-08-03 14:00:52', NULL, NULL, '您正在登录邮箱Email');

-- ----------------------------
-- Table structure for post
-- ----------------------------
DROP TABLE IF EXISTS `post`;
CREATE TABLE `post`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `status` int(11) NULL DEFAULT NULL,
  `post_id` bigint(20) NULL DEFAULT NULL,
  `community_id` bigint(20) NULL DEFAULT NULL,
  `author_id` bigint(20) NULL DEFAULT NULL,
  `title` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `content` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_post_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of post
-- ----------------------------
INSERT INTO `post` VALUES (1, '2020-08-09 09:58:39', '2020-08-09 09:58:39', NULL, 1, 14283784123846656, 1, 63105569170395136, '学习使我快乐', '只有学习才能变得更强');
INSERT INTO `post` VALUES (2, '2020-08-09 15:53:40', '2020-08-09 15:53:40', NULL, 1, 14373128436191232, 2, 63105569170395136, 'CSGO开箱子好上瘾', '花了钱不出金，我好气啊');
INSERT INTO `post` VALUES (3, '2020-08-09 15:54:08', '2020-08-09 15:54:08', NULL, 1, 14373246019309568, 3, 63105569170395136, 'IG牛逼', '打得好啊。。。');
INSERT INTO `post` VALUES (4, '2020-08-23 14:58:29', '2020-08-23 14:58:29', NULL, 1, 19432670719119360, 2, 63105569170395136, '投票功能真好玩', '12345');
INSERT INTO `post` VALUES (5, '2020-08-23 15:02:37', '2020-08-23 15:02:37', NULL, 1, 19433711036534784, 2, 63105569170395136, '投票功能真好玩2', '12345');
INSERT INTO `post` VALUES (6, '2020-08-23 15:04:26', '2020-08-23 15:04:26', NULL, 1, 19434165682311168, 2, 63105569170395136, '投票功能真好玩2', '12345');
INSERT INTO `post` VALUES (7, '2020-08-30 04:27:23', '2020-08-30 04:27:23', NULL, 1, 21810561880690688, 2, 63105569170395136, '看图说话', '4321');
INSERT INTO `post` VALUES (8, '2020-08-30 04:27:52', '2020-08-30 04:27:52', NULL, 1, 21810685746876416, 3, 63105569170395136, '永远不要高估自己', '做个普通人也挺难');
INSERT INTO `post` VALUES (9, '2020-08-30 04:28:35', '2020-08-30 04:28:35', NULL, 1, 21810865955147776, 1, 63105569170395136, '你知道泛型是什么吗？', '不知道泛型是什么却一直在问泛型什么时候出');
INSERT INTO `post` VALUES (10, '2020-08-30 04:28:52', '2020-08-30 04:28:52', NULL, 1, 21810938202034176, 1, 63105569170395136, '国庆假期哪里玩？', '走遍四海，还是威海。');
INSERT INTO `post` VALUES (11, '2025-06-24 07:18:21', '2025-06-24 07:18:21', NULL, 1, 63165801196163072, 5, 63105569170395136, '基于ADP的最优控制', ' 结合强化学习。。。');
INSERT INTO `post` VALUES (12, '2025-06-24 15:32:46', '2025-06-24 15:32:46', NULL, 1, 14283784123846656, 1, 28018727488323585, '学习使我快乐', '只有学习才能变得更强');
INSERT INTO `post` VALUES (13, '2025-06-24 15:32:46', '2025-06-24 15:32:46', NULL, 1, 14373128436191232, 2, 28018727488323585, 'CSGO开箱子好上瘾', '花了钱不出金，我好气啊');
INSERT INTO `post` VALUES (14, '2025-06-24 15:32:46', '2025-06-24 15:32:46', NULL, 1, 14373246019309568, 3, 28018727488323585, 'IG牛逼', '打得好啊。。。');
INSERT INTO `post` VALUES (15, '2025-06-24 15:32:46', '2025-06-24 15:32:46', NULL, 1, 19432670719119360, 2, 28018727488323585, '投票功能真好玩', '12345');
INSERT INTO `post` VALUES (16, '2025-08-03 06:32:57', '2025-08-03 06:32:57', NULL, 1, 77649893493051392, 3, 77395559970770944, 'CSGO已经回不去了', ' CS2时代到来');

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `user_id` bigint(20) NULL DEFAULT NULL,
  `age` tinyint(3) UNSIGNED NULL DEFAULT NULL,
  `user_name` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `email` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `password_digest` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `gender` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `token` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `avatar` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uni_user_user_name`(`user_name`) USING BTREE,
  INDEX `idx_user_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, '2025-08-02 13:42:20', '2025-08-03 05:59:09', NULL, 77395559970770944, 23, 'lucien', '1020263522@qq.com', '$2a$12$SJhwMh0/8XOxAv56HicuBuz7A3zl0VkfZ/52tjnSR0P9jNgzX8qda', '男', '', 'http://t0elomjvk.hd-bkt.clouddn.com/avatar/lucien_1754200623873951000.jpg');
INSERT INTO `user` VALUES (2, '2025-08-03 14:35:40', '2025-08-03 14:35:42', NULL, 28018727488323585, 21, 'ggg', NULL, '$2a$12$SJhwMh0/8XOxAv56HicuBuz7A3zl0VkfZ/52tjnSR0P9jNgzX8qda', '女', NULL, 'http://t0elomjvk.hd-bkt.clouddn.com/avatar/lucien_1754200623873951000.jpg');
INSERT INTO `user` VALUES (3, '2025-08-03 14:36:14', '2025-08-03 14:36:16', NULL, 63105569170395136, 11, 'qqq', NULL, '$2a$12$SJhwMh0/8XOxAv56HicuBuz7A3zl0VkfZ/52tjnSR0P9jNgzX8qda', '男', NULL, 'http://t0elomjvk.hd-bkt.clouddn.com/avatar/lucien_1754200623873951000.jpg');

SET FOREIGN_KEY_CHECKS = 1;
