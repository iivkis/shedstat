include .env

click-db:
	clickhouse client --host localhost --port 9000 --user $(CLICKHOUSE_USER) --password $(CLICKHOUSE_PASSWORD) --database $(CLICKHOUSE_NAME)