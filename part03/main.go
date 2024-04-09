package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
	"log"
	"os"
)

const (
	MINI_APP_LATEST_ID     string = "wx7d3da0c61f5198f6"
	MINI_APP_LATEST_SECRET string = "7df2a2f3c8e8446f7842a3445b77464a"
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

	png, err := os.Create("data.png")
	_, err2 := png.Write(res)

	if err2 != nil {
		log.Fatal(err2)
	}

	// 创建s3协议实例（这里不是指定aws的s3，而是所有s3协议的）
	sess := session.Must(session.NewSession(
		&aws.Config{
			// replaced
			Region:      aws.String("ap-beijing"),
			Credentials: credentials.NewStaticCredentials("your-AccessKeyID", "your-SecretAccessKey", ""),
			Endpoint:    aws.String("cos.ap-beijing.myqcloud.com"),
		},
	))
	CloudStorage := s3.New(sess)

	file, err := os.Open("race.png")
	if err != nil {
		fmt.Println("os.Ope:", err)
	}
	defer file.Close()

	// 将内存中的文件内容转换为 io.Reader
	_, err = CloudStorage.PutObjectWithContext(context.Background(), &s3.PutObjectInput{
		// replaced
		Bucket: aws.String("your-bucket"),
		// 这里是从 桶-id后开始写入，如下：
		// https://cos.ctra.top/ctra/xxx/xxx/xxx.jpeg
		// 我们需要传入 从 /ctra 开始的路径地址
		Key:  aws.String("/ctra/xxx/xxx/xxx.jpeg"),
		Body: file,
	})
	if err != nil {
		fmt.Println("failed to upload object！", err)
	}

	fmt.Println("QR code saved successfully!")

}
