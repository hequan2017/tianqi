-- ============================================
-- 模型训练模块 - 简化版SQL（不使用变量）
-- 适用于所有MySQL版本
-- 执行前请先备份数据库！
-- ============================================

-- 第一步：插入父菜单（模型训练）
-- 记录插入后的ID，用于后续子菜单
INSERT INTO `sys_base_menus` (`created_at`, `updated_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `keep_alive`, `default_menu`, `close_tab`)
VALUES (NOW(), NOW(), 0, 0, 'modeltraining', 'modeltraining', 0, 'view/routerHolder.vue', 2, '模型训练', 'data-analysis', 0, 0, 0);

-- 查看父菜单ID (假设返回的ID是 XXX，请记下实际ID)
SELECT id FROM sys_base_menus WHERE `name` = 'modeltraining' AND `parent_id` = 0;

-- ============================================
-- 第二步：插入子菜单
-- 请将下面的 @PARENT_ID@ 替换为上面查询到的父菜单ID
-- ============================================

-- 子菜单：数据集管理
INSERT INTO `sys_base_menus` (`created_at`, `updated_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `keep_alive`, `default_menu`, `close_tab`)
VALUES (NOW(), NOW(), 1, @PARENT_ID@, 'dataset', 'dataset', 0, 'view/modeltraining/dataset/index.vue', 1, '数据集管理', 'folder-opened', 0, 0, 0);

-- 子菜单：训练任务
INSERT INTO `sys_base_menus` (`created_at`, `updated_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `keep_alive`, `default_menu`, `close_tab`)
VALUES (NOW(), NOW(), 1, @PARENT_ID@, 'trainingTask', 'trainingTask', 0, 'view/modeltraining/training/index.vue', 2, '训练任务', 'video-play', 0, 0, 0);

-- 子菜单：创建训练任务（隐藏菜单）
INSERT INTO `sys_base_menus` (`created_at`, `updated_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `title`, `icon`, `keep_alive`, `default_menu`, `close_tab`)
VALUES (NOW(), NOW(), 1, @PARENT_ID@, 'createTrainingTask', 'createTrainingTask', 1, 'view/modeltraining/training/createTask.vue', 3, '创建训练任务', 'plus', 0, 0, 0);

-- 查看所有新菜单ID
SELECT id, name, path, parent_id FROM sys_base_menus WHERE `name` IN ('modeltraining', 'dataset', 'trainingTask', 'createTrainingTask');

-- ============================================
-- 第三步：插入菜单权限关联
-- 请将下面的菜单ID替换为实际ID
-- ============================================

