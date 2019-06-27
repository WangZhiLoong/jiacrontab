package admin

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"

	"jiacrontab/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/cors"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/iwannay/log"
	"github.com/kataras/iris/context"
)

func newApp(adm *Admin) *iris.Application {

	app := iris.New()
	app.UseGlobal(newRecover(adm))
	app.Logger().SetLevel("debug")
	app.Use(logger.New())
	app.StaticEmbeddedGzip("/", "./assets/", GzipAsset, GzipAssetNames)
	cfg := adm.getOpts()

	wrapHandler := func(h func(ctx *myctx)) context.Handler {
		return func(c iris.Context) {
			h(wrapCtx(c, adm))
		}
	}

	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Jwt.SigningKey), nil
		},

		Extractor: func(ctx iris.Context) (string, error) {
			token, err := url.QueryUnescape(ctx.GetHeader(cfg.Jwt.Name))
			return token, err
		},
		Expiration: true,

		ErrorHandler: func(c iris.Context, data string) {
			ctx := wrapCtx(c, adm)
			if ctx.RequestPath(true) != "/user/login" && ctx.RequestPath(true) != "/user/signUp" {
				ctx.respAuthFailed(errors.New("认证失败"))
				return
			}
			ctx.Next()
		},

		SigningMethod: jwt.SigningMethodHS256,
	})

	crs := cors.New(cors.Options{
		Debug:            true,
		AllowedHeaders:   []string{"Content-Type", "Token"},
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
	})

	app.Use(crs)
	app.AllowMethods(iris.MethodOptions)
	app.Get("/", func(ctx iris.Context) {
		if adm.initModel {
			ctx.SetCookieKV("ready", "true", func(c *http.Cookie) {
				c.HttpOnly = false
			})
		} else {
			ctx.SetCookieKV("ready", "false", func(c *http.Cookie) {
				c.HttpOnly = false
			})
		}
		ctx.Header("Content-Type", "text/html")
		ctx.Header("Content-Encoding", "gzip")
		asset, err := GzipAsset("assets/index.html")
		if err != nil {
			log.Error(err)
		}
		ctx.WriteGzip(asset)
	})

	if adm.initModel {
		adm := app.Party("/adm")
		{
			adm.Use(jwtHandler.Serve)
			adm.Post("/crontab/job/list", wrapHandler(GetJobList))
			adm.Post("/crontab/job/get", wrapHandler(GetJob))
			adm.Post("/crontab/job/log", wrapHandler(GetRecentLog))
			adm.Post("/crontab/job/edit", wrapHandler(EditJob))
			adm.Post("/crontab/job/action", wrapHandler(ActionTask))
			adm.Post("/crontab/job/exec", wrapHandler(ExecTask))

			adm.Post("/config/get", wrapHandler(GetConfig))
			adm.Post("/config/mail/send", wrapHandler(SendTestMail))
			adm.Post("/system/info", wrapHandler(SystemInfo))

			adm.Post("/daemon/job/list", wrapHandler(GetDaemonJobList))
			adm.Post("/daemon/job/action", wrapHandler(ActionDaemonTask))
			adm.Post("/daemon/job/edit", wrapHandler(EditDaemonJob))
			adm.Post("/daemon/job/get", wrapHandler(GetDaemonJob))
			adm.Post("/daemon/job/log", wrapHandler(GetRecentDaemonLog))

			adm.Post("/group/list", wrapHandler(GetGroupList))
			adm.Post("/group/edit", wrapHandler(EditGroup))

			adm.Post("/node/list", wrapHandler(GetNodeList))
			adm.Post("/node/delete", wrapHandler(DeleteNode))
			adm.Post("/node/group_node", wrapHandler(GroupNode))

			adm.Post("/user/activity_list", wrapHandler(GetActivityList))
			adm.Post("/user/job_history", wrapHandler(GetJobHistory))
			adm.Post("/user/audit_job", wrapHandler(AuditJob))
			adm.Post("/user/stat", wrapHandler(UserStat))
			adm.Post("/user/signup", wrapHandler(Signup))
			adm.Post("/user/edit", wrapHandler(EditUser))
			adm.Post("/user/group_user", wrapHandler(GroupUser))
			adm.Post("/user/list", wrapHandler(GetUserList))
		}

		app.Post("/user/login", wrapHandler(Login))
	}

	app.Post("/app/init", wrapHandler(InitApp))

	debug := app.Party("/debug")
	{
		debug.Get("/stat", wrapHandler(stat))
		debug.Get("/pprof/", wrapHandler(indexDebug))
		debug.Get("/pprof/{key:string}", wrapHandler(pprofHandler))
	}

	return app
}

// InitApp 初始化应用
func InitApp(ctx *myctx) {
	var (
		err     error
		user    models.User
		reqBody InitAppReqParams
	)

	if err = ctx.Valid(&reqBody); err != nil {
		ctx.respParamError(err)
		return
	}

	if ret := models.DB().Take(&user, "group_id=?", 1); ret.Error == nil && ret.RowsAffected > 0 {
		ctx.respNotAllowed()
		return
	}

	user.Username = reqBody.Username
	user.Passwd = reqBody.Passwd
	user.Root = true
	user.GroupID = models.SuperGroup.ID
	user.Mail = reqBody.Mail

	if err = user.Create(); err != nil {
		ctx.respBasicError(err)
		return
	}

	ctx.respSucc("", true)
}
