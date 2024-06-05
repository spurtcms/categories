package categories

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// Check category name already exists
func (cate *Categories) CheckCategroyGroupName(id int, name string) (bool, error) {

	var category TblCategories

	err := Categorymodel.CheckCategoryGroupName(category, id, name, cate.DB)

	if err != nil {

		return false, err

	}

	return true, nil
}

/*Get All cateogry with parents and subcategory*/
func (cate *Categories) AllCategoriesWithSubList() (arrangecategories []Arrangecategories, CategoryNames []string) {

	getallparentcat, _ := Categorymodel.GetAllParentCategory(cate.DB)

	var AllCategorieswithSubCategories []Arrangecategories

	for _, Group := range getallparentcat {

		GetData, _ := Categorymodel.GetCategoryTree(Group.Id, cate.DB)

		var pid int

		for _, categories := range GetData {

			var addcat Arrangecategories

			var individualid []CatgoriesOrd

			pid = categories.ParentId

		LOOP:
			for _, GetParent := range GetData {

				var indivi CatgoriesOrd

				if pid == GetParent.Id {

					pid = GetParent.ParentId

					indivi.Id = GetParent.Id

					indivi.Category = GetParent.CategoryName

					individualid = append(individualid, indivi)

					if pid != 0 {

						goto LOOP //this loop get looped until parentid=0

					}
				}

			}

			var ReverseOrder Arrangecategories

			addcat.Categories = append(addcat.Categories, individualid...)

			var singlecat []CatgoriesOrd

			for i := len(addcat.Categories) - 1; i >= 0; i-- {

				var Sing CatgoriesOrd

				Sing.Id = addcat.Categories[i].Id

				Sing.Category = addcat.Categories[i].Category

				singlecat = append(singlecat, Sing)

			}

			var newcate CatgoriesOrd

			newcate.Id = categories.Id

			newcate.Category = categories.CategoryName

			addcat.Categories = append(addcat.Categories, newcate)

			singlecat = append(singlecat, newcate)

			ReverseOrder.Categories = singlecat

			AllCategorieswithSubCategories = append(AllCategorieswithSubCategories, ReverseOrder)
		}

	}

	/*This for Channel category show also individual group*/
	var FinalCategoryList []Arrangecategories

	for _, val := range AllCategorieswithSubCategories {

		if len(val.Categories) > 1 {

			var infinalarray Arrangecategories

			infinalarray.Categories = append(infinalarray.Categories, val.Categories...)

			FinalCategoryList = append(FinalCategoryList, infinalarray)
		}

	}

	var Categorynames []string

	for _, val := range FinalCategoryList {

		var name string

		for index, cat := range val.Categories {

			if len(val.Categories)-1 == index {

				name += cat.Category
			} else {
				name += cat.Category + " / "
			}

		}

		Categorynames = append(Categorynames, name)

	}

	return FinalCategoryList, Categorynames

}

// Check Sub category name already exists
func (cate *Categories) CheckSubCategroyName(id []int, Currentcategoryid int, name string) (bool, error) {

	if Autherr := AuthandPermission(cate); Autherr != nil {

		return false, Autherr
	}

	_, err := Categorymodel.CheckSubCategoryName(id, Currentcategoryid, name, cate.DB)

	if err != nil {

		return false, err
	}

	return true, nil
}

func contains(slice []int, item int) (check bool, index int) {

	for i, val := range slice {
		if val == item {
			return true, i
		}
	}

	return false, -1
}

func (cate *Categories) DeleteChannelsubCategories(db *gorm.DB, categoryid int) error {

	subcategoryIds, subrowId, err1 := Categorymodel.GetChannelCategoryids(&TblChannelCategorie{}, db)
	if err1 != nil {
		return err1
	}

	subCategoryIdInt := make([][]int, len(subcategoryIds))

	for i, subCategoryId := range subcategoryIds {
		values := strings.Split(subCategoryId, ",")
		intValues := make([]int, len(values))
		for j, value := range values {
			intValues[j], _ = strconv.Atoi(value)
		}
		subCategoryIdInt[i] = intValues
	}

	err1 = Categorymodel.DeleteChannelCategoryids(&TblChannelCategorie{}, subCategoryIdInt, subrowId, categoryid, db)
	if err1 != nil {
		return err1
	}

	return nil
}

