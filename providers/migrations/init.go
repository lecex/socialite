package migrations

import (
	cpd "github.com/lecex/socialite/proto/config"
	userPD "github.com/lecex/socialite/proto/user"
	db "github.com/lecex/socialite/providers/database"
)

func init() {
	config()
	socialiteUser()
	user()
}

// config 配置数据迁移
func config() {
	config := &cpd.Config{}
	if !db.DB.HasTable(&config) {
		db.DB.Exec(`
			CREATE TABLE configs (
				id int(11) unsigned NOT NULL AUTO_INCREMENT,
				name varchar(32) NOT NULL COMMENT '驱动名称',
				driver varchar(32) NOT NULL COMMENT '驱动标识',
				client_id varchar(80) NOT NULL COMMENT '客户ID',
				client_secret varchar(80) NOT NULL COMMENT '客户密钥',	
				redirect varchar(180) NOT NULL COMMENT '回调地址',	
				status tinyint(1) DEFAULT 0 COMMENT '状态',
				created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (id),
				UNIQUE KEY driver_and_client_id (driver,client_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;
		`)
	}
}

// socialiteUser 用户数据迁移
func socialiteUser() {
	socialiteUser := &userPD.SocialiteUser{}
	if !db.DB.HasTable(&socialiteUser) {
		db.DB.Exec(`
			CREATE TABLE socialite_users (
			id varchar(36) NOT NULL,
			origin varchar(64) DEFAULT NULL,
			oauth_id varchar(64) DEFAULT NULL,
			content json DEFAULT NULL COMMENT '请求返回内容',
			created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY origin_oauth_id (origin,oauth_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;
		`)
	}
}

// user 用户数据迁移
func user() {
	user := &userPD.User{}
	if !db.DB.HasTable(&user) {
		db.DB.Exec(`
			CREATE TABLE users (
			id varchar(36) NOT NULL,
			socialite_user_id varchar(36) NOT NULL,
			PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;
		`)
	}
}
