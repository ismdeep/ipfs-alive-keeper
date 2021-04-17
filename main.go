package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/ismdeep/ismdeep-go-utils/args_util"
	"github.com/ismdeep/log"
	"net/http"
	"os"
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
			log.Error("URL LOAD ERROR [RETRY]", "err", err)
			continue
		}
		break
	}

	if contentType != "text/html" {
		log.Info("STOP ON contentType != html/text", "url", url)
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
		//log.Info("node info", "tmpNode", tmpNode.Data, "text", htmlquery.InnerText(tmpNode))
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

	if args_util.Exists("--help") {
		fmt.Println(HelpMsg())
		return
	}

	ipfsUrl := os.Args[1]
	GetLinks(ipfsUrl)
}
