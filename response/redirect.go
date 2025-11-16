package response

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/3086953492/gokit/errors"
)

// RedirectTemporary 临时重定向（HTTP 302）
// c: gin 上下文
// targetURL: 目标重定向地址
// err: 可选的错误对象，若非 nil 则会在 URL 中添加 error 和 error_description 参数
// params: 自定义查询参数（如 map[string]string{"code":"123","state":"ok"}），会覆盖原 URL 中的同名参数
//
// 参数优先级（从低到高）：
//  1. 原 URL 中的查询参数
//  2. params 中的自定义参数（覆盖原 URL 中的同名参数）
//  3. 错误相关参数 error 和 error_description（覆盖 params 中的同名参数）
//
// 若暂不需要额外参数，可传入 nil 或 map[string]string{}
func RedirectTemporary(c *gin.Context, targetURL string, err error, params map[string]string) {
	appErr := normalizeError(err)
	finalURL := buildRedirectURL(targetURL, appErr, params)
	c.Redirect(http.StatusFound, finalURL)
}

// RedirectPermanent 永久重定向（HTTP 301）
// c: gin 上下文
// targetURL: 目标重定向地址
// err: 可选的错误对象，若非 nil 则会在 URL 中添加 error 和 error_description 参数
// params: 自定义查询参数（如 map[string]string{"code":"123","state":"ok"}），会覆盖原 URL 中的同名参数
//
// 参数优先级（从低到高）：
//  1. 原 URL 中的查询参数
//  2. params 中的自定义参数（覆盖原 URL 中的同名参数）
//  3. 错误相关参数 error 和 error_description（覆盖 params 中的同名参数）
//
// 若暂不需要额外参数，可传入 nil 或 map[string]string{}
func RedirectPermanent(c *gin.Context, targetURL string, err error, params map[string]string) {
	appErr := normalizeError(err)
	finalURL := buildRedirectURL(targetURL, appErr, params)
	c.Redirect(http.StatusMovedPermanently, finalURL)
}

// buildRedirectURL 构建带参数的重定向 URL
// 支持添加自定义参数和错误信息到目标 URL 的查询参数中，并自动进行 URL 编码
//
// rawURL: 原始目标 URL
// appErr: 可选的错误对象，若非 nil 则会添加 error 和 error_description 参数
// params: 自定义查询参数，会覆盖原 URL 中的同名参数
//
// 返回构建好的完整 URL 字符串
func buildRedirectURL(rawURL string, appErr *errors.AppError, params map[string]string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		// 解析失败时返回原始 URL，保证容错
		return rawURL
	}

	// 获取现有查询参数
	query := parsedURL.Query()

	// 添加自定义参数（会覆盖原有同名参数）
	for k, v := range params {
		query.Set(k, v)
	}

	// 若存在错误信息，添加错误相关参数（优先级最高，会覆盖 params 中的同名参数）
	if appErr != nil {
		cfg := getConfig()
		errorCode := cfg.OAuthErrorCodeMapper(appErr.Type)
		query.Set("error", errorCode)
		query.Set("error_description", appErr.Message)
	}

	// 将修改后的查询参数编码并写回 URL
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String()
}

// normalizeError 规范化错误对象
// 将普通 error 包装为 AppError，如果已经是 AppError 则直接返回
// 如果 err 为 nil，则返回 nil
func normalizeError(err error) *errors.AppError {
	if err == nil {
		return nil
	}

	var appErr *errors.AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	// 若 err 不是 AppError，则包装为内部错误
	cfg := getConfig()
	return errors.Internal().Msg(cfg.FallbackErrorMessage).Err(err).Build()
}
