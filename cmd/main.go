package main

import (
	"app/internal/appconst"
	"app/internal/engine"
	"app/internal/mq"
	"fmt"
	"gitlab.jiandan100.cn/webdev/gocomm/pkg"
	"gitlab.jiandan100.cn/webdev/gocomm/pkg/appcore"
	_ "gitlab.jiandan100.cn/webdev/gocomm/pkg/appcore"
	"gitlab.jiandan100.cn/webdev/gocomm/pkg/cache"
	"gitlab.jiandan100.cn/webdev/gocomm/pkg/db"
	"gitlab.jiandan100.cn/webdev/gocomm/pkg/mqueue"
	"log"
	"path/filepath"
)

func main() {
	eng := engine.NewEngine()
	eng.RegisterHooks(appcore.StageAfterStop, func() error {
		db.Close()
		mqueue.Close()
		fmt.Println("exit app ...")
		return nil
	})
	if appconst.UseDB {
		loadDB(eng.ConfPath, eng)
	}

	if appconst.UseCache {
		loadCache(eng.ConfPath, eng)
	}

	if appconst.UseMqProducer || appconst.UseMqConsumer {
		loadMq(eng.ConfPath, eng)
	}
	if appconst.UseMqConsumer {
		loadMqConsumer()
	}
	loadApp()
	if err := eng.Run(); err != nil {
		log.Fatal(err)
	}
}

func loadDB(confPath string, eng *engine.Engine) {
	if ok := db.Init(filepath.Join(confPath, "database.yml")); ok != nil {
		eng.InitOK = false
	}
}

func loadCache(confPath string, eng *engine.Engine) {
	if ok := cache.Init(filepath.Join(confPath, "cache-sample.yml")); ok != nil {
		eng.InitOK = false
	}
}

func loadMq(confPath string, eng *engine.Engine) {
	if ok := mqueue.Init(filepath.Join(confPath, "mq-sample.yml")); ok != nil {
		eng.InitOK = false
	}
}

func loadMqConsumer() {
	mqconsumer.LoadConsumerHandler()
}

func loadApp() {
	pkg.SetAppName(appconst.AppName)
}
