package main

import (
	"errors"
	"fmt"
	"go-sse/seeder"
	"go-sse/util"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	for _, dataSeeder := range []struct {
		name string
		seed func(*gorm.DB)
	}{
		{name: "userSeeder", seed: seeder.UserSeeder},
	} {
		// check name in seeders table
		result := util.DB.Model(&seeder.Seeder{}).Where("name = ?", dataSeeder.name).First(&dataSeeder)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			panic(result.Error)
		}
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Printf("%s already executed\n", dataSeeder.name)
			continue
		}

		fmt.Printf("executing %s...\n", dataSeeder.name)
		dataSeeder.seed(util.DB)

		// add name to seeders table
		Id, _ := uuid.NewRandom()
		util.DB.Create(&seeder.Seeder{
			Id:        Id,
			Name:      dataSeeder.name,
			CreatedAt: time.Now(),
		})
		fmt.Printf("%s executed\n", dataSeeder.name)
	}

	fmt.Println("seeder done")
}
