package discord

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
)

var dgClient *discordgo.Session

func SetUp(token string) {
	// 创建一个新的 Discord 会话
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logx.Errorf("无法创建 Discord 会话: %v", err)
	}

	dgClient = dg
}

func GetDiscord() *discordgo.Session {
	return dgClient
}

func HandleDiscordInteractions(r *http.Request) *discordgo.Interaction {
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logx.Errorf("HandleDiscordInteractions 读取请求体失败: %v", err)
		return nil
	}
	defer r.Body.Close()
	interaction, err := UnmarshalJSON(body)
	if err != nil {
		logx.Errorf("HandleDiscordInteractions 解析请求体失败: %v, body:%s", err, string(body))
		return nil
	}
	if interaction.Member != nil {
		interaction.User = interaction.Member.User
	}
	return interaction
}

func UnmarshalJSON(data []byte) (i *discordgo.Interaction, err error) {
	if err = json.Unmarshal(data, &i); err != nil {
		logx.Errorf("Interaction.UnmarshalJSON,failed to unmarshal interaction, i error: %s", err.Error())
		return nil, err
	}
	switch i.Type {
	case discordgo.InteractionApplicationCommand, discordgo.InteractionApplicationCommandAutocomplete:
		i.ApplicationCommandData()
	case discordgo.InteractionMessageComponent:
		i.MessageComponentData()
	case discordgo.InteractionModalSubmit:
		i.ModalSubmitData()
	}
	return i, nil
}
