package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Ericwyn/cf-tools/cf"
	"github.com/cloudflare/cloudflare-go"
	"net"
	"os"
	"strings"
)

var apiKey = flag.String("apiKey", "", "cf 的 api Key")
var apiMail = flag.String("apiMail", "", "cf 的 apiMail")
var zoneId = flag.String("zoneId", "", "域名的 zoneId")

var dnsName = flag.String("name", "", "需要设置 dns 的域名, 可直接使用 ipv4/ipv6 替代本机 ip")
var dnsHost = flag.String("content", "", "解析到的 ip")
var dnsType = flag.String("type", "A", "dns 类型")
var dnsTtl = flag.Int("ttl", 120, "dns ttl")
var dnsProxied = flag.Bool("proxied", false, "使用 proxy(cf cdn 中转)")
var dnsProxiable = flag.Bool("proxiable", false, "使用 proxy(cf cdn 中转)")

var ipNow = flag.Bool("ip", false, "查看本机 ip 地址")
var ipv4Now = flag.Bool("ipv4", false, "查看本机 ip 地址")
var ipv6Now = flag.Bool("ipv6", false, "查看本机 ip 地址")

var allowSetLocalIp = flag.Bool("allowLocal", false, "允许设置为 C 类地址(192:xxx/fe80:xxx)")

var api *cloudflare.API

func main() {
	flag.Parse()

	if *ipNow {
		fmt.Println("ipv4:", getIp(false))
		fmt.Println("ipv6:", getIp(true))
		return
	}
	if *ipv4Now {
		fmt.Print(getIp(false))
		return
	}
	if *ipv6Now {
		fmt.Print(getIp(true))
		return
	}

	if *apiKey == "" || *apiMail == "" {
		fmt.Println("缺少 apiKey/apiMail 参数, 请使用 -h 查看说明")
		return
	}
	if *dnsName == "" {
		fmt.Println("请使用 -name 指定需要设置的域名")
	}
	if *dnsName == "" {
		fmt.Println("请使用 -content 指定 dns 的 ip")
	}
	if *zoneId == "" {
		fmt.Println("请使用 -zoneId 指定配置的域名")
	}

	// Construct a new API object
	api = cf.GetCfAPI2(*apiKey, *apiMail)

	// Fetch user details on the account
	user := cf.GetUserMsg(api)
	// Print user details
	//fmt.Println(user.Accounts)
	fmt.Println("当前cf账户邮箱:", user.Email)
	fmt.Println("===============================")
	fmt.Println("设置如下:")
	fmt.Println("name:\t", *dnsName)
	fmt.Println("host:\t", *dnsHost)
	fmt.Println("type:\t", *dnsType)
	fmt.Println("ttl:\t", *dnsTtl)
	fmt.Println("proxied:\t", *dnsProxied)
	fmt.Println("proxiable:\t", *dnsProxiable)
	fmt.Println("===============================")
	records, err := api.DNSRecords(context.Background(), *zoneId, cloudflare.DNSRecord{
		Type: *dnsType,
		Name: *dnsName,
	})
	if err != nil {
		fmt.Println("查询已有记录失败")
		//panic(err)
	} else {
		fmt.Println("\n")

		if *dnsHost == "ipv4" {
			*dnsHost = getIp(false)
			fmt.Println("查询得到本机 IPV4 地址:", *dnsHost)
		}
		if *dnsHost == "ipv6" {
			*dnsHost = getIp(true)
			fmt.Println("查询得到本机 IPV6 地址:", *dnsHost)
		}

		if *dnsName == "" {
			fmt.Println("获取本地 ip 失败")
			return
		}

		if isLocalAddress(*dnsHost) && !*allowSetLocalIp {
			fmt.Println("不允许设置为本地 ip 地址")
			return
		}

		newRec := cloudflare.DNSRecord{
			Type:    *dnsType,
			Name:    *dnsName,
			Content: *dnsHost,
			TTL:     *dnsTtl,
			Proxied: dnsProxied,
		}

		if len(records) != 1 {
			fmt.Println("未查询到已设置记录, 创建新记录")
			createNewRecords(newRec)
		} else {
			fmt.Println("查询到已有记录, 当前记录如下")
			fmt.Println("===============================")
			fmt.Println("name:\t", records[0].Name)
			fmt.Println("host:\t", records[0].Content)
			fmt.Println("type:\t", records[0].Type)
			fmt.Println("ttl:\t", records[0].TTL)
			fmt.Println("proxied:\t", *records[0].Proxied)
			fmt.Println("proxiable:\t", records[0].Proxiable)
			fmt.Println("===============================")
			updateRecords(records[0].ID, newRec)
		}
		//for i,record := range  records {
		//	fmt.Println("记录", i)
		//	fmt.Println(record.Content)
		//	fmt.Println(record.ZoneName)
		//	fmt.Println(record.ZoneID)
		//}
	}
}

func updateRecords(recordId string, record cloudflare.DNSRecord) {
	err := api.UpdateDNSRecord(context.Background(), *zoneId, recordId, record)
	if err == nil {
		fmt.Println("设置成功")
	} else {
		fmt.Println("设置失败")
		fmt.Println(err)
	}
}

func createNewRecords(record cloudflare.DNSRecord) {
	resp, err := api.CreateDNSRecord(context.Background(), *zoneId, record)
	if err != nil {
		fmt.Println("设置失败")
		fmt.Println(err)
		for _, info := range resp.Errors {
			fmt.Println(info.Code, ":", info.Message)
		}
	} else {
		fmt.Println("设置成功")
	}
}

func getIp(isIpv6 bool) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	resIp := ""
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ipParse := ipnet.IP.String()

			isIp := false
			if isIpv6 {
				isIp = IsIPv6(ipParse)
			} else {
				isIp = IsIPv4(ipParse)
			}

			if ipParse != "" && isIp {
				if resIp == "" {
					resIp = ipParse
				} else {
					// 如果当前 ip 已经存在，但是新的 ip 并不是本地地址的话，可以覆盖掉
					if !isLocalAddress(ipParse) {
						resIp = ipParse
					}
				}
			}
		}
	}
	return resIp
}

func getIpv6() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			//if ipnet.IP.To4() != nil {
			//	return ipnet.IP.String()
			//}
			ipParse := ipnet.IP.String()
			if ipParse != "" && IsIPv6(ipParse) {
				return ipParse
			}
		}
	}
	return ""
}

func IsIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

// 判断是否为本地地址
func isLocalAddress(address string) bool {
	return strings.Index(address, "fe80") == 0 ||
		strings.Index(address, "127") == 0 ||
		strings.Index(address, "192") == 0
}
