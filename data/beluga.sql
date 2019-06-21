/*
Navicat MySQL Data Transfer

Source Server         : 自己服务器
Source Server Version : 50726
Source Host           : 121.42.44.23:3306
Source Database       : beluga

Target Server Type    : MYSQL
Target Server Version : 50726
File Encoding         : 65001

Date: 2019-06-13 10:06:45
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
  PRIMARY KEY (`id`),
  KEY `token` (`token`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

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
) ENGINE=InnoDB AUTO_INCREMENT=248 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='配置';

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
) ENGINE=InnoDB AUTO_INCREMENT=351 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='配置详细记录';

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
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

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
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

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
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='配置中心节点配置';

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
) ENGINE=InnoDB AUTO_INCREMENT=79 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='配置操作';

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
) ENGINE=InnoDB AUTO_INCREMENT=530 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='配置发布|回滚记录';

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
) ENGINE=InnoDB AUTO_INCREMENT=211258 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

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
) ENGINE=InnoDB AUTO_INCREMENT=102 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='配置版本';

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
) ENGINE=InnoDB AUTO_INCREMENT=677 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='操作日志';
