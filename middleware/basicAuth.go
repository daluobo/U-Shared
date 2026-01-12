package middleware

import (
	"encoding/base64"
	"strings"

	"go-shared/config"

	"github.com/kataras/iris/v12"
)

// BasicAuth 中间件
func BasicAuth(ctx iris.Context) {
	if !config.AppConfig.BasicAuth.Enabled {
		ctx.Next()
		return
	}

	// 获取Authorization头
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		askForCredentials(ctx)
		return
	}

	// 检查Basic认证格式
	if !strings.HasPrefix(authHeader, "Basic ") {
		askForCredentials(ctx)
		return
	}

	// 解码Base64
	encoded := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid authorization header"})
		return
	}

	// 解析用户名和密码
	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		askForCredentials(ctx)
		return
	}

	username, password := credentials[0], credentials[1]

	// 验证用户
	if !validateUser(username, password) {
		askForCredentials(ctx)
		return
	}

	// 将用户信息存储到上下文中
	ctx.Values().Set("auth_user", username)
	ctx.Next()
}

func validateUser(username, password string) bool {
	if config.AppConfig.BasicAuth.Enabled == false {
		return true
	}

	for _, user := range config.AppConfig.BasicAuth.Users {
		if user.Username == username {
			return password == user.Password
		}
	}
	return false
}

func askForCredentials(ctx iris.Context) {
	ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
	ctx.StatusCode(iris.StatusUnauthorized)
	ctx.JSON(iris.Map{
		"error": "Authentication required",
		"code":  "UNAUTHORIZED",
	})
	ctx.StopExecution()
}
