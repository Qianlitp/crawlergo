package filter

import (
	"testing"

	"github.com/Qianlitp/crawlergo/pkg/config"
	"github.com/Qianlitp/crawlergo/pkg/model"

	"github.com/stretchr/testify/assert"
)

var (
	// queryUrls = []string{
	// 	"http://test.nil.local.com/cctv/abcd?keyword=crawlergocrawlergo&end=1",
	// 	"http://test.nil.local.com/cctv/abcd?keyword=crawlergocrawlergo&end=1",
	// }

	fragmentUrls = []string{
		// 基准组
		"http://testhtml5.vuwm.com/latest#/page/1",
		"http://testhtml5.vuwm.com/latest#/page/search?keyword=Crawlergo&source=2&demo=1423&c=afa",
		// 被标记成 {{long}}
		"http://testhtml5.vuwm.com/latest#/page/search/fasdfsdafsdfsdfsdfasfsfasfafdsafssfasdfsd",

		// 对照组
		"http://testhtml5.vuwm.com/latest#/page/2",
		// 不应该被标记成 {{long}}
		"http://testhtml5.vuwm.com/latest#/page/search?keyword=CrawlergoCrawlergoCrawlergo&source=1&demo=1255&c=afa",
	}

	// completeUrls = []string{
	// 	"https://test.local.com:1234/adfatd/123456/sx14xi?user=crawlergo&pwd=fa1424&end=1#/user/info",
	// }
	smart = NewSmartFilter(NewSimpleFilter(""), true)
)

func TestDoFilter_countFragment(t *testing.T) {
	reqs := []model.Request{}
	for _, fu := range fragmentUrls {
		url, err := model.GetUrl(fu)
		assert.Nil(t, err)
		reqs = append(reqs, model.GetRequest(config.GET, url))
	}
	// #/page/1 和 #/page/2 是同一种类型
	assert.Equal(t, smart.calcFragmentID(reqs[0].URL.Fragment), smart.calcFragmentID(reqs[3].URL.Fragment))
	assert.Equal(t, smart.calcFragmentID(reqs[1].URL.Fragment), smart.calcFragmentID(reqs[4].URL.Fragment))
	for _, rq := range reqs[:2] {
		// 第一次出现都不应该过滤
		assert.Equal(t, smart.DoFilter(&rq), false)
	}
	for _, rq := range reqs[3:] {
		// 同类型出现第二次，应该被过滤
		assert.Equal(t, smart.DoFilter(&rq), true)
	}
}
