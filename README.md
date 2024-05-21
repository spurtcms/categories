# categories Package

The 'Categories' package empowers administrators to define and structure content categories seamlessly. This package facilitates efficient content organization and user-friendly navigation within Golang projects, ensuring a streamlined browsing experience for users.



## Features

- Administrators can seamlessly navigate through the system's organizational structure by utilizing functions like CategoryGroupList, which retrieves an exhaustive list of existing category groups, while CreateCategoryGroup function allows the addition of new ones. 
- Updating and deleting category groups are facilitated by UpdateCategoryGroup and DeleteCategoryGroup functions respectively, ensuring adaptability and clutter-free organization. 
- ListCategory function furnishes a detailed rundown of available categories, while AddCategory empowers administrators to append new ones within existing groups. 
- Subcategory refinement is enabled through UpdateSubCategory and DeleteSubCategory functions, offering precision in content classification. 
- Detailed subcategory information is retrievable via GetSubCategoryDetails, and FilterSubCategory enables targeted management actions.
- CheckCategoryGroupName and CheckSubCategoryName functions guarantee the uniqueness of category and subcategory names respectively, upholding data integrity. 
- There is a master function AllCategoriesWithSubList that provides a holistic view of the category ecosystem, ensuring administrators possess the necessary tools for streamlined and effective content organization within the CMS.

# Installation

``` bash
go get github.com/spurtcms/categories
```


# Usage Example

``` bash
import (
	"github.com/spurtcms/auth"
	"github.com/spurtcms/categories"
)

func main() {

	Auth := auth.AuthSetup(auth.Config{
		UserId:     1,
		ExpiryTime: 2,
		SecretKey:  "SecretKey@123",
		DB: &gorm.DB{},
		RoleId: 1,
	})

	token, _ := Auth.CreateToken()

	Auth.VerifyToken(token, SecretKey)

	permisison, _ := Auth.IsGranted("Categories Group", auth.CRUD)

	Category := categories.CategoriesSetup(categories.Config{
		DB:               &gorm.DB{},
		AuthEnable:       true,
		PermissionEnable: true,
		Auth:             Auth,
	})

	//categorygroup
	if permisison {

		//list categorygroup
		Categorygrouplist, count, err := Category.CategoryGroupList(10, 0, categories.Filter{})
		fmt.Println(Categorygrouplist, count, err)

		//create categorygroup
		cerr := Category.CreateCategoryGroup(categories.CategoryCreate{
			CategoryName: "Default Group",
			CategorySlug: "default_group",
		})

		if cerr != nil {

			fmt.Println(cerr)
		}

		//update categorygroup
		uerr := Category.UpdateCategoryGroup(categories.CategoryCreate{
			Id:           1,
			CategoryName: "Default Group",
			CategorySlug: "default_group",
		})

		if uerr != nil {

			fmt.Println(uerr)
		}

		// delete categorygroup
		derr := Category.DeleteCategoryGroup(1,1)

		if derr != nil {

			fmt.Println(derr)
		}
	}

	cpermisison, _ := Auth.IsGranted("Categories", auth.CRUD)

	if cpermisison {

		//category list
		Categorylist, fnlist, parentcategory, count, err := Category.ListCategory(categories.CategoriesListReq{Limit: 10, Offset: 0})
		fmt.Println(Categorylist, fnlist, parentcategory, count, err)

		//create category
		cerr := Category.AddCategory(categories.CategoryCreate{
			CategoryName: "Default Category",
			 CategorySlug: "default_category",
			  ParentId: 1,
			})

			if cerr != nil {

				fmt.Println(cerr)
			}

		//update category
		uerr := Category.UpdateSubCategory(categories.CategoryCreate{
			Id: 5,
			CategoryName: "Default Category",
			CategorySlug: "default_category",
			ParentId: 2,
		})

		if uerr != nil {

			fmt.Println(uerr)
		}

		//delete category
		derr := Category.DeleteSubCategory(2, 1)

		if derr != nil {

			fmt.Println(derr)
		}

	}
}

```
# Getting help
If you encounter a problem with the package,please refer [Please refer [(https://www.spurtcms.com/documentation/cms-admin)] or you can create a new Issue in this repo[https://github.com/spurtcms/categories/issues]. 