func (cate *Categories) DeleteEntriessubCategories(db *gorm.DB, categoryid int) error {

	subEntryIds, subCategoryRowIds, subEntryRowIds, err2 := Categorymodel.GetEntryCategoryids(&TblChannelEntrie{}, categoryid, db)
	if err2 != nil {
		return err2
	}

	var subEntriesIdInt [][]int

	for i, _ := range subEntryIds {
		values := subEntryIds[i]
		intValues := make([]int, len(values))
		if len(values) == 1 {
			intValues[0], _ = strconv.Atoi(values)
			subEntriesIdInt = append(subEntriesIdInt, intValues)
		} else if len(values) > 1 {
			splittedVal := strings.Split(values, ",")
			intValues := make([]int, len(splittedVal))
			for j, value := range splittedVal {
				intValues[j], _ = strconv.Atoi(value)
			}
			subEntriesIdInt = append(subEntriesIdInt, intValues)
		}
	}

	subCategoryRowIdInt := make([]int, len(subCategoryRowIds))
	subEntryRowIdInt := make([]int, len(subEntryRowIds))

	for i, val := range subCategoryRowIds {
		subCategoryRowIdInt[i], _ = strconv.Atoi(val)
	}

	for i, val := range subEntryRowIds {
		subEntryRowIdInt[i], _ = strconv.Atoi(val)
	}

	var updatedSubEntryId string
	var updatedSubRowId int

	for _, rowId := range subCategoryRowIdInt {

		for i, subEntryIdInt := range subEntriesIdInt {

			ok, index := contains(subEntryIdInt, rowId)
			fmt.Println("ok", ok)

			if ok {
				updatedSubRowId = subEntryRowIdInt[i]
				var newSlice []int
				newSlice = append(newSlice, subEntryIdInt[:index]...)
				newSlice = append(newSlice, subEntryIdInt[index+1:]...)
				subEntriesIdInt[i] = newSlice
				updatedSubEntryId = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(newSlice)), ","), "[]")
				err2 := Categorymodel.DeleteEntryCategoryids(&TblChannelEntrie{}, updatedSubEntryId, updatedSubRowId, db)
				if err2 != nil {
					fmt.Println("err2:", err2)
					return err2
				}

			}

		}
	}
	return nil
}

func DeleteGroupchannelcategoryid(Categoryid int, DB *gorm.DB) error {

	channelIds, rowIds, err1 := Categorymodel.GetChannelCategoryids(&TblChannelCategorie{}, DB)
	if err1 != nil {
		return err1
	}

	channelIdInt := make([][]int, len(channelIds))

	for i, channelId := range channelIds {
		values := strings.Split(channelId, ",")
		intValues := make([]int, len(values))
		for j, value := range values {
			intValues[j], _ = strconv.Atoi(value)
		}
		channelIdInt[i] = intValues
	}

	err1 = Categorymodel.DeleteChannelCategoryids(&TblChannelCategorie{}, channelIdInt, rowIds, Categoryid, DB)
	if err1 != nil {
		return err1
	}

	return nil
}

func DeleteGroupEntriesCategoryid(Categoryid int, DB *gorm.DB) error {

	entryIds, categoryRowIds, entryRowIds, err2 := Categorymodel.GetEntryCategoryids(&TblChannelEntrie{}, Categoryid, DB)
	if err2 != nil {

		return err2
	}

	var entryIdInt [][]int

	for i, _ := range entryIds {
		values := entryIds[i]
		intValues := make([]int, len(values))
		if len(values) == 1 {
			intValues[0], _ = strconv.Atoi(values)
			entryIdInt = append(entryIdInt, intValues)
		} else if len(values) > 1 {
			splittedVal := strings.Split(values, ",")
			intValues := make([]int, len(splittedVal))
			for j, value := range splittedVal {
				intValues[j], _ = strconv.Atoi(value)
			}
			entryIdInt = append(entryIdInt, intValues)
		}

	}

	categoryRowIdInt := make([]int, len(categoryRowIds))
	entryRowIdInt := make([]int, len(entryRowIds))

	for i, val := range categoryRowIds {
		categoryRowIdInt[i], _ = strconv.Atoi(val)
	}

	for i, val := range entryRowIds {
		entryRowIdInt[i], _ = strconv.Atoi(val)
	}

	var updatedEntryId string
	var updateRowId int

	for _, rowid := range categoryRowIdInt {

		for i, entryId := range entryIdInt {
			ok, index := contains(entryId, rowid)

			if ok {
				updateRowId = entryRowIdInt[i]
				var newSlice []int
				newSlice = append(newSlice, entryId[:index]...)
				newSlice = append(newSlice, entryId[index+1:]...)
				entryIdInt[i] = newSlice
				updatedEntryId = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(newSlice)), ","), "[]")
				err2 := Categorymodel.DeleteEntryCategoryids(&TblChannelEntrie{}, updatedEntryId, updateRowId, DB)
				if err2 != nil {
					return err2
				}
			}
		}
	}
	return nil
}
