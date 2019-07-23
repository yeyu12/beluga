/*
Navicat MySQL Data Transfer

Source Server         : 自己服务器
Source Server Version : 50726
Source Host           : 121.42.44.23:3306
Source Database       : beluga

Target Server Type    : MYSQL
Target Server Version : 50726
File Encoding         : 65001

Date: 2019-07-23 09:37:31
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for beluga_account
-- ----------------------------
DROP TABLE IF EXISTS `beluga_account`;
CREATE TABLE `beluga_account` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `passwd` char(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `nickname` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `avatar` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `gender` tinyint(1) DEFAULT NULL,
  `age` tinyint(2) DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  `create_time` datetime DEFAULT NULL,
  `token_expiry_time` datetime DEFAULT NULL,
  `token` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `configuration_num` int(10) unsigned DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `token` (`token`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for beluga_configuration
-- ----------------------------
DROP TABLE IF EXISTS `beluga_configuration`;
CREATE TABLE `beluga_configuration` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `appid` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `namespace_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `key` text COLLATE utf8_unicode_ci NOT NULL,
  `val` text COLLATE utf8_unicode_ci NOT NULL,
  `remake` text COLLATE utf8_unicode_ci,
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `appid` (`appid`),
  KEY `namespace_name` (`namespace_name`),
  FULLTEXT KEY `key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=279 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='配置';

-- ----------------------------
-- Table structure for beluga_configuration_log
-- ----------------------------
DROP TABLE IF EXISTS `beluga_configuration_log`;
CREATE TABLE `beluga_configuration_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `appid` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `namespace_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `version` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `account_id` int(11) DEFAULT NULL,
  `key` text COLLATE utf8_unicode_ci,
  `val` text COLLATE utf8_unicode_ci,
  `remake` text COLLATE utf8_unicode_ci,
  `type` tinyint(1) DEFAULT '0',
  `is_release` tinyint(1) DEFAULT '0',
  `old_val` text COLLATE utf8_unicode_ci,
  `old_remake` text COLLATE utf8_unicode_ci,
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `appid` (`appid`),
  KEY `namespace_name` (`namespace_name`),
  KEY `version` (`version`),
  KEY `account_id` (`account_id`),
  KEY `is_release` (`is_release`),
  KEY `type` (`type`),
  FULLTEXT KEY `key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=498 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='配置详细记录';

-- ----------------------------
-- Table structure for beluga_configuration_namespace
-- ----------------------------
DROP TABLE IF EXISTS `beluga_configuration_namespace`;
CREATE TABLE `beluga_configuration_namespace` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(11) DEFAULT NULL,
  `namespace_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `account_id` int(11) DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `project_id` (`project_id`),
  KEY `namespace_name` (`namespace_name`),
  KEY `account_id` (`account_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for beluga_configuration_node
-- ----------------------------
DROP TABLE IF EXISTS `beluga_configuration_node`;
CREATE TABLE `beluga_configuration_node` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ip` varchar(15) COLLATE utf8_unicode_ci DEFAULT NULL,
  `node_conf_id` text COLLATE utf8_unicode_ci,
  `conf_update_time` datetime DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  `conf_update_account_id` int(10) unsigned DEFAULT NULL,
  `remake` text COLLATE utf8_unicode_ci,
  `is_delete` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `ip` (`ip`),
  FULLTEXT KEY `node_conf_id` (`node_conf_id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for beluga_configuration_node_conf
-- ----------------------------
DROP TABLE IF EXISTS `beluga_configuration_node_conf`;
CREATE TABLE `beluga_configuration_node_conf` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `conf` text COLLATE utf8_unicode_ci,
  `account_id` int(10) unsigned DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='配置中心节点配置';

-- ----------------------------
-- Table structure for beluga_configuration_operation
-- ----------------------------
DROP TABLE IF EXISTS `beluga_configuration_operation`;
CREATE TABLE `beluga_configuration_operation` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `appid` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `namespace_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `version` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `account_id` int(11) DEFAULT NULL,
  `operation_type` tinyint(1) DEFAULT '0',
  `rollback` tinyint(1) DEFAULT '0',
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `appid` (`appid`),
  KEY `namespace_name` (`namespace_name`),
  KEY `version` (`version`),
  KEY `account_id` (`account_id`),
  KEY `operation_type` (`operation_type`)
) ENGINE=InnoDB AUTO_INCREMENT=88 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='配置操作';

-- ----------------------------
-- Table structure for beluga_configuration_operation_log
-- ----------------------------
DROP TABLE IF EXISTS `beluga_configuration_operation_log`;
CREATE TABLE `beluga_configuration_operation_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `configuration_operation_id` int(11) NOT NULL,
  `key` text COLLATE utf8_unicode_ci,
  `val` text COLLATE utf8_unicode_ci,
  `account_id` int(10) unsigned DEFAULT NULL,
  `remake` text COLLATE utf8_unicode_ci,
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `configuration_operation_id` (`configuration_operation_id`),
  FULLTEXT KEY `key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=565 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='配置发布|回滚记录';

-- ----------------------------
-- Table structure for beluga_configuration_project
-- ----------------------------
DROP TABLE IF EXISTS `beluga_configuration_project`;
CREATE TABLE `beluga_configuration_project` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '项目名',
  `account_id` int(10) unsigned DEFAULT NULL COMMENT '创建人id',
  `create_time` datetime DEFAULT NULL,
  `appid` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '项目appid',
  PRIMARY KEY (`id`),
  KEY `appid` (`appid`),
  KEY `project_name` (`project_name`)
) ENGINE=InnoDB AUTO_INCREMENT=211263 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Table structure for beluga_configuration_version
-- ----------------------------
DROP TABLE IF EXISTS `beluga_configuration_version`;
CREATE TABLE `beluga_configuration_version` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `p_id` int(11) DEFAULT NULL,
  `version` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `namespace_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `appid` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_release` tinyint(1) DEFAULT '0',
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `p_id` (`p_id`),
  KEY `appid` (`appid`),
  KEY `namespace_name` (`namespace_name`),
  KEY `version` (`version`)
) ENGINE=InnoDB AUTO_INCREMENT=111 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='配置版本';

-- ----------------------------
-- Table structure for beluga_operation_log
-- ----------------------------
DROP TABLE IF EXISTS `beluga_operation_log`;
CREATE TABLE `beluga_operation_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `c` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `params` text COLLATE utf8_unicode_ci,
  `account_id` int(11) DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  `ident` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=2381 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='操作日志';

-- ----------------------------
-- Table structure for beluga_task
-- ----------------------------
DROP TABLE IF EXISTS `beluga_task`;
CREATE TABLE `beluga_task` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `create_time` datetime DEFAULT NULL,
  `start_time` datetime DEFAULT NULL COMMENT '开始执行时间',
  `consume_time` double(5,4) unsigned DEFAULT NULL COMMENT '执行消耗时间',
  `overtime` int(10) unsigned DEFAULT '0' COMMENT '超时时间，s为单位,0不限制超时时间',
  `last_exec_type` tinyint(1) DEFAULT '0' COMMENT '上次执行状态（0无状态，1执行成功，-1执行失败',
  `task_type` tinyint(1) unsigned DEFAULT '1' COMMENT '任务类型（1主任务，2子任务',
  `rely` tinyint(1) DEFAULT '0' COMMENT '依赖关系（1强依赖，0弱依赖',
  `subtasks_id` text COLLATE utf8_unicode_ci COMMENT '子任务id，以英文都好隔开',
  `cron` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT 'crontab表达式',
  `task_exec_type` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '任务执行类型（shell,http',
  `exec_task_node_type` tinyint(1) DEFAULT '1' COMMENT '执行任务节点类型（1随机，0指定。shell必须要指定',
  `exec_task_node_id` text COLLATE utf8_unicode_ci COMMENT '执行任务节点id',
  `cmd` text COLLATE utf8_unicode_ci COMMENT '执行的shell命令或http请求地址',
  `http_type` varchar(4) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT 'http请求类型，post，get',
  `task_fail_num` int(10) unsigned DEFAULT '0' COMMENT '任务失败次数，0不限制',
  `task_fail_retry_time` int(10) unsigned DEFAULT '0' COMMENT '任务失败重试时间间隔s为单位，0失败则执行',
  `task_notice` tinyint(1) unsigned DEFAULT '0' COMMENT '任务通知，0不通知，1失败通知，2总是通知，3关键字通知',
  `notice_type` tinyint(1) unsigned DEFAULT '0' COMMENT '通知类型(0无通知，1邮件通知，2webhook通知',
  `keyword_notice` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '关键字通知',
  `remake` text COLLATE utf8_unicode_ci COMMENT '备注',
  `account_id` int(11) DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1' COMMENT '任务状态（-1删除，1启动，2停止',
  `next_exec_time` datetime DEFAULT NULL COMMENT '下次执行时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='任务';

-- ----------------------------
-- Table structure for beluga_task_log
-- ----------------------------
DROP TABLE IF EXISTS `beluga_task_log`;
CREATE TABLE `beluga_task_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` int(10) unsigned DEFAULT NULL COMMENT '任务id',
  `task_exec_type` tinyint(1) unsigned DEFAULT '1' COMMENT '任务执行状态（1成功，0失败',
  `create_time` datetime DEFAULT NULL,
  `end_time` datetime DEFAULT NULL COMMENT '执行结束时间',
  `node_ip` char(15) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '执行的节点ip',
  `consume_time` double(5,4) unsigned DEFAULT '0.0000' COMMENT '执行消耗时间,秒为单位',
  `err` text COLLATE utf8_unicode_ci,
  `task_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `cmd` text COLLATE utf8_unicode_ci,
  `output` text COLLATE utf8_unicode_ci COMMENT '输出结果',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=89224 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='任务执行日志';

-- ----------------------------
-- Table structure for beluga_task_node
-- ----------------------------
DROP TABLE IF EXISTS `beluga_task_node`;
CREATE TABLE `beluga_task_node` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ip` varchar(15) COLLATE utf8_unicode_ci DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  `remake` text COLLATE utf8_unicode_ci,
  `is_delete` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `ip` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

INSERT INTO `beluga`.`beluga_account` (`id`, `username`, `passwd`, `nickname`, `avatar`, `gender`, `age`, `status`, `create_time`, `token_expiry_time`, `token`, `configuration_num`) VALUES ('1', 'admin', '31bf43eee77df03c65f4ec90629b3e9c', 'admin', '', '0', '0', '1', '2019-04-16 18:07:20', '2019-07-22 23:36:23', 'a77d764e7fb0700395ae5bdfc1d6e3e2', '2');