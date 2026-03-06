package utils

import (
	"github.com/cloudwego/eino/schema"
	"github.com/wwwzy/CloudAI/model"
)

//---------------------------
//@brief 将多条Message model转为schema.Message
//---------------------------
func MessageList2ChatHistory(mess []*model.Message) (history []*schema.Message) {
	for _, m := range mess {
		history = append(history, message2MessagesTemplate(m))
	}
	return
}

//---------------------------
//@brief 将单条Message model转为schema.Message
//---------------------------
func message2MessagesTemplate(mess *model.Message) *schema.Message {
	return &schema.Message{
		Role:    schema.RoleType(mess.Role),
		Content: mess.Content,
	}
}