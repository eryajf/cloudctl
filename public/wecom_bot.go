package public

import (
	botApi "github.com/electricbubble/wecom-bot-api"
)

func SendMsg(botkey, msg string) error {
	return botApi.NewWeComBot(botkey).PushMarkdownMessage(msg)
}
