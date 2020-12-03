include .env

migrate/init:
	mysql -u root -h 127.0.0.1 --protocol tcp -e "create database \`$(DB_NAME)\`" -p$(MYSQL_ROOT_PASSWORD)

migrate/new:
	@goose -dir migrations mysql "$(DATASOURCE)" create "$(NAME)" sql

migrate/up:
	@goose -dir migrations mysql "$(DATASOURCE)" up
