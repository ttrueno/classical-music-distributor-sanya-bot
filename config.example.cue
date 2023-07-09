#BotConfig: {
   tg_bot_api_token: string
}

#DbConnConfig: {
   dsn: string
   conn_pool_config: #ConnPoolConfig
}

#ConnPoolConfig: {
   max_conns: int
   max_conn_idle_time: string
   max_conn_life_time: string
}

bot_config: #BotConfig & {
         tg_bot_api_token: "<BOT_API_TOKEN>"
        }

db_conn_config: #DbConnConfig & {
                     dsn: "<DB_DSN>"
                     conn_pool_config: #ConnPoolConfig & {
                        max_conns: 50
                        max_conn_idle_time: "5m"
                        max_conn_life_time: "15m"
                     }
                }
