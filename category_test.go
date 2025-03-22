package categories

import (
	"fmt"
	"log"
	"testing"

	"github.com/spurtcms/auth"
)

// test listcategory function
func TestListCategory(t *testing.T) {

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

	permisison, _ := Auth.IsGranted("Categories", auth.CRUD, "1")

	Category := CategoriesSetup(Config{
		DB:               db,
		AuthEnable:       true,
		PermissionEnable: true,
		Auth:             Auth,
	})
	if permisison {

		Categorylist, fnlist, parentcategory, count, err := Category.ListCategory(10, 0, Filter{}, 1, "1")

		if err != nil {

			panic(err)
		}

		fmt.Println(Categorylist, fnlist, parentcategory, count)
	} else {

		log.Println("permissions enabled not initialised")

	}

}

// test addcategory function
func TestAddCategory(t *testing.T) {

	db, _ := DBSetup()

	Category := CategoriesSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})
	err := Category.AddCategory(CategoryCreate{CategoryName: "Default", CategorySlug: "default", ParentId: 1, TenantId: "1"})

	if err != nil {

		panic(err)
	}

}

// test updatesubcategory function
func TestUpdateSubCategory(t *testing.T) {

	db, _ := DBSetup()

	Category := CategoriesSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})
	err := Category.UpdateSubCategory(CategoryCreate{Id: 10, CategoryName: "Default", CategorySlug: "default", ParentId: 2}, "1")

	if err != nil {

		panic(err)
	}

}

// test deletesubcategory function
func TestDeleteSubCategory(t *testing.T) {

	db, _ := DBSetup()

	Category := CategoriesSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})
	err := Category.DeleteSubCategory(2, 1, "1")

	if err != nil {

		panic(err)
	}

}

// test getsubcategorydetails function
func TestGetSubCategoryDetails(t *testing.T) {

	db, _ := DBSetup()

	Category := CategoriesSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})
	category, err := Category.GetSubCategoryDetails(2, "1")

	if err != nil {

		panic(err)
	}

	log.Println(category)

}
