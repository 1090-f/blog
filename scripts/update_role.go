package main

import (
	"blog/config"
	"blog/pkg/database"
	"blog/internal/model"
	"log"
)

func main() {
	cfg := config.MustLoad()
	db := database.MustOpen(cfg.Database)

	// Update admin user role
	result := db.Model(&model.User{}).Where("username = ?", "admin").Update("role", "admin")
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	log.Printf("Updated %d user(s) to admin role", result.RowsAffected)
}
