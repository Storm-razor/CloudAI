package milvus

import "github.com/wwwzy/CloudAI/pkgs/consts"

var (
	defaultSearchFields = []string{
		consts.FieldNameID,
		consts.FieldNameContent,
		consts.FieldNameKBID,
		consts.FieldNameDocumentID,
		consts.FieldNameMetadata,
	}
)
