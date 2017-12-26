/* Gin Service with two endpoints:
- POST /pages-statistics (attach pdf file)
- POST /words-appearances (attach pdf file)
*/

package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/pages-statistics", func(c *gin.Context) {
		fileHeader, _ := c.FormFile("file")
		file, _ := fileHeader.Open()
		_, pagesStatistics, _ := processBookData(file)

		c.JSON(200, pagesStatistics)
	})

	router.POST("/words-appearances", func(c *gin.Context) {
		fileHeader, _ := c.FormFile("file")
		file, _ := fileHeader.Open()
		wordsAppearances, _, _ := processBookData(file)

		c.JSON(200, wordsAppearances)
	})

	return router
}

func main() {
	router := setupRouter()
	router.Run(":8080")
}
