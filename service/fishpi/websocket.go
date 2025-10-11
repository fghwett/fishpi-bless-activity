package fishpi

import (
	"log/slog"
	"stick/pkg/utils"
	"time"

	"github.com/duke-git/lancet/v2/retry"
	"github.com/lxzan/gws"
)

func (service *Service) start() {
	if err := service.reconnect(); err != nil {
		service.logger.Error("【聊天室】首次连接失败", slog.Any("err", err))
		return
	}
}

func (service *Service) connect() error {

	res, err := service.GetChatroomNodeGet()
	if err != nil {
		return err
	}

	eh := newEventHandler(service.app, service.onClose)
	socket := new(gws.Conn)
	if socket, _, err = gws.NewClient(eh, &gws.ClientOption{
		Addr: res.Data,
		PermessageDeflate: gws.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
	}); err != nil {
		return err
	}
	service.conn = socket

	go socket.ReadLoop()

	return nil
}

func (service *Service) onClose() {
	retryTimes := 50
	if err := retry.Retry(
		service.reconnect,
		retry.RetryTimes(uint(retryTimes)),
		retry.RetryWithCustomBackoff(&utils.RetryStrategy{
			Base:       3 * time.Second,
			Interval:   time.Second,
			Multiplier: 2,
			Max:        10 * time.Minute,
		}),
	); err != nil {
		service.logger.Error("【聊天室】多次重连失败", slog.Int("times", retryTimes), slog.Any("err", err))

		go service.onClose()
	}
}

func (service *Service) reconnect() error {
	if err := service.connect(); err != nil {
		service.logger.Error("【聊天室】重连失败", slog.Any("err", err))
		return err
	}
	return nil
}
