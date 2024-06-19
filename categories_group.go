package categories

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

/*List Category Group*/
func (cat *Categories) CategoryGroupList(limit int, offset int, filter Filter) (Categorylist []TblCategories, categorycount int64, err error) {

	if AuthError := AuthandPermission(cat); AuthError != nil {

		return []TblCategories{}, 0, AuthError
	}

	Categorymodel.DataAccess = cat.DataAccess
	Categorymodel.Userid = cat.UserId

	_, Total_categories, _ := Categorymodel.CategoryGroupList(0, 0, filter, true, cat.DB)

	categorygrplist, _, cerr := Categorymodel.CategoryGroupList(offset, limit, filter, true, cat.DB)

	if cerr != nil {

		return []TblCategories{}, 0, cerr
	}

	return categorygrplist, Total_categories, nil

}

/*Add Category Group*/
func (cat *Categories) CreateCategoryGroup(req CategoryCreate) error {

	if AuthError := AuthandPermission(cat); AuthError != nil {

		return AuthError
	}

	if req.CategoryName == "" {

		return ErrorCategoryName
	}

	var category TblCategories

	category.CategoryName = req.CategoryName

	category.CategorySlug = strings.ToLower(strings.ReplaceAll(strings.TrimRight(req.CategoryName, " "), " ", "-"))

	category.Description = req.Description

	category.CreatedBy = req.CreatedBy

	category.ParentId = 0

	category.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := Categorymodel.CreateCategory(&category, cat.DB)

	if err != nil {

		return err
	}

	return nil

}

/*UpdateCategoryGroup*/
func (cat *Categories) UpdateCategoryGroup(req CategoryCreate) error {

	if AuthError := AuthandPermission(cat); AuthError != nil {

		return AuthError
	}

	if req.Id <= 0 || req.CategoryName == "" {

		return ErrorCategoryName
	}
	var category TblCategories

	category.Id = req.Id

	category.CategoryName = req.CategoryName

	category.Description = req.Description

	category.CategorySlug = strings.ToLower(strings.ReplaceAll(strings.TrimRight(req.CategoryName, " "), " ", "-"))

	category.ModifiedBy = req.ModifiedBy

	category.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := Categorymodel.UpdateCategory(&category, cat.DB)

	if err != nil {

		return err
	}

	return nil

}

/*DeleteCategoryGroup*/
func (cat *Categories) DeleteCategoryGroup(Categoryid int, modifiedby int) error {

	if AuthError := AuthandPermission(cat); AuthError != nil {

		return AuthError
	}

	GetData, _ := Categorymodel.GetCategoryTree(Categoryid, cat.DB)

	var individualid []int

	for _, GetParent := range GetData {

		indivi := GetParent.Id

		individualid = append(individualid, indivi)
	}

	derr := DeleteGroupchannelcategoryid(Categoryid, cat.DB)
	if derr != nil {
		fmt.Println(derr)
	}

	dcat := DeleteGroupEntriesCategoryid(Categoryid, cat.DB)
	if derr != nil {
		fmt.Println(dcat)
	}

	// spacecategory := individualid[0]

	var category TblCategories

	category.DeletedBy = modifiedby

	category.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	category.IsDeleted = 1

	err := Categorymodel.DeleteallCategoryById(&category, individualid, cat.DB)

	if err != nil {

		return err
	}

	return nil

}

/*DeleteCategoryGroup*/
func (cat *Categories) MultiSelectDeleteCategoryGroup(Categoryids []int, modifiedby int) error {

	var individualid []int

	for _, Categoryid := range Categoryids {

		GetData, _ := Categorymodel.GetCategoryTree(Categoryid, cat.DB)

		for _, GetParent := range GetData {

			indivi := GetParent.Id

			individualid = append(individualid, indivi)
		}

	}

	spacecategory := individualid

	multiChannelIds, multiRowids, err1 := Categorymodel.GetChannelCategoryids(&TblChannelCategorie{}, cat.DB)

	if err1 != nil {
		return err1
	}

	multiChannelIdsInt := make([][]int, len(multiChannelIds))

	for i, multiChannelId := range multiChannelIds {
		values := strings.Split(multiChannelId, ",")
		intValues := make([]int, len(values))
		for j, value := range values {
			intValues[j], _ = strconv.Atoi(value)
		}
		multiChannelIdsInt[i] = intValues
	}

	err1 = Categorymodel.MultiDeleteChannelCategoryids(&TblChannelCategorie{}, multiChannelIdsInt, multiRowids, Categoryids, cat.DB)

	if err1 != nil {
		return err1
	}

	multiEntryids, multiCategoryRowIds, multiEntryRowIds, err2 := Categorymodel.MultiGetEntryCategoryids(&TblChannelEntrie{}, Categoryids, cat.DB)
	if err2 != nil {
		return err2
	}

	var multiEntryIdInt [][]int

	for i, _ := range multiEntryids {
		values := multiEntryids[i]
		intValues := make([]int, len(values))
		if len(values) == 1 {
			intValues[0], _ = strconv.Atoi(values)
			multiEntryIdInt = append(multiEntryIdInt, intValues)
		} else if len(values) > 1 {
			splittedVal := strings.Split(values, ",")
			intValues := make([]int, len(splittedVal))
			for j, value := range splittedVal {
				intValues[j], _ = strconv.Atoi(value)
			}
			multiEntryIdInt = append(multiEntryIdInt, intValues)
		}

	}

	multiCategoryRowIdInt := make([]int, len(multiCategoryRowIds))
	multiEntryRowIdInt := make([]int, len(multiEntryRowIds))

	for i, val := range multiCategoryRowIds {
		multiCategoryRowIdInt[i], _ = strconv.Atoi(val)
	}

	for i, val := range multiEntryRowIds {
		multiEntryRowIdInt[i], _ = strconv.Atoi(val)
	}

	var multiUpdatedEntryId string
	var multiUpdateRowId int

	for _, rowid := range multiCategoryRowIdInt {

		for i, entryId := range multiEntryIdInt {
			ok, index := contains(entryId, rowid)

			if ok {
				multiUpdateRowId = multiEntryRowIdInt[i]
				var newSlice []int
				newSlice = append(newSlice, entryId[:index]...)
				newSlice = append(newSlice, entryId[index+1:]...)
				multiEntryIdInt[i] = newSlice
				multiUpdatedEntryId = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(newSlice)), ","), "[]")
				err2 := Categorymodel.DeleteEntryCategoryids(&TblChannelEntrie{}, multiUpdatedEntryId, multiUpdateRowId, cat.DB)
				if err2 != nil {
					return err2
				}
			}
		}
	}

	var category TblCategories

	category.DeletedBy = modifiedby

	category.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	category.IsDeleted = 1

	err := Categorymodel.DeleteallCategoryByIds(&category, Categoryids, spacecategory, cat.DB)

	if err != nil {

		return err
	}

	return nil

}

