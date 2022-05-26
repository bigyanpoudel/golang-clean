package infrastructure

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migration struct {
	logger Logger
	env    Env
	db     Database
}

func NewMigration(env Env, logger Logger, db Database) Migration {
	return Migration{
		env:    env,
		logger: logger,
		db:     db,
	}
}

//Migrate -> migrates all table
func (m Migration) Migrate() {
	m.logger.Zap.Info("Migrating schemas...")

	USER := m.env.DbUsername
	PASS := m.env.DbPassword
	HOST := m.env.DbHost
	PORT := m.env.DbPort
	DBNAME := m.env.DbName

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)

	if m.env.DbType == "cloudsql" {
		dsn = fmt.Sprintf(
			"%s:%s@unix(/cloudsql/%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			USER,
			PASS,
			HOST,
			DBNAME,
		)
	}

	// m.db.AutoMigrate(&models.User{}, &models.Post{})
	migrations, err := migrate.New("file://migration/", "mysql://"+dsn)

	if err != nil {
		panic(err)
	}

	m.logger.Zap.Info("--- Running Migration ---")
	_ = migrations.Steps(1000)
	fmt.Println("Migration Error:::", err)
	// err = migrations.Steps(1000)
	// if err != nil {
	// 	m.logger.Zap.Error("Error in migration: ", err.Error())
	// }
}
