package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	CreateWXAQRCodeURL   = "https://api.weixin.qq.com/cgi-bin/wxaapp/createwxaqrcode?access_token=%s"
	GetWXACodeURL        = "https://api.weixin.qq.com/wxa/getwxacode?access_token=%s"
	GetWXACodeUnlimitURL = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
)

const (
	// 替换成你的小程序 appiId 和 secretKey
	MINI_APP_LATEST_ID     string = "wxxxxf6"
	MINI_APP_LATEST_SECRET string = "7dfxxx64a"
)

// 通过 普通的源生请求
func main() {
	// 要生成二维码的路径和参数
	_, err := GetQRCode(58)
	if err != nil {
		fmt.Println("Error getting QR code:", err)
		return
	}

}

func GetQRCode(id int) ([]byte, error) {
	//上面生成的access code 判断为空时重新请求
	accessToken, err := requestToken(MINI_APP_LATEST_ID, MINI_APP_LATEST_SECRET)
	if err != nil {
		return nil, err
	}
	strUrl := fmt.Sprintf(CreateWXAQRCodeURL, accessToken)

	// todo: 这里有bug待完善
	parm := make(map[string]string)
	// path 要以?开始
	parm["path"] = fmt.Sprintf("?pathName=/race/pages/group&id=%d", id)
	// page 起始位置不需要 /
	parm["page"] = fmt.Sprintf("pages/home/index")
	jsonStr, err := json.Marshal(parm)
	if err != nil {
		return nil, errors.New("json Marshal QRCode paramter err :" + err.Error())
	}
	req, err := http.NewRequest("POST", strUrl, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return nil, errors.New("get QRCode err :" + err.Error())
	}
	// 发起 post 请求
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("get QRCode err :" + err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("get QRCode err :" + err.Error())
	}

	// 生成图片
	png, err := os.Create("data.png")
	_, err2 := png.Write(body)

	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("QR code saved successfully!")
	return nil, nil
}

func requestToken(appid, secret string) (string, error) {
	u, err := url.Parse("https://api.weixin.qq.com/cgi-bin/token")
	if err != nil {
		log.Fatal(err)
	}
	paras := &url.Values{}
	//设置请求参数
	paras.Set("appid", appid)
	paras.Set("secret", secret)
	paras.Set("grant_type", "client_credential")
	u.RawQuery = paras.Encode()
	resp, err := http.Get(u.String())
	//关闭资源
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", errors.New("request token err :" + err.Error())
	}

	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return "", errors.New("request token response json parse err :" + err.Error())
	}
	if jMap["errcode"] == nil || jMap["errcode"] == 0 {
		accessToken, _ := jMap["access_token"].(string)
		return accessToken, nil
	} else {
		//返回错误信息
		errcode := jMap["errcode"].(string)
		errmsg := jMap["errmsg"].(string)
		err = errors.New(errcode + ":" + errmsg)
		return "", err
	}
}
