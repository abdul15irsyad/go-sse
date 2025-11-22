package user

import (
	"go-sse/util"
	"time"

	"github.com/google/uuid"
)

func CreateUser(dto CreateUserDTO) (User, error) {
	newUuid, _ := uuid.NewRandom()
	var password *string = nil
	if dto.Password != nil {
		hashPassword, err := util.HashPassword(*dto.Password)
		if err != nil {
			return User{}, err
		}
		password = &hashPassword
	}

	newUser := User{
		Id:        newUuid,
		Name:      dto.Name,
		Username:  dto.Username,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := util.DB.Model(&User{}).Create(newUser).Error; err != nil {
		return User{}, err
	}

	return newUser, nil
}

func GetPaginatedFriends(authUserId uuid.UUID, page int, limit int) ([]User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	users := []User{}
	query := util.DB.Model(&User{}).Where("id <> ?", authUserId)
	if err := query.Limit(limit).Offset((page - 1) * limit).Order("name ASC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func GetUserByUsername(username string) (User, error) {
	var user User
	if err := util.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUser(id uuid.UUID) (User, error) {
	var user User
	if err := util.DB.Model(&User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}
