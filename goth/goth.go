package goth

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	"github.com/shengcodehub/auxiliary/goth/providers/instagram"
	"github.com/shengcodehub/auxiliary/goth/providers/riot"
	"net/http"
)

type Conf struct {
	Key       string
	Discord   DiscordConf
	Instagram InstagramConf
	Twitter   TwitterConf
	Riot      RiotConf
	Twitch    TwitchConf
	Facebook  FacebookConf
}

type DiscordConf struct {
	Key         string
	ClientKey   string
	Secret      string
	CallbackURL string
	Scopes      []string
}

type InstagramConf struct {
	ClientKey   string
	Secret      string
	CallbackURL string
	Scopes      []string
}

type TwitterConf struct {
	ClientKey   string
	Secret      string
	CallbackURL string
}

type RiotConf struct {
	ClientKey   string
	Secret      string
	CallbackURL string
	Scopes      []string
}

type TwitchConf struct {
	ClientKey   string
	Secret      string
	CallbackURL string
}

type FacebookConf struct {
	ClientKey   string
	Secret      string
	CallbackURL string
	Scopes      []string
}

func Setup(c Conf) {
	cookieStore := sessions.NewCookieStore([]byte(c.Key))
	cookieStore.Options.Secure = true
	cookieStore.Options.SameSite = http.SameSiteNoneMode
	gothic.Store = cookieStore

	goth.UseProviders(
		discord.New(c.Discord.ClientKey, c.Discord.Secret, c.Discord.CallbackURL, c.Discord.Scopes...),
		riot.New(c.Riot.ClientKey, c.Riot.Secret, c.Riot.CallbackURL, c.Riot.Scopes...),
		instagram.New(c.Instagram.ClientKey, c.Instagram.Secret, c.Instagram.CallbackURL, c.Instagram.Scopes...),
		twitter.New(c.Twitter.ClientKey, c.Twitter.Secret, c.Twitter.CallbackURL),
		twitch.New(c.Twitch.ClientKey, c.Twitch.Secret, c.Twitch.CallbackURL),
		facebook.New(c.Facebook.ClientKey, c.Facebook.Secret, c.Facebook.CallbackURL, c.Facebook.Scopes...),
	)
}
