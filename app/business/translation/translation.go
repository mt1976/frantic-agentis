package translation

import (
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/logHandler"
	trnsl8r "github.com/mt1976/trnsl8r_connect"
)

var trnsServerHost string
var trnsServerPort int
var trnsServerProtocol string
var trnsLocale string
var cfg *commonConfig.Settings
var appName string

func init() {
	logHandler.TranslationLogger.Printf("[%] Initialised", "frantic-agentis")
	cfg = commonConfig.Get()
	trnsServerProtocol = cfg.GetTranslationServer_Protocol()
	trnsServerHost = cfg.GetTranslationServer_Host()
	trnsServerPort = cfg.GetTranslationServer_Port()
	trnsLocale = cfg.GetApplication_Locale()
	appName = cfg.GetApplication_Name()
}

func Get(in string) string {
	// Validate the input data

	tl8 := trnsl8r.NewRequest().WithProtocol(trnsServerProtocol).WithHost(trnsServerHost).WithPort(trnsServerPort).WithLogger(logHandler.TranslationLogger).FromOrigin(appName)

	resp, err := tl8.Get(in)
	if err != nil {
		logHandler.ErrorLogger.Println(err.Error())
		return in
	}
	if resp.Information != "" {
		logHandler.WarningLogger.Println(resp.Information)
	}
	if resp.Translated == "" {
		logHandler.ErrorLogger.Println("no response from translation service")
		return in
	}
	return resp.Translated
}
