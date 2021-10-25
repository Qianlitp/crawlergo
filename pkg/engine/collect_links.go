package engine

import (
	"crawlergo/pkg/config"
	"crawlergo/pkg/logger"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"regexp"
)

/**
最后收集所有的链接
*/
func (tab *Tab) collectLinks() {
	go tab.collectHrefLinks()
	go tab.collectObjectLinks()
	go tab.collectCommentLinks()
}

func (tab *Tab) collectHrefLinks() {
	defer tab.collectLinkWG.Done()
	// 收集 src href data-url 属性值
	attrNameList := []string{"src", "href", "data-url", "data-href"}
	for _, attrName := range attrNameList {
		var attrs []map[string]string
		_ = chromedp.Run(*tab.Ctx, chromedp.AttributesAll(fmt.Sprintf(`[%s]`, attrName), &attrs, chromedp.ByQueryAll, chromedp.AtLeast(0)))
		for _, attrMap := range attrs {
			tab.AddResultUrl(config.GET, attrMap[attrName], config.FromDOM)
		}
	}
}

func (tab *Tab) collectObjectLinks() {
	defer tab.collectLinkWG.Done()
	// 收集 object[data] links
	var attrs []map[string]string
	if err := chromedp.Run(*tab.Ctx, chromedp.AttributesAll(`object[data]`, &attrs, chromedp.ByQueryAll, chromedp.AtLeast(0))); err != nil {
		logger.Logger.Debug(err)
		return
	}
	for _, attrMap := range attrs {
		tab.AddResultUrl(config.GET, attrMap["data"], config.FromDOM)
	}
}

func (tab *Tab) collectCommentLinks() {
	defer tab.collectLinkWG.Done()
	// 收集注释中的链接
	var nodes []*cdp.Node
	if err := chromedp.Run(*tab.Ctx, chromedp.Nodes(`//comment()`, &nodes, chromedp.BySearch)); err != nil {
		logger.Logger.Debug("get comment nodes err")
		logger.Logger.Debug(err)
		return
	}
	urlRegex := regexp.MustCompile(config.URLRegex)
	for _, node := range nodes {
		content := node.NodeValue
		urlList := urlRegex.FindAllString(content, -1)
		for _, url := range urlList {
			tab.AddResultUrl(config.GET, url, config.FromComment)
		}
	}
}
