package config

// func DatabaseConnection() *gorm.DB {
// 	err := godotenv.Load(".env")

// 	if err != nil {
// 		log.Fatalf("Error: .Env file is not found.")
// 	}

// 	DBDriver := os.Getenv("DB_DRIVER")
// 	DBHost := os.Getenv("DB_HOST")
// 	DBUser := os.Getenv("DB_USER")
// 	DBPassword := os.Getenv("DB_PASSWORD")
// 	DBName := os.Getenv("DB_NAME")
// 	DBPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

// 	DBURL := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%d", DBHost, DBUser, DBName, DBPassword, DBPort)

// 	database, err := gorm.Open(postgres.Open(DBURL), &gorm.Config{})

// 	if err != nil {
// 		fmt.Println("Cannot connect to database ", DBDriver)
// 		log.Fatal("Error: Connection Error ", err)
// 	} else {
// 		fmt.Println("Message: We are connected to the database...")
// 	}

// 	return database
// }
