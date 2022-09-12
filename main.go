package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Article struct {
	gorm.Model
	Title string
	Slug  string `gorm:"unique_index"`
	Desc  string `sql:"type:text;"`
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open("mysql", "root:@tcp(localhost:3306)/gin_gorm?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	//Migrate schema
	db.AutoMigrate(&Article{})

	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		articles := v1.Group("/article")
		{
			articles.GET("/", getHome)
			articles.GET("/:slug", getArticle)
			articles.POST("/", postArticle)
		}

	}
	router.Run()
}

func getHome(ctx *gin.Context) {
	items := []Article{}
	db.Find(&items)

	ctx.JSON(200, gin.H{
		"status": "Success",
		"data":   items,
	})
}

func getArticle(ctx *gin.Context) {
	slug := ctx.Param("slug")
	var item Article

	if db.First(&item, "slug = ?", slug).RecordNotFound() {
		ctx.JSON(404, gin.H{
			"status":  "Error",
			"message": "record not found",
		})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{
		"status": "Success",
		"data":   item,
	})
}

func postArticle(ctx *gin.Context) {
	item := Article{
		Title: ctx.PostForm("title"),
		Desc:  ctx.PostForm("desc"),
		Slug:  slug.Make(ctx.PostForm("title")),
	}

	db.Create(&item)

	ctx.JSON(200, gin.H{
		"status": "Success post",
		"data":   item,
	})
}
