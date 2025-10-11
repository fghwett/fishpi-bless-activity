package fishpi

import (
	"context"
	"log/slog"
	"slices"
	"stick/model"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/lxzan/gws"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/tidwall/gjson"
)

type eventHandler struct {
	app     core.App
	onClose func()

	collection *core.Collection

	logger *slog.Logger
}

func newEventHandler(app core.App, onClose func()) *eventHandler {
	collection, _ := app.FindCollectionByNameOrId(`chatroom_messages`)

	return &eventHandler{
		app:     app,
		onClose: onClose,

		collection: collection,

		logger: app.Logger().WithGroup("chatroom_event_handler"),
	}
}

func (eh *eventHandler) OnOpen(socket *gws.Conn) {
	eh.logger.Debug("【聊天室】连接成功")

	heartbeat := NewHeartbeat(eh.app, socket)
	socket.Session().Store("heartbeat", heartbeat)

	go heartbeat.Start()
}

func (eh *eventHandler) OnClose(socket *gws.Conn, err error) {
	eh.logger.Debug("【聊天室】连接关闭", slog.Any("err", err))

	if eh.onClose == nil {
		return
	}
	// 如果有心跳，则停止心跳
	if value, exist := socket.Session().Load("heartbeat"); exist {
		if heartbeat, ok := value.(*Heartbeat); ok {
			heartbeat.Stop()
			socket.Session().Delete("heartbeat")
		}
	}

	eh.onClose()
}

func (eh *eventHandler) OnPing(socket *gws.Conn, payload []byte) {}

func (eh *eventHandler) OnPong(socket *gws.Conn, payload []byte) {}

func (eh *eventHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer func() {
		_ = message.Close()
	}()
	if message.Opcode != gws.OpcodeText {
		return
	}
	//eh.logger.Debug("【聊天室】收到消息", slog.String("message", string(message.Bytes())))

	result := gjson.ParseBytes(message.Bytes())
	oId := result.Get("oId").String()
	msgType := result.Get("type").String()
	userOId := result.Get("userOId").String()
	userName := result.Get("userName").String()
	userNickname := result.Get("userNickname").String()
	md := result.Get("md").String()
	client := result.Get("client").String()
	//created := result.Get("time").String()
	t, _ := time.ParseInLocation(time.DateTime, result.Get("time").String(), time.Local)
	created, _ := types.ParseDateTime(t)

	result.ForEach(func(key, value gjson.Result) bool {
		if slices.Contains([]string{
			"oId", "type", "userOId", "userName", "userNickname", "md", "client", "time",
		}, key.String()) {
			return false
		}
		return true
	})

	if oId == "" {
		time.Now().Unix()
		oId = convertor.ToString(time.Now().UnixMilli())
	}
	if created.IsZero() {
		created = types.NowDateTime()
	}

	chatroomMessage := model.NewChatroomMessage(core.NewRecord(eh.collection))
	chatroomMessage.SetOId(oId)

	chatroomMessage.SetType(msgType)
	chatroomMessage.SetUserOId(userOId)
	chatroomMessage.SetUserName(userName)
	chatroomMessage.SetUserNickname(userNickname)
	chatroomMessage.SetMd(md)
	chatroomMessage.SetClient(client)
	chatroomMessage.SetMessage(result.String())
	chatroomMessage.SetCreated(created)

	if err := eh.app.Save(chatroomMessage); err != nil {
		eh.logger.Error("【聊天室】保存消息失败", slog.Any("err", err), slog.String("message", string(message.Bytes())), slog.Any("record", chatroomMessage))
		return
	}
}

type Heartbeat struct {
	app             core.App
	conn            *gws.Conn
	onlineHookId    string
	terminateHookId string
	ctx             context.Context
	cancel          context.CancelFunc

	lastTime time.Time
}

func NewHeartbeat(app core.App, conn *gws.Conn) *Heartbeat {
	heartbeat := &Heartbeat{
		app:      app,
		conn:     conn,
		lastTime: time.Now(),
	}

	heartbeat.ctx, heartbeat.cancel = context.WithCancel(context.Background())

	return heartbeat
}

func (heartbeat *Heartbeat) Start() {

	heartbeat.onlineHookId = heartbeat.app.OnRecordAfterCreateSuccess("chatroom_messages").BindFunc(heartbeat.onMessage)
	heartbeat.app.OnTerminate().BindFunc(heartbeat.onTerminate)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-heartbeat.ctx.Done():
			heartbeat.app.Logger().Warn("【聊天室】手动断开链接")
			heartbeat.app.OnRecordAfterCreateSuccess("chatroom_messages").Unbind(heartbeat.onlineHookId)
			heartbeat.app.OnTerminate().Unbind(heartbeat.terminateHookId)
			_ = heartbeat.conn.WriteClose(1000, []byte("disconnect"))
			return
		case <-ticker.C:
			if time.Since(heartbeat.lastTime) > 120*time.Second {
				heartbeat.app.Logger().Warn("【聊天室】心跳超时 准备重连")
				heartbeat.app.OnRecordAfterCreateSuccess("chatroom_messages").Unbind(heartbeat.onlineHookId)
				heartbeat.app.OnTerminate().Unbind(heartbeat.terminateHookId)
				_ = heartbeat.conn.WriteClose(1000, []byte("heartbeat timeout"))
				return
			}
		}
	}
}

func (heartbeat *Heartbeat) Stop() {
	heartbeat.cancel()
}

func (heartbeat *Heartbeat) onMessage(event *core.RecordEvent) error {
	heartbeat.lastTime = time.Now()

	return event.Next()
}

func (heartbeat *Heartbeat) onTerminate(event *core.TerminateEvent) error {
	heartbeat.Stop()

	return event.Next()
}
