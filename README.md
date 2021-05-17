# CloudFlareTools

对接 CloudFlare 的 go sdk， 提供命令行方式设定 DNS 域名解析记录

## 编译

```shell
 git clone https://github.com/Ericwyn/CloudFlareTools.git
 go build main/cftools.go
```

## 参数
```shell
Usage of cftools.exe:
  -allowLocal
        允许设置为 C 类地址(192:xxx/fe80:xxx)
  -apiKey string
        cf 的 api Key
  -apiMail string
        cf 的 apiMail
  -content string
        解析到的 ip
  -ip
        查看本机 ip 地址
  -ipv4
        查看本机 ip 地址
  -ipv6
        查看本机 ip 地址
  -name string
        需要设置 dns 的域名, 可直接使用 ipv4/ipv6 替代本机 ip
  -proxiable
        使用 proxy(cf cdn 中转)
  -proxied
        使用 proxy(cf cdn 中转)
  -ttl int
        dns ttl (default 120)
  -type string
        dns 类型 (default "A")
  -zoneId string
        域名的 zoneId

```

## 使用示例


- 将 www.domain.com 的 A 记录为地址 127.0.0.1

    ```shell
    cftools -apiMail
            $API_MAIL
            -apiKey
            $API_KEY
            -zoneId
            $ZONE_ID
            -content
            127.0.0.1
            -name
            www.domain.com
            -ttl
            100
            -type
            A
    ```
    - 需要从 Cloudflare 网站获取以下配置
        - apiMail 用户的 api mail (登录邮箱)
        - apiKey 用户的 api key (api token)
        - zoneId 域名绑定的 zone id
    
 - 将 www.domain.com 的 AAAA 记录为本机 ipv6 地址
    ```shell
    cftools -apiMail
            $API_MAIL
            -apiKey
            $API_KEY
            -zoneId
            $ZONE_ID
            -content
            ipv6
            -name
            www.domain.com
            -ttl
            100
            -type
            AAAA
    ```

 - 将 www.domain.com 的 AAAA 记录设置为本机 ipv6 地址并开启 cdn 转发

    ```shell
    cftools -apiMail
            $API_MAIL
            -apiKey
            $API_KEY
            -zoneId
            $ZONE_ID
            -content
            ipv6
            -name
            www.domain.com
            -ttl
            100
            -type
            AAAA
           -proxied 
           -proxiable
    ```
   

可配置 cron 定时任务，定时更新解析记录