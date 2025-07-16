package task_chain

import (
	"github.com/Yostardev/gf"
	"github.com/tiancheng92/task-chain/internal/lark"
	"github.com/tiancheng92/task-chain/internal/llm"
	"github.com/tiancheng92/task-chain/internal/modal"
	"gorm.io/gorm"
)

type initSetting struct {
	mysqlDsn        string
	db              *gorm.DB
	larkAppID       string
	larkAppSecret   string
	superLinkUrlFmt string
	llmBaseUrl      string
	llmApiKey       string
	llmModel        string
	plugins         []string
	customPlugins   []CustomPlugin
}

type initOption func(setting *initSetting)

func SetDBByDsn(dsn string) initOption {
	return func(setting *initSetting) {
		setting.mysqlDsn = dsn

		if !(gf.ArrayContains(setting.plugins, "use_db")) {
			setting.plugins = append(setting.plugins, "use_db")
		}
	}
}

func SetDB(db *gorm.DB) initOption {
	return func(setting *initSetting) {
		setting.db = db

		if !(gf.ArrayContains(setting.plugins, "use_db")) {
			setting.plugins = append(setting.plugins, "use_db")
		}
	}
}

func SetLarkAppInfo(appID, appSecret, superLinkUrlFmt string) initOption {
	return func(setting *initSetting) {
		setting.larkAppID = appID
		setting.larkAppSecret = appSecret
		setting.superLinkUrlFmt = superLinkUrlFmt

		if !(gf.ArrayContains(setting.plugins, "send_lark_msg")) {
			setting.plugins = append(setting.plugins, "send_lark_msg")
		}
	}
}

func SetLLMInfo(baseUrl, apiKey, modelName string) initOption {
	return func(setting *initSetting) {
		setting.llmBaseUrl = baseUrl
		setting.llmApiKey = apiKey
		setting.llmModel = modelName

		if !(gf.ArrayContains(setting.plugins, "use_llm")) {
			setting.plugins = append(setting.plugins, "use_llm")
		}
	}
}

func SetCustomPlugins(cpl ...CustomPlugin) initOption {
	return func(setting *initSetting) {
		setting.customPlugins = cpl
	}
}

func Init(opts ...initOption) {
	var setting initSetting
	for i := range opts {
		opts[i](&setting)
	}

	if setting.mysqlDsn != "" {
		modal.InitByDSN(setting.mysqlDsn)
	}

	if setting.db != nil {
		modal.InitByDB(setting.db)
	}

	lark.SetAppInfo(setting.larkAppID, setting.larkAppSecret, setting.superLinkUrlFmt)
	llm.SetLLMInfo(setting.llmBaseUrl, setting.llmApiKey, setting.llmModel)
	plugins = setting.plugins
	customPluginList = setting.customPlugins
}
