package pdf

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"code.sajari.com/docconv/v2"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
)

type options struct {
	toPages *bool
}

func WithToPages(toPages bool) parser.Option {
	return parser.WrapImplSpecificOptFn(func(opts *options) {
		opts.toPages = &toPages
	})
}

type Config struct {
	ToPages bool
}

type DocconvPDFParser struct {
	ToPages bool
}

// ---------------------------
// @brief 新建一个基于 docconv 库的 pdf 解析器
// ---------------------------
func NewDocconvPDFParser(ctx context.Context, config *Config) (*DocconvPDFParser, error) {
	if config == nil {
		config = &Config{}
	}
	return &DocconvPDFParser{ToPages: config.ToPages}, nil
}

// ---------------------------
// @brief 实现eino组件的Parse接口
// ---------------------------
func (pp *DocconvPDFParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	commonOpts := parser.GetCommonOptions(nil, opts...)

	specificOpts := parser.GetImplSpecificOptions(&options{
		toPages: &pp.ToPages,
	}, opts...)

	log.Println("解析PDF文档...")
	res, meta, err := docconv.ConvertPDF(reader)
	if err != nil {
		log.Printf("PDF解析错误: %v\n", err)
		return nil, fmt.Errorf("PDF解析失败: %w", err)
	}

	log.Printf("PDF解析完成，文本长度: %d字符\n", len(res))
	log.Printf("PDF元数据: %+v\n", meta)

	// 检查解析结果Parse
	if len(res) < 100 { // 至少需要100个字符才算有效
		log.Println("PDF解析结果太短或为空")
		if len(res) == 0 {
			return nil, fmt.Errorf("PDF解析结果为空，可能是扫描PDF或无文本内容")
		}
	}

	if *specificOpts.toPages {
		pages := strings.Split(res, "\f")
		if len(pages) == 1 {
			pages = []string{res}
		}
		documents := make([]*schema.Document, 0, len(pages))
		for _, p := range pages {
			pt := strings.TrimSpace(p)
			if len(pt) == 0 {
				continue
			}
			var sb strings.Builder
			sb.Grow(len(pt))
			sb.WriteString(pt)
			ptCopy := sb.String()

			meta := make(map[string]any, len(commonOpts.ExtraMeta))
			for k, v := range commonOpts.ExtraMeta {
				meta[k] = v
			}

			documents = append(documents, &schema.Document{
				Content:  ptCopy,
				MetaData: meta,
			})
		}
		return documents, nil
	}

	return []*schema.Document{{
		Content:  res,
		MetaData: commonOpts.ExtraMeta,
	}}, nil

}
