package response

import (
	"errors"
	"net/http"
	"net/url"
)

// Manager 响应管理器，提供统一的响应构建与写入能力。
// Manager 是线程安全的（内部 options 只读）。
type Manager struct {
	opts *Options
}

// NewManager 创建响应管理器。
func NewManager(opts ...Option) *Manager {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return &Manager{opts: o}
}

// --------------------------------------------------------------------
// Problem 构建 / 写入
// --------------------------------------------------------------------

// BuildProblem 根据 err 构建 RFC7807 Problem。
// 若 err 实现了 response.Error 接口，则从中提取信息；否则使用 fallback。
func (m *Manager) BuildProblem(err error, instance string) Problem {
	var respErr Error
	if errors.As(err, &respErr) {
		return m.buildFromRespError(respErr, instance)
	}
	// 普通 error，使用 fallback
	return Problem{
		Type:     m.opts.ProblemTypePrefix + m.opts.FallbackCode,
		Title:    m.opts.FallbackTitle,
		Status:   m.opts.FallbackStatus,
		Detail:   "",
		Instance: instance,
	}
}

// buildFromRespError 从 response.Error 构建 Problem。
func (m *Manager) buildFromRespError(e Error, instance string) Problem {
	code := e.Code()
	if code == "" {
		code = m.opts.FallbackCode
	}

	status := e.Status()
	if status == 0 {
		status = m.opts.FallbackStatus
	}

	title := e.Title()
	if title == "" {
		title = code
	}

	problemType := e.ProblemType()
	if problemType == "" {
		problemType = m.opts.ProblemTypePrefix + code
	}

	return Problem{
		Type:     problemType,
		Title:    title,
		Status:   status,
		Detail:   e.Detail(),
		Instance: instance,
	}
}

// WriteProblem 将 Problem 写入 JSONWriter，Content-Type 由 provider 设置为 application/problem+json。
// 注：核心包不关心 header，由 provider 在调用 w.JSON 前设置。
func (m *Manager) WriteProblem(w JSONWriter, err error, instance string) {
	p := m.BuildProblem(err, instance)
	w.JSON(p.Status, p)
}

// --------------------------------------------------------------------
// 成功响应
// --------------------------------------------------------------------

// WriteSuccess 写入成功响应。
// message 为空时使用默认消息。
func (m *Manager) WriteSuccess(w JSONWriter, message string, data any) {
	if message == "" {
		message = m.opts.DefaultSuccessMessage
	}
	w.JSON(http.StatusOK, SuccessBody{
		Message: message,
		Data:    data,
	})
}

// WritePaginated 写入分页成功响应。
func (m *Manager) WritePaginated(w JSONWriter, items any, total int64, page, pageSize int) {
	totalPages := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPages++
	}
	w.JSON(http.StatusOK, PaginatedBody{
		Message: m.opts.DefaultPaginatedMessage,
		Data: PaginatedData{
			Items:      items,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	})
}

// --------------------------------------------------------------------
// 重定向（支持 OAuth2 错误参数注入）
// --------------------------------------------------------------------

// BuildRedirectURL 构建重定向 URL。
// params 为自定义查询参数；若 err != nil 且配置了 OAuth2ErrorMapper，则注入 error/error_description。
func (m *Manager) BuildRedirectURL(rawURL string, err error, params map[string]string) string {
	parsed, parseErr := url.Parse(rawURL)
	if parseErr != nil {
		return rawURL
	}

	query := parsed.Query()

	// 添加自定义参数
	for k, v := range params {
		query.Set(k, v)
	}

	// 若存在错误且配置了 mapper，则注入 OAuth2 参数
	if err != nil && m.opts.OAuth2ErrorMapper != nil {
		if code, desc, ok := m.opts.OAuth2ErrorMapper(err); ok {
			query.Set("error", code)
			query.Set("error_description", desc)
		}
	}

	parsed.RawQuery = query.Encode()
	return parsed.String()
}

// WriteRedirect 执行重定向。
// status 应为 301/302/307/308；err 可为 nil。
func (m *Manager) WriteRedirect(w RedirectWriter, status int, rawURL string, err error, params map[string]string) {
	location := m.BuildRedirectURL(rawURL, err, params)
	w.Redirect(status, location)
}

// WriteRedirectTemporary 临时重定向（302）。
func (m *Manager) WriteRedirectTemporary(w RedirectWriter, rawURL string, err error, params map[string]string) {
	m.WriteRedirect(w, http.StatusFound, rawURL, err, params)
}

// WriteRedirectPermanent 永久重定向（301）。
func (m *Manager) WriteRedirectPermanent(w RedirectWriter, rawURL string, err error, params map[string]string) {
	m.WriteRedirect(w, http.StatusMovedPermanently, rawURL, err, params)
}
