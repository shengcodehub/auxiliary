package i18n

import (
	"fmt"
	"github.com/spf13/cast"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	DefaultLanguage = language.English
	bundle          *i18n.Bundle
	localeMap       sync.Map
	mu              sync.Mutex
)

func Setup(localePath string) {
	bundle = i18n.NewBundle(DefaultLanguage)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	reader := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fmt.Println("load locale file: ", path)
		bundle.MustLoadMessageFile(path)
		return nil
	}

	err := filepath.Walk(localePath, reader)
	if err != nil {
		return
	}
}

func GetLocale(lang string) *i18n.Localizer {
	if locale, ok := localeMap.Load(lang); ok {
		return locale.(*i18n.Localizer)
	}

	mu.Lock()
	defer mu.Unlock()

	locale := i18n.NewLocalizer(bundle, lang)
	localeMap.Store(lang, locale)
	return locale
}

func SetLang(key string, lang string) {
	if lan, ok := localeMap.Load(key); ok {
		if cast.ToString(lan) == lang {
			return
		}
	}
	localeMap.Store(key, lang)
}

func GetLang(key string) string {
	if lang, ok := localeMap.Load(key); ok {
		return cast.ToString(lang)
	}
	return ""
}
