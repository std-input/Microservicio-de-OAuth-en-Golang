package database

import (
	"adbr.xx/auth_microservice/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabase() {

	var err error
	// sqlite
	DB, err = gorm.Open(sqlite.Open("auth.db"), &gorm.Config{})

	// PGSQL
	/*pgsql_conf := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		configs.Get("DB_HOST"), configs.Get("DB_USER"), configs.Get("DB_PASSWORD"),
		configs.Get("DB_NAME"), configs.Get("DB_PORT"), configs.Get("DB_SSLMODE"),
	)*/
	//DB, err = gorm.Open(postgres.Open(pgsql_conf), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(models.User{}, models.RefreshToken{})
}

func GetUserById(id string) (int64, models.SafeUser) {
	user := models.User{}
	var userSafe models.SafeUser
	r := DB.Where("id = ?", id).Omit("email").First(&user).RowsAffected
	userSafe = models.SafeUser{
		ID:          user.ID,
		Name:        user.Name,
		FamilyName:  user.FamilyName,
		Picture:     user.Picture,
		Description: user.Description,
		Role:        user.Role,
	}
	return r, userSafe
}

func DeleteUser(id string) int64 {
	return DB.Delete(models.User{ID: id}).RowsAffected
}

func UpdateUser(user *models.User) int64 {
	return DB.Select("description").Updates(user).RowsAffected
}

func GetOrCreateUser(user *models.User) int64 {
	result := DB.Where("id = ?", user.ID).FirstOrInit(&user)
	if result.RowsAffected == 0 {
		user.Role = "user"
		DB.Create(&user)
	}
	return result.RowsAffected
}
