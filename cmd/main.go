package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/wwwzy/CloudAI/config"
	"github.com/wwwzy/CloudAI/internal/controller"
	"github.com/wwwzy/CloudAI/internal/dao"
	"github.com/wwwzy/CloudAI/internal/dao/history"
	"github.com/wwwzy/CloudAI/internal/database"
	"github.com/wwwzy/CloudAI/internal/middlerware"
	"github.com/wwwzy/CloudAI/internal/router"
	"github.com/wwwzy/CloudAI/internal/service"
)

func main() {
	config.InitConfig()
	ctx := context.Background()

	db, _ := database.GetDB()

	userDao := dao.NewUserDao(db)
	userService := service.NewUserService(userDao)
	userController := controller.NewUserController(userService)

	fileDao := dao.NewFileDao(db)
	fileService := service.NewFileService(fileDao)
	fileController := controller.NewFileController(fileService)

	milvusClient, _ := database.InitMilvus(ctx)
	defer milvusClient.Close()

	modelDao := dao.NewModelDao(db)
	modelService := service.NewModelService(modelDao)
	modelController := controller.NewModelController(modelService)

	kbDao := dao.NewKnowledgeBaseDao(db)
	kbService := service.NewKBService(kbDao, fileService, modelDao)
	kbController := controller.NewKBController(kbService, fileService)

	msgDao := history.NewMsgDao(db)
	convDao := history.NewConvDao(db)
	historyService := service.NewHistoryService(convDao, msgDao)

	agentDao := dao.NewAgentDao(db)
	agentService := service.NewAgentService(agentDao, modelService, kbService, kbDao, modelDao, historyService)
	agentController := controller.NewAgentController(agentService)

	conversationService := service.NewConversationService(agentService, historyService)
	conversationController := controller.NewConversationController(conversationService)

	r := gin.Default()
	r.Use(middlerware.SetupCORS()) //配置跨域请求
	router.SetUpRouters(r,
		userController,
		fileController,
		kbController,
		modelController,
		agentController,
		conversationController,
	)

	port := config.GetConfig().Server.Port
	r.Run(":" + port)
}
