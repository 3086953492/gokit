package problem

import (
	"github.com/gin-gonic/gin"
)

// Fail 输出 RFC 7807 Problem Details 响应。
//
// 调用方显式传入 status/title/detail/ptype，problem 只负责组装并输出。
//
//	problem.Fail(c, 403, "Forbidden", "no permission", "https://api.example.com/errors/forbidden")
//	problem.Fail(c, 400, "Bad Request", "invalid params", "") // ptype 为空则使用 "about:blank"
func Fail(c *gin.Context, status int, title, detail, ptype string, opts ...Option) {
	o := defaultOptions()
	for _, fn := range opts {
		fn(o)
	}

	if ptype == "" {
		ptype = "about:blank"
	}

	p := &Problem{
		Type:       ptype,
		Title:      title,
		Status:     status,
		Detail:     detail,
		Extensions: o.Extensions,
	}

	writeProblem(c, p, o)
}

// writeProblem 写入 Problem 响应（header + body）
func writeProblem(c *gin.Context, p *Problem, o *Options) {
	if o.Instance != "" {
		p.Instance = o.Instance
	} else {
		p.Instance = c.Request.URL.Path
	}

	// 如果有扩展字段，需要展平到顶层
	if len(p.Extensions) > 0 {
		body := map[string]any{
			"type":     p.Type,
			"title":    p.Title,
			"status":   p.Status,
			"detail":   p.Detail,
			"instance": p.Instance,
		}
		for k, v := range p.Extensions {
			body[k] = v
		}
		c.Header("Content-Type", ContentType)
		c.JSON(p.Status, body)
		return
	}

	c.Header("Content-Type", ContentType)
	c.JSON(p.Status, p)
}

