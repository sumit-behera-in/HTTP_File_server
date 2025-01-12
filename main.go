package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sumit-behera-in/HTTPFileServer/controller"
	"github.com/sumit-behera-in/HTTPFileServer/storage"
	"github.com/sumit-behera-in/goLogger"
)

func main() {
	server := gin.Default()
	logger, err := goLogger.NewLogger("HttpFileServer", "", 1, 4, "IST")
	if err != nil {
		panic(err)
	}

	storageOption := storage.StorageOptions{
		StorageRoot:       "HTTP_FILESERVER_PORT_4000",
		PathTransformFunc: storage.CASPathTransformFunc,
		Logger:            logger,
	}

	storageController := controller.NewStorageController(storageOption)
	storageController.RegisterRouterGroup(server.Group("/v1"))
	log.Fatal(server.Run(":4000"))
}
