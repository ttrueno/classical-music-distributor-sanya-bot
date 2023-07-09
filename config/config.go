package config

type Config struct {
	BotConfig    BotConfig    `yaml:"bot_config"`
	DbConnConfig DbConnConfig `yaml:"db_conn_config"`
}

type BotConfig struct {
	BotApiToken string `yaml:"tg_bot_api_token"`
}

type DbConnConfig struct {
	Dsn            string         `yaml:"dsn"`
	ConnPoolConfig ConnPoolConfig `yaml:"conn_pool_config"`
}

type ConnPoolConfig struct {
	MaxConns        int    `yaml:"max_conns"`
	MaxConnIdleTime string `yaml:"max_conn_idle_time"`
	MaxConnLifeTime string `yaml:"max_conn_life_time"`
}
