package tr

import(
	"github.com/nicksnyder/go-i18n/i18n"
)

var(
	Tr i18n.TranslateFunc
)

func init() {
	// TO DO : find another way maybe with assets, go bind-data
	i18n.MustLoadTranslationFile("tr/en-us.all.yaml")
	i18n.MustLoadTranslationFile("tr/fr-be.all.yaml")
	Tr, _ = i18n.Tfunc("fr-be")
}
