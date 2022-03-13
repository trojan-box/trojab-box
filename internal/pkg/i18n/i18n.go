package i18n_lib

import (
	"github.com/BurntSushi/toml"
	"github.com/aresprotocols/trojan-box/internal/pkg/constant"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"path/filepath"
)

func GetLocalizer(lang string) *i18n.Localizer {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	availableLangs := make(map[string]string)
	availableLangs["en"] = constant.LangEn
	availableLangs["zh_CN"] = constant.LangZhCn

	if _, ok := availableLangs[lang]; !ok {
		lang = constant.LangEn
	}
	currDir, _ := filepath.Abs("./config/i18n")
	localLangFile := currDir + "/broadcast." + lang + ".toml"
	bundle.MustLoadMessageFile(localLangFile)
	return i18n.NewLocalizer(bundle, lang)
}
