package sqlmigrate

func (migrate *SqlMigration) NewMigration(name string) error {
	return migrate.db.NewMigration(name)
}

// apply all migrations
func (migrate *SqlMigration) Up() error {
	return migrate.db.CreateAndMigrate()
}

// rollback 1 latest migration
func (migrate *SqlMigration) Down() error {
	return migrate.db.Rollback()
}
