package intl

import "sort"

var languageMap = map[string]string{
	// Arabic
	"ar": "العربية",
	// Azerbaijani
	"az": "Azərbaycanca",
	// English
	"en": "English",
	// Spanish
	"es": "Español",
	// French
	"fr": "Français",
	// Hungarian
	"hu": "magyar",
	// Italian
	"it": "Italiano",
	// Japanese
	"ja": "日本語",
	// Korean
	"ko": "한국어",
	// Mongolian
	"mn": "Монгол хэл",
	// Polish
	"pl": "Polski",
	// Portuguese (Brazil)
	"pt-BR": "Português do Brasil",
	// Russian
	"ru": "Русский",
	// Turkish
	"tr": "Türkçe",
	// Ukrainian
	"uk": "Українська",
	// Simplified Chinese
	"zh-Hans": "简体中文",
	// Traditional Chinese
	"zh-Hant": "繁體中文",
}

type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func LanguageMap() map[string]string {
	return languageMap
}

func LanguageList() []Language {
	var keys []string
	for k := range languageMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var list []Language
	for _, k := range keys {
		list = append(list, Language{
			Code: k,
			Name: languageMap[k],
		})
	}
	return list
}
