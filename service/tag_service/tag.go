package tag_service

import (
	"blog-service/models"
	"blog-service/pkg/e"
	"blog-service/pkg/export"
	"blog-service/pkg/gredis"
	"blog-service/pkg/logging"
	"blog-service/service/cache_service"
	"encoding/json"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var tags []models.Tag
	var cacheTags []models.Tag
	cache := cache_service.Tag{
		State:    t.State,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cache.GetTagsKey()
	num, err := gredis.RedisClient.Exists(key).Result()
	if err != nil {
		logging.LogrusObj.Info(err)
	}
	if num > 0 {
		data, err := gredis.RedisClient.Get(key).Bytes()
		if err != nil {
			logging.LogrusObj.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	} else {
		code := e.ERROR_CHECK_EXIST_ARTICLE_FAIL
		logging.LogrusObj.Info(code)
		return nil, err
	}
	tags, err = models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}
	gredis.RedisClient.Set(key, tags, 3600)
	return tags, nil
}
func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}
	return maps
}

func (t *Tag) Count() (int64, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) Create() error {
	return models.CreateTags(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Update() error {
	data := make(map[string]interface{})
	data["name"] = t.Name
	data["modified_by"] = t.ModifiedBy
	if t.State >= 0 {
		data["state"] = t.State
	}
	return models.UpdateTags(t.ID, data)
}

func (t *Tag) Delete() error {
	return models.DeleteTags(t.ID)
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("tags")
	if err != nil {
		return "", err
	}
	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, v := range tags {
		values := []string{
			strconv.FormatInt(v.ID, 10),
			v.Name,
			v.CreatedBy,
			v.CreatedOn.String(),
			v.ModifiedBy,
			v.ModifiedOn.String(),
		}
		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}
	time := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags_" + time + ".xlsx"

	fullPath := export.GetExcelFullPath() + filename
	err = file.Save(fullPath)
	if err != nil {
		return "", err
	}
	return filename, nil
}