// Delete Multiselect subcategory//
func (cat *Categories) MultiselectSubCategoryDelete(Categoryids []int, modifiedby int) error {

	multiSubCategoryIds, multiSubRowIds, err1 := Categorymodel.MultiselectGetChannelCategoryids(&TblChannelCategorie{}, cat.DB)
	if err1 != nil {
		return err1
	}

	multiSubCategoryIdsInt := make([][]int, len(multiSubCategoryIds))

	for i, multiChannelId := range multiSubCategoryIds {
		values := strings.Split(multiChannelId, ",")
		intValues := make([]int, len(values))
		for j, value := range values {
			intValues[j], _ = strconv.Atoi(value)
		}
		multiSubCategoryIdsInt[i] = intValues
	}

	err1 = Categorymodel.MultiDeleteChannelCategoryids(&TblChannelCategorie{}, multiSubCategoryIdsInt, multiSubRowIds, Categoryids, cat.DB)
	if err1 != nil {
		return err1
	}

	multiSubEntryIds, multiSubCategoryRowIds, multiSubEntryRowIds, err2 := Categorymodel.MultiGetEntryCategoryids(&TblChannelEntrie{}, Categoryids, cat.DB)
	if err2 != nil {
		return err2
	}

	var multiSubEntryIdsInt [][]int

	for i, _ := range multiSubEntryIds {
		values := multiSubEntryIds[i]
		intValues := make([]int, len(values))
		if len(values) == 1 {
			intValues[0], _ = strconv.Atoi(values)
			multiSubEntryIdsInt = append(multiSubEntryIdsInt, intValues)
		} else if len(values) > 1 {
			splittedVal := strings.Split(values, ",")
			intValues := make([]int, len(splittedVal))
			for j, value := range splittedVal {
				intValues[j], _ = strconv.Atoi(value)
			}
			multiSubEntryIdsInt = append(multiSubEntryIdsInt, intValues)
		}

	}

	multiSubCategoryRowIdInt := make([]int, len(multiSubCategoryRowIds))
	multiSubEntryRowIdInt := make([]int, len(multiSubEntryRowIds))

	for i, val := range multiSubCategoryRowIds {
		multiSubCategoryRowIdInt[i], _ = strconv.Atoi(val)
	}

	for i, val := range multiSubEntryRowIds {
		multiSubEntryRowIdInt[i], _ = strconv.Atoi(val)
	}

	var multiSubUpdatedEntryId string
	var multiSubUpdateRowId int

	for _, rowid := range multiSubCategoryRowIdInt {

		for i, entryId := range multiSubEntryIdsInt {
			ok, index := contains(entryId, rowid)

			if ok {
				multiSubUpdateRowId = multiSubEntryRowIdInt[i]
				var newSlice []int
				newSlice = append(newSlice, entryId[:index]...)
				newSlice = append(newSlice, entryId[index+1:]...)
				multiSubEntryIdsInt[i] = newSlice
				multiSubUpdatedEntryId = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(newSlice)), ","), "[]")
				err2 := Categorymodel.DeleteEntryCategoryids(&TblChannelEntrie{}, multiSubUpdatedEntryId, multiSubUpdateRowId, cat.DB)
				if err2 != nil {
					return err2
				}
			}
		}
	}

	var category TblCategories

	category.DeletedBy = modifiedby

	category.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	category.IsDeleted = 1

	err := Categorymodel.DeleteCategoryByIds(&category, Categoryids, cat.DB)
	if err != nil {
		return err
	}

	return nil
}

func (cate CategoryModel) MultiselectGetChannelCategoryids(channelcategory *TblChannelCategorie, DB *gorm.DB) (categories []string, rowIds []string, err error) {

	// var categoryId string
	var categoryId []string
	var rowId []string

	result := DB.Debug().Table("tbl_channel_categories").Pluck("category_id", &categoryId)
	if result.Error != nil {
		return
	}

	result = DB.Debug().Table("tbl_channel_categories").Pluck("id", &rowId)
	if result.Error != nil {
		return
	}

	return categoryId, rowId, nil

}
