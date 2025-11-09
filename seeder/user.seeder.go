package seeder

import (
	"fmt"
	"go-sse/user"
	"go-sse/util"
	"reflect"
	"time"

	"github.com/bxcodec/faker"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserSeeder(db *gorm.DB) {
	users := []user.User{
		newUser("Luffy"),
		newUser("Zoro"),
		newUser("Sanji"),
	}

	for range 20 - len(users) {
		randomUuid, _ := uuid.NewRandom()
		hashedPassword, _ := util.HashPassword("Qwerty123")
		name, _ := faker.GetPerson().Name(reflect.Value{})
		nameSlug := util.Slugify(name.(string))
		randomDate := util.RandomDate(time.Now().AddDate(0, 0, -1), time.Now())
		user := user.User{
			Id:        randomUuid,
			Name:      name.(string),
			Username:  &nameSlug,
			Password:  &hashedPassword,
			CreatedAt: randomDate,
			UpdatedAt: randomDate,
		}
		users = append(users, user)
	}

	result := db.Model(&user.User{}).Create(&users)
	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Printf("%d users inserted\n", result.RowsAffected)
}

func newUser(name string) user.User {
	newUuid, _ := uuid.NewRandom()
	hashedPassword, _ := util.HashPassword("Qwerty123")
	randomDate := util.RandomDate(time.Now().AddDate(0, 0, -1), time.Now())
	username := util.Slugify(name)
	return user.User{
		Id:        newUuid,
		Name:      name,
		Username:  &username,
		Password:  &hashedPassword,
		CreatedAt: randomDate,
		UpdatedAt: randomDate,
	}

}
