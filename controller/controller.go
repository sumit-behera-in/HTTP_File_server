package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sumit-behera-in/HTTPFileServer/storage"
)

type StorageController struct {
	storage storage.Storage
}

func NewStorageController(StorageOptions storage.StorageOptions) *StorageController {
	return &StorageController{
		storage: *storage.NewStorage(StorageOptions),
	}
}

func (sc *StorageController) RegisterRouterGroup(rg *gin.RouterGroup) {
	userRoute := rg.Group("/fileserver")
	userRoute.GET("/:key", sc.readFile)
	userRoute.POST("/:key", sc.writeFile)
	userRoute.PATCH("/:key", sc.updateFile)
	userRoute.DELETE("/:key", sc.deleteFile)
}

func (sc *StorageController) readFile(ctx *gin.Context) {
	key := ctx.Param("key")
	if key == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "key is required to get a particular file"})
		return
	}

	reader, err := sc.storage.ReadStream(key)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("failed to fetching file with key : %s", key)})
		return
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("failed to Read file with error : %s", err.Error())})
		return
	}

	contentType := http.DetectContentType(data)

	// Set the appropriate Content-Type and return the file data
	ctx.Data(http.StatusOK, contentType, data)
}

func (sc *StorageController) writeFile(ctx *gin.Context) {
	key := ctx.Param("key")

	if key == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "key is required to get a particular file"})
		return
	}

	// Parse the file from the request body
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve the file"})
		return
	}

	// Open the uploaded file
	fileData, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to open the uploaded file"})
		return
	}
	defer fileData.Close()

	// check if it is written successfully
	if sc.storage.WriteStream(key, fileData) != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save the file content"})
		return
	}

	// Respond with success
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("File %s uploaded successfully", key)})

}

func (sc *StorageController) updateFile(ctx *gin.Context) {

}

func (sc *StorageController) deleteFile(ctx *gin.Context) {
	key := ctx.Param("key")
	if key == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "key is required to delete a particular file"})
		return
	}

	err := sc.storage.Delete(key)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("failed to fetching file with error : %s", err.Error())})
		return
	}

	// Respond with success
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("File %s deleted successfully", key)})
}
