package migration

import (
	"belajar/database"
	"belajar/model/entity"
	"fmt"
	"log"
)

func RunMigration() {
	err := database.DB.AutoMigrate(&entity.User{}, &entity.Book{}, &entity.Category{}, &entity.Photo{})
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Database Migrate")
}