-- 假设: modeltraining父菜单ID = M1, dataset ID = M2, trainingTask ID = M3, createTrainingTask ID = M4
-- 为超级管理员(888)添加菜单权限
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`, `sys_base_menu_id`) VALUES
(888, @M1@),
(888, @M2@),
(888, @M3@),
(888, @M4@);

-- ============================================
-- 第四步：插入API数据
-- ============================================

INSERT INTO `sys_apis` (`created_at`, `updated_at`, `path`, `description`, `api_group`, `method`) VALUES
-- 数据集管理 API
(NOW(), NOW(), '/modeltraining/dataset/createDataset', '创建数据集', '数据集管理', 'POST'),
(NOW(), NOW(), '/modeltraining/dataset/deleteDataset', '删除数据集', '数据集管理', 'DELETE'),
(NOW(), NOW(), '/modeltraining/dataset/deleteDatasetByIds', '批量删除数据集', '数据集管理', 'DELETE'),
(NOW(), NOW(), '/modeltraining/dataset/updateDataset', '更新数据集', '数据集管理', 'PUT'),
(NOW(), NOW(), '/modeltraining/dataset/findDataset', '查询数据集详情', '数据集管理', 'GET'),
(NOW(), NOW(), '/modeltraining/dataset/getDatasetList', '获取数据集列表', '数据集管理', 'GET'),
(NOW(), NOW(), '/modeltraining/dataset/getDatasetDataSource', '获取数据集数据源', '数据集管理', 'GET'),
(NOW(), NOW(), '/modeltraining/dataset/createVersion', '创建数据集版本', '数据集管理', 'POST'),
(NOW(), NOW(), '/modeltraining/dataset/deleteVersion', '删除数据集版本', '数据集管理', 'DELETE'),
(NOW(), NOW(), '/modeltraining/dataset/getVersionList', '获取数据集版本列表', '数据集管理', 'GET'),
(NOW(), NOW(), '/modeltraining/dataset/publishDataset', '发布数据集', '数据集管理', 'POST'),
(NOW(), NOW(), '/modeltraining/dataset/uploadFile', '上传数据集文件', '数据集管理', 'POST'),
(NOW(), NOW(), '/modeltraining/dataset/uploadVersionFile', '上传版本文件', '数据集管理', 'POST'),
-- 训练任务 API
(NOW(), NOW(), '/modeltraining/trainingTask/createTask', '创建训练任务', '训练任务', 'POST'),
(NOW(), NOW(), '/modeltraining/trainingTask/deleteTask', '删除训练任务', '训练任务', 'DELETE'),
(NOW(), NOW(), '/modeltraining/trainingTask/deleteTaskByIds', '批量删除训练任务', '训练任务', 'DELETE'),
(NOW(), NOW(), '/modeltraining/trainingTask/updateTask', '更新训练任务', '训练任务', 'PUT'),
(NOW(), NOW(), '/modeltraining/trainingTask/findTask', '查询训练任务详情', '训练任务', 'GET'),
(NOW(), NOW(), '/modeltraining/trainingTask/getTaskList', '获取训练任务列表', '训练任务', 'GET'),
(NOW(), NOW(), '/modeltraining/trainingTask/startTask', '启动训练任务', '训练任务', 'POST'),
(NOW(), NOW(), '/modeltraining/trainingTask/stopTask', '停止训练任务', '训练任务', 'POST'),
(NOW(), NOW(), '/modeltraining/trainingTask/markCompleted', '手动标记完成', '训练任务', 'POST'),
(NOW(), NOW(), '/modeltraining/trainingTask/startService', '启动推理服务', '训练任务', 'POST'),
(NOW(), NOW(), '/modeltraining/trainingTask/stopService', '停止推理服务', '训练任务', 'POST'),
(NOW(), NOW(), '/modeltraining/trainingTask/chatCompletion', '模型对话测试', '训练任务', 'POST'),
(NOW(), NOW(), '/modeltraining/trainingTask/getTaskLogs', '获取训练日志', '训练任务', 'GET'),
(NOW(), NOW(), '/modeltraining/trainingTask/getTaskDataSource', '获取训练任务数据源', '训练任务', 'GET'),
(NOW(), NOW(), '/modeltraining/trainingTask/getDefaultParams', '获取默认训练参数', '训练任务', 'GET'),
-- 模型测试历史 API
(NOW(), NOW(), '/modeltraining/modelTest/createTestHistory', '创建测试历史', '模型测试', 'POST'),
(NOW(), NOW(), '/modeltraining/modelTest/deleteTestHistory', '删除测试历史', '模型测试', 'DELETE'),
(NOW(), NOW(), '/modeltraining/modelTest/clearTestHistory', '清空测试历史', '模型测试', 'DELETE'),
(NOW(), NOW(), '/modeltraining/modelTest/getTestHistoryList', '获取测试历史列表', '模型测试', 'GET');

-- ============================================
-- 第五步：插入Casbin权限规则
-- ============================================

-- 为超级管理员(888)添加权限
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
-- 数据集管理权限
('p', '888', '/modeltraining/dataset/createDataset', 'POST'),
('p', '888', '/modeltraining/dataset/deleteDataset', 'DELETE'),
('p', '888', '/modeltraining/dataset/deleteDatasetByIds', 'DELETE'),
('p', '888', '/modeltraining/dataset/updateDataset', 'PUT'),
('p', '888', '/modeltraining/dataset/findDataset', 'GET'),
('p', '888', '/modeltraining/dataset/getDatasetList', 'GET'),
('p', '888', '/modeltraining/dataset/getDatasetDataSource', 'GET'),
('p', '888', '/modeltraining/dataset/createVersion', 'POST'),
('p', '888', '/modeltraining/dataset/deleteVersion', 'DELETE'),
('p', '888', '/modeltraining/dataset/getVersionList', 'GET'),
('p', '888', '/modeltraining/dataset/publishDataset', 'POST'),
('p', '888', '/modeltraining/dataset/uploadFile', 'POST'),
('p', '888', '/modeltraining/dataset/uploadVersionFile', 'POST'),
-- 训练任务权限
('p', '888', '/modeltraining/trainingTask/createTask', 'POST'),
('p', '888', '/modeltraining/trainingTask/deleteTask', 'DELETE'),
('p', '888', '/modeltraining/trainingTask/deleteTaskByIds', 'DELETE'),
('p', '888', '/modeltraining/trainingTask/updateTask', 'PUT'),
('p', '888', '/modeltraining/trainingTask/findTask', 'GET'),
('p', '888', '/modeltraining/trainingTask/getTaskList', 'GET'),
('p', '888', '/modeltraining/trainingTask/startTask', 'POST'),
('p', '888', '/modeltraining/trainingTask/stopTask', 'POST'),
('p', '888', '/modeltraining/trainingTask/markCompleted', 'POST'),
('p', '888', '/modeltraining/trainingTask/startService', 'POST'),
('p', '888', '/modeltraining/trainingTask/stopService', 'POST'),
('p', '888', '/modeltraining/trainingTask/chatCompletion', 'POST'),
('p', '888', '/modeltraining/trainingTask/getTaskLogs', 'GET'),
('p', '888', '/modeltraining/trainingTask/getTaskDataSource', 'GET'),
('p', '888', '/modeltraining/trainingTask/getDefaultParams', 'GET'),
-- 模型测试历史权限
('p', '888', '/modeltraining/modelTest/createTestHistory', 'POST'),
('p', '888', '/modeltraining/modelTest/deleteTestHistory', 'DELETE'),
('p', '888', '/modeltraining/modelTest/clearTestHistory', 'DELETE'),
('p', '888', '/modeltraining/modelTest/getTestHistoryList', 'GET');

-- ============================================
-- 第六步：创建数据表（如果不存在）
-- ============================================

CREATE TABLE IF NOT EXISTS `dataset` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL COMMENT '数据集名称',
  `type` varchar(20) NOT NULL COMMENT '数据集类型',
  `format` varchar(50) DEFAULT NULL COMMENT '数据格式',
  `train_method` varchar(20) DEFAULT NULL COMMENT '训练方式',
  `storage_path` varchar(255) DEFAULT NULL COMMENT '存储路径',
  `description` varchar(500) DEFAULT NULL COMMENT '描述',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '创建者ID',
  `latest_version` varchar(20) DEFAULT NULL COMMENT '最新版本',
  `data_count` bigint DEFAULT 0 COMMENT '数据量',
  `import_status` varchar(20) DEFAULT 'pending' COMMENT '导入状态',
  `publish_status` tinyint(1) DEFAULT 0 COMMENT '发布状态',
  PRIMARY KEY (`id`),
  KEY `idx_dataset_deleted_at` (`deleted_at`),
  KEY `idx_dataset_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='数据集表';

