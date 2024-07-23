package ck

import (
	"database/sql"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/gowins/dionysus/log"
	"github.com/spf13/viper"
)

var ckClient *sql.DB

func Setup() {
	client, err := sql.Open("clickhouse", viper.GetString("Clickhouse.Log.DataSource"))
	if err != nil {
		log.Fatalf("open clickhouse failed %s", err.Error())
	}
	ckClient = client
}

func GetCk() *sql.DB {
	return ckClient
}
