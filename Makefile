include .env

mgup: # Migrate the DB to the most recent version available
	@echo goose migration up
	goose -dir ${MIGRATION_DIR} mysql "${DB_USER}:${DB_PASS}@/${DB_NAME}?parseTime=true" up
mgdown: # Roll back the version by 1
	@echo goose migration down
	goose -dir ${MIGRATION_DIR} mysql "${DB_USER}:${DB_PASS}@/${DB_NAME}?parseTime=true" down
mgreset: # Roll back all migrations
	@echo goose migration reset
	goose -dir ${MIGRATION_DIR} mysql "${DB_USER}:${DB_PASS}@/${DB_NAME}?parseTime=true" reset
mgredo: # Roll back all migrations
	@echo goose migration redo
	goose -dir ${MIGRATION_DIR} mysql "${DB_USER}:${DB_PASS}@/${DB_NAME}?parseTime=true" redo
mgstatus: # Dump the migration status for the current DB
	@echo goose migration redo
	goose -dir ${MIGRATION_DIR} mysql "${DB_USER}:${DB_PASS}@/${DB_NAME}?parseTime=true" status
mgmk: # make migration
	@echo goose migration up
	goose -dir ${MIGRATION_DIR} create ${m} sql