CREATE TABLE `chat_message` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'メッセージの識別子',
  `category`    VARCHAR(128) NOT NULL COMMENT 'メッセージのカテゴリ',
  `message`     Text NOT NULL COMMENT 'メッセージ',
  `role`        VARCHAR(128) NOT NULL COMMENT 'ロール',
  `created`     DATETIME(6) NOT NULL COMMENT 'レコード作成日時',
  `modified`    DATETIME(6) NOT NULL COMMENT 'レコード修正日時',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='チャットメッセージ';

CREATE TABLE `chat_summary` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '要約の識別子',
  `category`    VARCHAR(128) NOT NULL COMMENT '要約のカテゴリ',
  `summary`     Text NOT NULL COMMENT '要約',
  `created`     DATETIME(6) NOT NULL COMMENT 'レコード作成日時',
  `modified`    DATETIME(6) NOT NULL COMMENT 'レコード修正日時',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='要約';
