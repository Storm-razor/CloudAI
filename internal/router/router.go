package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wwwzy/CloudAI/internal/controller"
	"github.com/wwwzy/CloudAI/internal/middlerware"
)

func SetUpRouters(r *gin.Engine,
	uc *controller.UserController,
	fc *controller.FileController,
	kc *controller.KBController,
	mc *controller.ModelController,
	ac *controller.AgentController,
	cc *controller.ConversationController) {
	api := r.Group("/api")
	{
		publicUser := api.Group("/user")
		{
			publicUser.POST("/register", uc.Register)
			publicUser.POST("/login", uc.Login)
		}

		file := api.Group("/files")
		file.Use(middlerware.JWTAuth())
		{
			file.POST("/upload", fc.Upload)
			file.GET("/page", fc.PageList)
			file.GET("/download", fc.Download)
			file.DELETE("/delete", fc.Delete)
			file.POST("/folder", fc.CreateFolder)
			file.POST("/move", fc.BatchMove)
			file.GET("/search", fc.Search)
			file.PUT("/rename", fc.Rename)
			file.GET("/path", fc.GetPath)
			file.GET("/idPath", fc.GetIDPath)
		}

		kb := api.Group("knowledge")
		kb.Use(middlerware.JWTAuth())
		{
			// KB
			kb.POST("/create", kc.Create)
			kb.DELETE("/delete", kc.Delete)
			kb.POST("/add", kc.AddExistFile)
			kb.POST("/addNew", kc.AddNewFile)
			kb.GET("/page", kc.PageList)
			kb.GET("/detail", kc.GetKBDetail)
			// Doc
			kb.GET("/docPage", kc.DocPage)
			kb.POST("/docDelete", kc.DeleteDocs)
			// RAG
			kb.POST("/retrieve", kc.Retrieve)
			kb.POST("/chat", kc.Chat)
			kb.POST("/stream", kc.ChatStream)
		}

		model := api.Group("model")
		model.Use(middlerware.JWTAuth())
		{
			model.POST("/create", mc.CreateModel)
			model.PUT("/update", mc.UpdateModel)
			model.DELETE("/delete", mc.DeleteModel)
			model.GET("/get", mc.GetModel)
			model.GET("/page", mc.PageModels)
			model.GET("/list", mc.ListModels)
		}

		agent := api.Group("agent")
		agent.Use(middlerware.JWTAuth())
		{
			agent.POST("/create", ac.CreateAgent)
			agent.POST("/update", ac.UpdateAgent)
			agent.DELETE("/delete", ac.DeleteAgent)
			agent.GET("/get", ac.GetAgentDetail)
			agent.GET("/page", ac.PageAgents)
			agent.POST("/execute/:id", ac.ExecuteAgent)
			agent.POST("/stream", ac.StreamExecuteAgent)
		}

		conv := api.Group("chat")
		conv.Use(middlerware.JWTAuth())
		{
			// 调试模式，不保存历史
			conv.POST("/debug", cc.DebugStreamAgent)
			// 会话相关功能
			conv.POST("/create", cc.CreateConversation)
			conv.POST("/stream", cc.StreamConversation)
			conv.GET("/list", cc.ListConversations)
			conv.GET("/list/agent", cc.ListAgentConversations)
			conv.GET("/history", cc.GetConversationHistory)
			conv.DELETE("/delete", cc.DeleteConversation)
		}
	}
}
