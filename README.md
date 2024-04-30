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


# Create new category group


# Getting help
If you encounter a problem with the package,please refer [Please refer [(https://www.spurtcms.com/documentation/cms-admin)] or you can create a new Issue in this repo[https://github.com/spurtcms/categories/issues]. 
