DROP TABLE IF EXISTS `authtoken`;
CREATE TABLE `authtoken` (
  `token` varchar(36) NOT NULL,
  `owner` int(11) unsigned NOT NULL,
  `expiry` int(11) unsigned NOT NULL,
  UNIQUE KEY `idx_user` (`token`)
) ENGINE=InnoDB CHARACTER SET utf8;

DROP TABLE IF EXISTS `invitetoken`;
CREATE TABLE `invitetoken` (
  `token` varchar(36) NOT NULL,
  `friend` int(11) unsigned NOT NULL,
  `account` int(11) unsigned,
  UNIQUE KEY `idx_token` (`token`),
  UNIQUE KEY `idx_user` (`account`)
) ENGINE=InnoDB CHARACTER SET utf8;
