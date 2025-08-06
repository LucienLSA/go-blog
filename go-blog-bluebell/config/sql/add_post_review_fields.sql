-- 添加帖子审核相关字段
ALTER TABLE post 
ADD COLUMN review_status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending' COMMENT '审核状态',
ADD COLUMN review_result TEXT COMMENT '审核结果详情',
ADD COLUMN reviewed_at TIMESTAMP NULL COMMENT '审核时间',
ADD COLUMN review_score DECIMAL(3,2) DEFAULT 0.00 COMMENT '审核评分',
ADD COLUMN review_reason VARCHAR(500) COMMENT '拒绝原因',
ADD COLUMN review_suggestions TEXT COMMENT '修改建议',
ADD COLUMN review_tags JSON COMMENT '内容标签';

-- 创建审核状态索引
CREATE INDEX idx_review_status ON post (review_status);
CREATE INDEX idx_reviewed_at ON post (reviewed_at);

-- 创建审核任务表
CREATE TABLE post_review_task (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    post_id BIGINT NOT NULL COMMENT '帖子ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    content TEXT NOT NULL COMMENT '审核内容',
    status ENUM('pending', 'processing', 'completed', 'failed') DEFAULT 'pending' COMMENT '任务状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    error_message TEXT COMMENT '错误信息',
    INDEX idx_post_id (post_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) COLLATE = utf8mb4_general_ci;

-- 创建审核日志表
CREATE TABLE post_review_log (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    post_id BIGINT NOT NULL COMMENT '帖子ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    action VARCHAR(50) NOT NULL COMMENT '操作类型',
    old_status VARCHAR(20) COMMENT '原状态',
    new_status VARCHAR(20) COMMENT '新状态',
    review_result TEXT COMMENT '审核结果',
    operator_id BIGINT COMMENT '操作者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_post_id (post_id),
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
) COLLATE = utf8mb4_general_ci; 