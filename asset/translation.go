package asset

import(
	"github.com/nicksnyder/go-i18n/i18n"
)

var(
	Tr i18n.TranslateFunc
)

// go-bindata -pkg asset .
func init() {
	// TO DO : find another way maybe with assets, go bind-data
	out1, _ := Asset("en-us.all.yaml")
	out2, _ := Asset("fr-be.all.yaml")
	i18n.ParseTranslationFileBytes("en-us.all.yaml",out1)
	i18n.ParseTranslationFileBytes("fr-be.all.yaml",out2)
	Tr, _ = i18n.Tfunc("fr-be")
}