package ginmodel

// CreateFolderReq 创建文件夹请求
type CreateFolderReq struct {
	Name     string  `json:"name" binding:"required"` // 文件夹名称
	ParentID *string `json:"parent_id,omitempty"`     // 父文件夹ID(根目录传null)
}

// BatchMoveRequest 批量移动文件请求
type BatchMoveRequest struct {
	FileIDs        []string `json:"files_pid" binding:"required"` // 需要移动的文件/文件夹ID列表
	TargetParentID string   `json:"target_pid"`                   // 目标父文件夹ID
}

// RenameRequest 重命名文件请求
type RenameRequest struct {
	FileID  string `json:"file_id" binding:"required"`  // 文件ID
	NewName string `json:"new_name" binding:"required"` // 新名称(包含扩展名)
}
