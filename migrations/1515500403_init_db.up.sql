CREATE TABLE `subs` ( `sfrom` INTEGER, `sto` INTEGER, PRIMARY KEY(`sfrom`,`sto`) );

CREATE TABLE `tweets` ( `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, `uid` INTEGER NOT NULL, `created_at` timestamp NOT NULL, `text` TEXT NOT NULL );

CREATE TABLE "users" ( `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, `name` TEXT NOT NULL, `nickname` TEXT NOT NULL UNIQUE, `password` BLOB NOT NULL );

CREATE INDEX `idx_tweets_uid` ON `tweets` ( `uid` );
