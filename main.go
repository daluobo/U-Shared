package main

import (
	_ "embed"
	"u-shared/config"
	"u-shared/controller"
	"u-shared/middleware"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/getlantern/systray"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/pio"
)

// 生成图标 windres -o app-icon.syso app-icon.rc
// 打包 命令 go build -ldflags="-H windowsgui"

var (
	app *iris.Application
)

func main() {
	app = iris.New()

	app.Use(logger.New())
	app.Use(recover.New())

	app.Logger().SetLevel("debug")
	initLog()

	initFolder()

	err := config.LoadConfig("config.yaml")
	if err != nil {
		golog.Fatal("Failed to load config:", err)
	}

	mvcApp := mvc.New(app)

	appParty := mvcApp.Party("/app")
	appParty.Router.Use(middleware.BasicAuth)
	appParty.Handle(new(controller.AppController))

	sharedParty := mvcApp.Party("/shared")
	sharedParty.Router.Use(middleware.BasicAuth, middleware.PathValidator)
	sharedParty.Handle(new(controller.SharedController))

	runtime.LockOSThread()
	systray.Run(onReady, onExit)
}

func initLog() {
	golog.DebugText("[ DBUG ]", pio.Green)
	golog.InfoText("[ INFO ]", pio.Cyan)
	golog.WarnText("[ WARN ]", pio.Yellow)
	golog.ErrorText("[ ERRO ]", pio.Red, pio.Reversed)
}

func initFolder() {
	for _, folder := range middleware.MediaFolders {
		_, err := os.Stat("./shared/" + folder)
		if err != nil {
			if os.IsNotExist(err) {
				golog.Infof("Folder does not exist, creating:", folder)
				err = os.MkdirAll("./shared/"+folder, 0755)
				if err != nil {
					golog.Errorf("Error creating folder:", err)
				}
			} else {
				golog.Errorf("Error checking folder:", err)
			}
		}
	}
}

func onReady() {
	// 设置托盘图标（Windows 需要 .ico 格式）
	systray.SetIcon(loadIcon()) // iconData 是 []byte 格式的图标
	systray.SetTitle("Shared Server")
	systray.SetTooltip("Shared Server")

	// 添加菜单项
	mQuit := systray.AddMenuItem("退出", "退出")

	// 监听菜单点击
	go func() {
		for range mQuit.ClickedCh {
			systray.Quit()                     // 先退出系统托盘
			time.Sleep(100 * time.Millisecond) // 给一点时间清理
			os.Exit(0)                         // 然后退出
		}
	}()

	app.Run(iris.Addr(":" + strconv.Itoa(config.AppConfig.Port)))
}

func onExit() {
	// 清理代码
}

func loadIcon() []byte {
	data, err := os.ReadFile("app-icon.ico")
	if err != nil {
		// 处理错误或使用内置图标
	}
	return data
}
