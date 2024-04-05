package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"
	"github.com/valyala/fasthttp"
)

var idx = 0

func Forward(redirectURL []string) fiber.Handler {
	return func(c fiber.Ctx) error {
		parth := c.Params("*")

		rurl, _ := url.Parse(redirectURL[idx%len(redirectURL)])
		idx += 1
		target := &url.URL{
			Scheme: rurl.Scheme,
			Host:   rurl.Host,
			Path:   parth, //
		}
		fmt.Println(target)
		copyReq := new(fasthttp.Request)
		c.Request().CopyTo(copyReq)
		fmt.Println(parth)
		fmt.Println(redirectURL)
		fmt.Println(string(copyReq.Body()))
		return proxy.Do(c, target.String())
	}
}

func main() {
	// 解析命令行参数，第一个参数是 port，后面的不确定个数的是 redirect_urls数组，获取 port 和 redirect_urls数组
	args := os.Args[1:] // 去除第一个元素（程序名）

	if len(args) < 2 {
		fmt.Println("Usage: program_name port redirect_url [redirect_url...]")
		os.Exit(1)
	}

	portStr := args[0]
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		fmt.Printf("Invalid port value '%s'\n", portStr)
		os.Exit(1)
	}

	redirects := args[1:]
	for _, rawURL := range redirects {
		_, err := url.Parse(rawURL)
		if err != nil {
			fmt.Printf("Invalid redirect URL '%s': %v\n", rawURL, err)
			os.Exit(1)
		}
	}
	proxy.ConfigDefault.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	app := fiber.New()

	app.Post("/*", Forward(redirects))
	fmt.Printf("serv on: %d \n\n", port)
	log.Fatal(app.Listen(fmt.Sprintf("0.0.0.0:%d", port)))
}
