package categories

import (
	"fmt"
	"log"
	"testing"

	"github.com/spurtcms/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var SecretKey = "Secret123"

// Db connection
func DBSetup() (*gorm.DB, error) {

	dbConfig := map[string]string{
		"username": "postgres",
		"password": "postgres",
		"host":     "localhost",
		"port":     "5432",
		"dbname":   "Spurtcms_V2",
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "user=" + dbConfig["username"] + " password=" + dbConfig["password"] +
			" dbname=" + dbConfig["dbname"] + " host=" + dbConfig["host"] +
			" port=" + dbConfig["port"] + " sslmode=disable TimeZone=Asia/Kolkata",
	}), &gorm.Config{})

	if err != nil {

		log.Fatal("Failed to connect to database:", err)

	}
	if err != nil {

		return nil, err

	}

	return db, nil
}

// test categorygrouplist function
func TestCategoryGroupList(t *testing.T) {

	db, _ := DBSetup()

	config := auth.Config{
		UserId: 1,
		// ExpiryTime: 2,
		ExpiryFlg: false,
		SecretKey: "Secret123",
		DB:        db,
		RoleId:    2,
	}

	Auth := auth.AuthSetup(config)

	token, _ := Auth.CreateToken()

	Auth.VerifyToken(token, SecretKey)

	permisison, _ := Auth.IsGranted("Categories Group", auth.CRUD, 1)

	Category := CategoriesSetup(Config{
		DB:               db,
		AuthEnable:       true,
		PermissionEnable: true,
		Auth:             Auth,
	})
	if permisison {

		Categorygrouplist, count, err := Category.CategoryGroupList(10, 0, Filter{}, 1)

		if err != nil {

			panic(err)
		}

		fmt.Println(Categorygrouplist, count)
	} else {

		log.Println("permissions enabled not initialised")

	}

}

// test createcategorygroup function
func TestCreateCategoryGroup(t *testing.T) {

	db, _ := DBSetup()

	Category := CategoriesSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})
	err := Category.CreateCategoryGroup(CategoryCreate{CategoryName: "indoor-sports", CategorySlug: "indoor_sports", Description: "type of indoor sports", TenantId: 1})

	if err != nil {

		panic(err)
	}

}

// test updatecategorygroupt function
func TestUpdateCategoryGroup(t *testing.T) {

	db, _ := DBSetup()

	Category := CategoriesSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})
	err := Category.UpdateCategoryGroup(CategoryCreate{Id: 1, CategoryName: "Default Category", CategorySlug: "default_category", Description: "Default Category"}, 1)

	if err != nil {

		panic(err)
	}

}

// test deletecategorygroup function
func TestDeleteCategoryGroup(t *testing.T) {

	db, _ := DBSetup()

	Category := CategoriesSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})
	err := Category.DeleteCategoryGroup(1, 1, 1)

	if err != nil {

		panic(err)
	}

}
