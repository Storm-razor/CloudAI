package service

import (
	"context"

	"github.com/cloudwego/eino/schema"
	ginmodel "github.com/wwwzy/CloudAI/gin_model"
	"github.com/wwwzy/CloudAI/internal/dao"
	"github.com/wwwzy/CloudAI/model"
)

const (
	InputToQuery   = "InputToQuery"
	InputToHistory = "InputToHistory"
	ChatTemplate   = "ChatTemplate"
	ChatModel      = "ChatModel"
	Retriever      = "Retriever"
	Agent          = "Agent"
)

type AgentService interface {
	CreateAgent(ctx context.Context, agent *model.Agent) error
	UpdateAgent(ctx context.Context, agent *model.Agent) error
	DeleteAgent(ctx context.Context, userID uint, agentID string) error
	GetAgent(ctx context.Context, userID uint, agentID string) (*model.Agent, error)
	ListAgents(ctx context.Context, userID uint) ([]*model.Agent, error)
	PageAgents(ctx context.Context, userID uint, page, size int) ([]*model.Agent, int64, error)
	ExecuteAgent(ctx context.Context, userID uint, agentID string, msg ginmodel.UserMessage) (string, error)
	StreamExecuteAgent(ctx context.Context, userID uint, agentID string, msg ginmodel.UserMessage) (*schema.StreamReader[*schema.Message], error)
}

// TODO: 实现HistoryService
type agentService struct {
	dao      dao.AgentDao
	modelSvc ModelService
	kbSvc    KBService
	kbDao    dao.KnowledgeBaseDao
	modelDao dao.ModelDao
	//historySvc HistoryService
}

func NewAgentService(dao dao.AgentDao,
	modelSvc ModelService,
	kbSvc KBService,
	kbDao dao.KnowledgeBaseDao,
	modelDao dao.ModelDao,
	//historySvc HistoryService,
) AgentService {
	// return &agentService{
	// 	dao:        ctx,
	// 	modelSvc:   modelSvc,
	// 	kbSvc:      kbSvc,
	// 	kbDao:      kbDao,
	// 	modelDao:   modelDao,
	// 	historySvc: historySvc,
	// }
	return nil
}
