package cf

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
	"os"
)

func GetCfAPI() *cloudflare.API {
	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		fmt.Println("无法创建 api 对象")
		panic(err)
	}
	return api
}

func GetCfAPI2(apiKey string, apiMail string) *cloudflare.API {
	api, err := cloudflare.New(apiKey, apiMail)
	if err != nil {
		fmt.Println("无法创建 api 对象")
		panic(err)
	}
	return api
}

func GetUserMsg(api *cloudflare.API) *cloudflare.User {
	user, err := api.UserDetails(context.Background())
	if err != nil {
		fmt.Println("无法获取 user msg")
		log.Fatal(err)
		return nil
	}
	return &user
}
