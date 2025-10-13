package controller

import (
	"bless-activity/model"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

const (
	ctxFishpiLoginUser = "login_user"
	ctxFishpiNext      = "next"
	ctxFishpiOpenId    = "openid"
	ctxFishpiUserInfo  = "fishpi_user_info"
)

type FishPiController struct {
	event *core.ServeEvent
	app   core.App

	logger *slog.Logger
}

func NewFishPiController(event *core.ServeEvent) *FishPiController {
	logger := event.App.Logger().With(
		slog.String("controller", "fishpi"),
	)

	controller := &FishPiController{
		event:  event,
		app:    event.App,
		logger: logger,
	}

	controller.registerRoutes()

	return controller
}

func (controller *FishPiController) registerRoutes() {
	fishpiGroup := controller.event.Router.Group("/fishpi")
	fishpiGroup.GET("/login", controller.Login)
	fishpiGroup.GET("/callback", controller.Callback).BindFunc(
		controller.CallbackVerify,
	)
}

func (controller *FishPiController) makeActionLogger(action string) *slog.Logger {
	return controller.logger.With(
		slog.String("action", action),
	)
}

func (controller *FishPiController) Login(event *core.RequestEvent) error {

	appUrl := event.App.Settings().Meta.AppURL
	callbackUrl := fmt.Sprintf("%s/fishpi/callback", appUrl)

	query := url.Values{}
	query.Set("openid.ns", "http://specs.openid.net/auth/2.0")
	query.Set("openid.mode", "checkid_setup")
	query.Set("openid.return_to", callbackUrl)
	query.Set("openid.realm", appUrl)
	query.Set("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")
	query.Set("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")

	addr := url.URL{
		Scheme:   "https",
		Host:     "fishpi.cn",
		Path:     "/openid/login",
		RawQuery: query.Encode(),
	}

	return event.Redirect(http.StatusFound, addr.String())
}

func (controller *FishPiController) CallbackVerify(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("callback verify").With(
		slog.String("path", event.Request.URL.String()),
	)

	info, err := event.RequestInfo()
	if err != nil {
		logger.Error("获取请求信息失败", slog.Any("err", err))
		return err
	}

	query := info.Query
	query["openid.mode"] = "check_authentication"

	resp := new(req.Response)
	if resp, err = req.C().R().
		SetBodyJsonMarshal(query).
		Post("https://fishpi.cn/openid/verify"); err != nil {
		logger.Error("发起验证请求失败", slog.Any("err", err))
		return err
	}
	valid := false
	arr := strings.Split(resp.String(), "\n")
	for _, line := range arr {
		if strings.HasPrefix(line, "is_valid:") {
			valid = strings.TrimPrefix(line, "is_valid:") == "true"
			break
		}
	}
	if !valid {
		logger.Error("验证失败", slog.String("resp", resp.String()))
		return errors.New("用户信息无效")
	}
	identity := query["openid.identity"]
	openid := filepath.Base(identity)

	result := new(FishpiUserInfoResult)
	if resp, err = req.C().R().
		SetSuccessResult(result).
		Get(fmt.Sprintf("https://fishpi.cn/api/user/getInfoById?userId=%s", openid)); err != nil {
		logger.Error("发起获取用户信息请求失败", slog.Any("err", err))
		return err
	}

	if result.Code != 0 {
		logger.Error("获取用户信息失败", slog.String("resp", resp.String()))
		return errors.New(result.Msg)
	}

	user := new(model.User)
	if err = event.App.RecordQuery(model.DbNameUsers).Where(dbx.HashExp{model.UsersFieldOId: openid}).One(user); err == nil {
		event.Set(ctxFishpiLoginUser, user)
		event.Set(ctxFishpiUserInfo, result.Data)
		event.Set(ctxFishpiNext, "login")
		return event.Next()
	} else if !errors.Is(err, sql.ErrNoRows) {
		logger.Error("查询用户信息失败", slog.Any("err", err))
		return err
	}

	event.Set(ctxFishpiOpenId, openid)
	event.Set(ctxFishpiUserInfo, result.Data)
	event.Set(ctxFishpiNext, "register")

	return event.Next()
}

type FishpiUserInfoResult struct {
	Msg  string          `json:"msg"`
	Code int             `json:"code"`
	Data *FishpiUserInfo `json:"data"`
}

type FishpiUserInfo struct {
	UserAvatarURL string `json:"userAvatarURL"`
	UserNickname  string `json:"userNickname"`
	UserName      string `json:"userName"`
}

func (fishpiUserInfo *FishpiUserInfo) Name() string {
	return fmt.Sprintf("(%s)%s", fishpiUserInfo.UserName, fishpiUserInfo.UserNickname)
}

func (controller *FishPiController) Callback(event *core.RequestEvent) error {
	if event.Get(ctxFishpiNext) == "login" {
		return controller.login(event)
	}
	return controller.register(event)
}

func (controller *FishPiController) login(event *core.RequestEvent) error {
	logger := controller.makeActionLogger("callback login").With(
		slog.String("path", event.Request.URL.String()),
	)
	user := event.Get(ctxFishpiLoginUser).(*model.User)
	fishpiUserInfo := event.Get(ctxFishpiUserInfo).(*FishpiUserInfo)

	logger = logger.With(slog.String("id", user.Id), slog.String("name", user.GetString("name")))

	// 更新用户资料
	if fishpiUserInfo.Name() != user.Name() || fishpiUserInfo.UserAvatarURL != user.Avatar() {
		if fishpiUserInfo.Name() != user.Name() {
			user.SetName(fishpiUserInfo.Name())
		}
		if fishpiUserInfo.UserAvatarURL != user.Avatar() {
			user.SetAvatar(fishpiUserInfo.UserAvatarURL)
		}

		if err := event.App.Save(user); err != nil {
			logger.Error("更新用户资料失败", slog.Any("user", user), slog.Any("fishpi_user_info", fishpiUserInfo), slog.Any("err", err))
			return err
		}
	}

	token, err := user.NewAuthToken()
	if err != nil {
		logger.Error("生成token失败", slog.Any("err", err))
		return err
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(user.Collection().AuthToken.DurationTime() / time.Second),
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	}
	event.SetCookie(cookie)
	return event.Redirect(http.StatusFound, "/?from=login")
}

func (controller *FishPiController) register(event *core.RequestEvent) error {
	//logger := controller.makeActionLogger("callback register").With(
	//	slog.String("path", event.Request.URL.String()),
	//)

	//openid := event.Get(ctxFishpiOpenId).(string)
	//fishpiUserInfo := event.Get(ctxFishpiUserInfo).(*FishpiUserInfo)
	//
	//userCollection, err := event.App.FindCollectionByNameOrId("users")
	//if err != nil {
	//	logger.Error("查找用户表失败", slog.Any("err", err))
	//	return err
	//}
	//var file *filesystem.File
	//if file, err = filesystem.NewFileFromURL(context.Background(), fishpiUserInfo.UserAvatarURL); err != nil {
	//	logger.Error("头像加载失败", slog.Any("err", err))
	//	return err
	//}
	//
	//user := core.NewRecord(userCollection)
	//user.SetEmail(fmt.Sprintf("%s@fishpi.cn", openid))
	//user.SetEmailVisibility(true)
	//user.SetVerified(true)
	//user.Set("name", fishpiUserInfo.Name())
	//user.Set("avatar", file)
	//user.Set("avatar_url", fishpiUserInfo.UserAvatarURL)
	//user.Set("role", "user")
	//user.SetRandomPassword()
	//
	//var authCollection *core.Collection
	//if authCollection, err = event.App.FindCollectionByNameOrId("externalAuths"); err != nil {
	//	logger.Error("查找用户表失败", slog.Any("err", err))
	//	return err
	//}
	//externalAuth := core.NewRecord(authCollection)
	//
	//externalAuth.Set("provider", "fishpi")
	//externalAuth.Set("openid", openid)
	//
	//var tokenCollection *core.Collection
	//if tokenCollection, err = event.App.FindCollectionByNameOrId("staticTokens"); err != nil {
	//	logger.Error("查找token表失败", slog.Any("err", err))
	//	return err
	//}
	//tokenRecord := core.NewRecord(tokenCollection)
	//
	//var walletCollection *core.Collection
	//if walletCollection, err = event.App.FindCollectionByNameOrId("wallets"); err != nil {
	//	logger.Error("查找钱包表失败", slog.Any("err", err))
	//	return err
	//}
	//wallet := core.NewRecord(walletCollection)
	//wallet.Set("level", 0)
	//wallet.Set("balance", 0)
	//
	//var token string
	//if err = event.App.RunInTransaction(func(txApp core.App) error {
	//	if e := txApp.Save(user); e != nil {
	//		return e
	//	}
	//	externalAuth.Set("user", user.Id)
	//	if e := txApp.Save(externalAuth); e != nil {
	//		return e
	//	}
	//	wallet.Set("user", user.Id)
	//	if e := txApp.Save(wallet); e != nil {
	//		return e
	//	}
	//	overdue := time.Hour * 24 * 30
	//	token, err = user.NewStaticAuthToken(overdue)
	//	if err != nil {
	//		logger.Error("创建token失败", slog.Any("err", err))
	//		return err
	//	}
	//	tokenRecord.Set("user", user.Id)
	//	tokenRecord.Set("token", token)
	//	tokenRecord.Set("expired", time.Now().Add(overdue))
	//	if e := txApp.Save(tokenRecord); e != nil {
	//		return e
	//	}
	//
	//	return nil
	//}); err != nil {
	//	bbb, _ := user.MarshalJSON()
	//	logger.Error("创建用户信息失败", slog.Any("user_record", user), slog.Any("user", string(bbb)), slog.Any("err", err))
	//	return err
	//}
	//
	//if clientId := event.Request.PathValue("client_id"); clientId != "" {
	//	if client, clientErr := event.App.SubscriptionsBroker().ClientById(clientId); clientErr == nil {
	//		// todo 错误处理
	//		if body, jsonErr := json.Marshal(map[string]any{
	//			"action": "register",
	//			"id":     user.Id,
	//			"token":  token,
	//			"is_new": true,
	//		}); jsonErr == nil {
	//			client.Send(subscriptions.Message{
	//				Name: "login",
	//				Data: body,
	//			})
	//		}
	//	}
	//}
	//
	//cookie := &http.Cookie{
	//	Name:     "token",
	//	Value:    tokenRecord.GetString("token"),
	//	Path:     "/",
	//	MaxAge:   int(tokenRecord.GetDateTime("expired").Sub(types.NowDateTime()).Seconds()),
	//	Secure:   true,
	//	HttpOnly: true,
	//	SameSite: http.SameSiteNoneMode,
	//}
	//event.SetCookie(cookie)
	//return event.JSON(http.StatusOK, map[string]any{
	//	"action": "register",
	//	"id":     user.Id,
	//	"token":  token,
	//	"is_new": true,
	//})
	return event.Redirect(http.StatusFound, "/?from=register")
}
