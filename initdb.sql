CREATE TABLE `comments` (
  `id` varchar(64) NOT NULL,
  `video_id` varchar(64) DEFAULT NULL,
  `author_id` int(11) DEFAULT NULL,
  `content` text,
  `time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `sessions` (
  `session_id` varchar(64) NOT NULL,
  `TTL` tinytext,
  `login_name` text,
  PRIMARY KEY (`session_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#
# alter table sessions add primary key (session_id(64));
#

CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `login_name` varchar(64) DEFAULT NULL,
  `pwd` text NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `login_name` (`login_name`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;


CREATE TABLE `video_del_rec` (
  `video_id` varchar(64) NOT NULL,
  PRIMARY KEY (`video_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;


CREATE TABLE `video_info` (
  `id` char(64) NOT NULL,
  `author_id` int(10) unsigned DEFAULT NULL,
  `name` text,
  `display_ctime` text,
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;