func ConnectDatabase() {
	dsn := os.Getenv("MYSQL_DSN")
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.User{},      
		&models.Company{},     
		&models.Job{},       
		&models.Profile{},  
		&models.Application{}, 
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database connection established and migrations completed")
}
