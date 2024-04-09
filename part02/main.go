package main

import (
	"fmt"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
	"log"
	"os"
)

const (

	// 替换成你的小程序 appiId 和 secretKey
	MINI_APP_LATEST_ID     string = "wxxxxf6"
	MINI_APP_LATEST_SECRET string = "7dfxxx64a"
)

func main() {
	// 初始化 Wechat 实例
	wc := wechat.NewWechat()
	//这里本地内存保存access_token，也可选择redis，memcache或者自定cache
	memory := cache.NewMemory()
	//token, _ := requestToken(MINI_APP_LATEST_ID, MINI_APP_LATEST_SECRET)
	cfg := &miniConfig.Config{
		AppID:     MINI_APP_LATEST_ID,
		AppSecret: MINI_APP_LATEST_SECRET,
		Cache:     memory,
	}

	//officialAccount := wc.GetOfficialAccount(cfg)
	mini := wc.GetMiniProgram(cfg)
	fmt.Println("------------")
	qcode := mini.GetQRCode()
	res, err := qcode.CreateWXAQRCode(qrcode.QRCoder{
		Page: "pages/home/index",
		Path: "?pathName=/race/pages/group&id=58",
		//CheckPath:  nil,
		Width: 300,
		//Scene: "pathName=/race/pages/group&id=58",
		//AutoColor:  false,
		//LineColor:  nil,
		//IsHyaline:  false,
		//EnvVersion: "",
	})
	if err != nil {
		return
	}
	// 生成图片
	png, err := os.Create("race.png")
	_, err2 := png.Write(res)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("QR code saved successfully!")

}
