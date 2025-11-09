package main

import (
	"fmt"
	"go-sse/notification"
	"go-sse/seeder"
	"go-sse/user"
	"go-sse/util"
)

func main() {
	// start migrate
	if !util.DB.Migrator().HasTable(&user.User{}) {
		util.DB.Migrator().CreateTable(&user.User{})
	}
	if !util.DB.Migrator().HasTable(&seeder.Seeder{}) {
		util.DB.Migrator().CreateTable(&seeder.Seeder{})
	}
	if !util.DB.Migrator().HasTable(&notification.Notification{}) {
		util.DB.Migrator().CreateTable(&notification.Notification{})
	}

	fmt.Println("migration done")
}
