package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var APIVersion = "10"

type DiscordWrapper struct {
	*discordgo.Session // 嵌入原 Session，继承所有方法
}

var (
	EndpointAPI                        = discordgo.EndpointDiscord + "api/v" + APIVersion + "/"
	EndpointApplications               = EndpointAPI + "applications/"
	EndpointEndpointApplicationsEmojis = func(aID string) string { return EndpointApplications + aID + "/emojis" }
	EndpointApplicationsEmojis         = func(aID string) string { return EndpointApplications + aID + "/emojis" }
)

func (s *DiscordWrapper) ApplicationsEmojiCreate(clientID string, data *discordgo.EmojiParams, options ...discordgo.RequestOption) (emoji *discordgo.Emoji, err error) {
	body, err := s.RequestWithBucketID("POST", EndpointEndpointApplicationsEmojis(clientID), data, EndpointApplicationsEmojis(clientID), options...)
	if err != nil {
		return
	}

	err = unmarshal(body, &emoji)
	return
}

func unmarshal(data []byte, v interface{}) error {
	err := discordgo.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("%w: %s", discordgo.ErrJSONUnmarshal, err)
	}
	return nil
}
