package v1

import "github.com/gin-gonic/gin"

type Tag struct{}

func NewTag() *Tag {
	return &Tag{}
}

func (t *Tag) GetTag(c *gin.Context) {}

func (t *Tag) ListTag(c *gin.Context) {}

func (t *Tag) CreateTag(c *gin.Context) {}

func (t *Tag) UpdateTag(c *gin.Context) {}

func (t *Tag) DeleteTag(c *gin.Context) {}
