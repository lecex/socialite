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
	seeds()
}

// config 配置数据迁移
func config() {
	config := &cpd.Config{}
	if !db.DB.HasTable(&config) {
		db.DB.Exec(`
			CREATE TABLE configs (
			name varchar(32) NOT NULL COMMENT '配置名称',
			value json DEFAULT NULL COMMENT '配置内容',
			PRIMARY KEY (name)
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
			user_id varchar(36) NOT NULL,
			created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY origin_openid (origin,openid)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;
		`)
	}
}

// user 用户数据迁移
func user() {
	user := &userPD.User{}
	if !db.DB.HasTable(&user) {
		db.DB.Exec(`
			CREATE TABLE socialite_users (
			id varchar(36) NOT NULL,
			name varchar(64) DEFAULT NULL,
			socialite_id varchar(36) NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY origin_openid (origin,openid)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;
		`)
	}
}

// seeds 填充文件
func seeds() {
	db.DB.Exec(`
		INSERT INTO configs ( name ) VALUES ('config')
	`)
}
