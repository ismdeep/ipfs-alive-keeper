package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/ismdeep/args"
	"github.com/ismdeep/ipfs-alive-keeper/config"
	"github.com/ismdeep/log"
	"net/http"
	"os"
	"sync"
)

func GetLinkType(url string) (string, error) {
	client := &http.Client{}
	r, err := client.Get(url)
	if err != nil {
		return "", err
	}
	return r.Header.Get("content-type"), nil
}

func GetLinks(url string) {
	log.Info("START LOAD URL", "url", url)
	contentType := ""
	for {
		var err error
		contentType, err = GetLinkType(url)
		if err != nil {
			//log.Error("URL LOAD ERROR [RETRY]", "err", err)
			continue
		}
		break
	}

	if contentType != "text/html" {
		//log.Info("STOP ON contentType != html/text", "url", url)
		for {
			client := http.Client{}
			_, err := client.Get(url)
			if err != nil {
				continue
			}
			break
		}
		return
	}

	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		log.Error("get links", "err", err)
		return
	}

	nodes := htmlquery.Find(doc, `//div[@class="table-responsive"]//table//tr`)
	for _, node := range nodes {
		tmpNode := htmlquery.FindOne(node, `.//td[2]//a`)
		if htmlquery.InnerText(tmpNode) == ".." {
			continue
		}

		GetLinks(fmt.Sprintf("https://ipfs.io%v", htmlquery.SelectAttr(tmpNode, "href")))
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(HelpMsg())
		return
	}

	if args.Exists("--help") {
		fmt.Println(HelpMsg())
		return
	}

	if !args.Exists("-c") {
		fmt.Println(HelpMsg())
		return
	}

	// 加载配置
	configPath := args.GetValue("-c")
	if err := config.Load(configPath); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, link := range config.DefaultConf.Links {
		wg.Add(1)
		go func(ipfsURL string) {
			for {
				GetLinks(ipfsURL)
			}
			wg.Done()
		}(link)
	}

	wg.Wait()
}