CREATE TABLE IF NOT EXISTS `dataset_version` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `dataset_id` bigint unsigned NOT NULL COMMENT '数据集ID',
  `version` varchar(20) NOT NULL COMMENT '版本号',
  `data_count` bigint DEFAULT 0 COMMENT '数据量',
  `storage_path` varchar(255) DEFAULT NULL COMMENT '存储路径',
  `description` varchar(500) DEFAULT NULL COMMENT '版本说明',
  `file_size` bigint DEFAULT NULL COMMENT '文件大小(字节)',
  `status` varchar(20) DEFAULT 'pending' COMMENT '状态',
  `file_name` varchar(255) DEFAULT NULL COMMENT '文件名称',
  `file_path` varchar(500) DEFAULT NULL COMMENT '文件路径',
  PRIMARY KEY (`id`),
  KEY `idx_dataset_version_deleted_at` (`deleted_at`),
  KEY `idx_dataset_version_dataset_id` (`dataset_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='数据集版本表';

CREATE TABLE IF NOT EXISTS `training_task` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL COMMENT '任务名称',
  `task_id` varchar(50) DEFAULT NULL COMMENT '任务ID',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '创建者ID',
  `base_model` varchar(100) DEFAULT NULL COMMENT '基础模型',
  `train_method` varchar(20) DEFAULT NULL COMMENT '训练方式',
  `train_type` varchar(20) DEFAULT NULL COMMENT '训练类型',
  `status` varchar(20) DEFAULT 'pending' COMMENT '训练状态',
  `train_dataset_id` bigint unsigned DEFAULT NULL COMMENT '训练集ID',
  `train_version_id` bigint unsigned DEFAULT NULL COMMENT '训练集版本ID',
  `val_dataset_id` bigint unsigned DEFAULT NULL COMMENT '验证集ID',
  `val_version_id` bigint unsigned DEFAULT NULL COMMENT '验证集版本ID',
  `val_split_ratio` double DEFAULT NULL COMMENT '验证集切分比例',
  `output_count` int DEFAULT 5 COMMENT '产出数量上限',
  `model_name` varchar(100) DEFAULT NULL COMMENT '输出模型名称',
  `checkpoint_interval` int DEFAULT NULL COMMENT 'Checkpoint保存间隔',
  `checkpoint_unit` varchar(10) DEFAULT NULL COMMENT 'Checkpoint间隔单位',
  `progress` int DEFAULT 0 COMMENT '训练进度(百分比)',
  `start_time` datetime(3) DEFAULT NULL COMMENT '开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '结束时间',
  `node_id` bigint unsigned DEFAULT NULL COMMENT '执行节点ID',
  `instance_id` bigint unsigned DEFAULT NULL COMMENT '实例ID',
  `host_port` int DEFAULT NULL COMMENT '训练容器端口',
  `container_id` varchar(128) DEFAULT NULL COMMENT '训练容器ID',
  `container_name` varchar(255) DEFAULT NULL COMMENT '训练容器名称',
  `checkpoint_path` varchar(512) DEFAULT NULL COMMENT '训练产出Checkpoint路径',
  `remark` varchar(1000) DEFAULT NULL COMMENT '备注信息',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_training_task_task_id` (`task_id`),
  KEY `idx_training_task_deleted_at` (`deleted_at`),
  KEY `idx_training_task_user_id` (`user_id`),
  KEY `idx_training_task_status` (`status`),
  KEY `idx_training_task_node_id` (`node_id`),
  KEY `idx_training_task_instance_id` (`instance_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='训练任务表';

CREATE TABLE IF NOT EXISTS `training_param` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `task_id` bigint unsigned NOT NULL COMMENT '任务ID',
  `batch_size` int DEFAULT NULL COMMENT '批次大小',
  `learning_rate` double DEFAULT NULL COMMENT '学习率',
  `n_epochs` int DEFAULT NULL COMMENT '训练轮数',
  `eval_steps` int DEFAULT NULL COMMENT '验证步数',
  `lora_alpha` int DEFAULT NULL COMMENT 'LoRa缩放系数',
  `lora_dropout` double DEFAULT NULL COMMENT 'LoRa Dropout',
  `lora_rank` int DEFAULT NULL COMMENT 'LoRa秩值',
  `lr_scheduler_type` varchar(50) DEFAULT NULL COMMENT '学习率调整策略',
  `max_length` int DEFAULT NULL COMMENT '序列长度',
  `warmup_ratio` double DEFAULT NULL COMMENT '学习率预热比例',
  `weight_decay` double DEFAULT NULL COMMENT '权重衰减',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_training_param_task_id` (`task_id`),
  KEY `idx_training_param_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='训练参数表';

CREATE TABLE IF NOT EXISTS `model_test_history` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `task_id` bigint unsigned NOT NULL COMMENT '训练任务ID',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '创建者ID',
  `question` text NOT NULL COMMENT '测试问题',
  `base_answer` text COMMENT '基础模型回复',
  `lora_answer` text COMMENT 'LoRA模型回复',
  `test_time` datetime(3) DEFAULT NULL COMMENT '测试时间',
  PRIMARY KEY (`id`),
  KEY `idx_model_test_history_deleted_at` (`deleted_at`),
  KEY `idx_model_test_history_task_id` (`task_id`),
  KEY `idx_model_test_history_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='模型测试历史表';