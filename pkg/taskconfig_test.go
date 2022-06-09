package pkg_test

import (
	"testing"
	"time"

	"github.com/Qianlitp/crawlergo/pkg"
	"github.com/Qianlitp/crawlergo/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestTaskConfigOptFunc(t *testing.T) {
	// 测试 https://github.com/Qianlitp/crawlergo/pull/101 修改的代码
	var taskConf pkg.TaskConfig
	for _, fn := range []pkg.TaskConfigOptFunc{
		pkg.WithTabRunTimeout(config.TabRunTimeout),
		pkg.WithMaxTabsCount(config.MaxTabsCount),
		pkg.WithMaxCrawlCount(config.MaxCrawlCount),
		pkg.WithDomContentLoadedTimeout(config.DomContentLoadedTimeout),
		pkg.WithEventTriggerInterval(config.EventTriggerInterval),
		pkg.WithBeforeExitDelay(config.BeforeExitDelay),
		pkg.WithEventTriggerMode(config.DefaultEventTriggerMode),
		pkg.WithIgnoreKeywords(config.DefaultIgnoreKeywords),
	} {
		fn(&taskConf)
	}

	// 应该都要等于默认配置
	assert.Equal(t, taskConf.TabRunTimeout, config.TabRunTimeout)
	assert.Equal(t, taskConf.MaxTabsCount, config.MaxTabsCount)
	assert.Equal(t, taskConf.MaxCrawlCount, config.MaxCrawlCount)
	assert.Equal(t, taskConf.DomContentLoadedTimeout, config.DomContentLoadedTimeout)
	assert.Equal(t, taskConf.EventTriggerInterval, config.EventTriggerInterval)
	assert.Equal(t, taskConf.BeforeExitDelay, config.BeforeExitDelay)
	assert.Equal(t, taskConf.EventTriggerMode, config.DefaultEventTriggerMode)
	assert.Equal(t, taskConf.IgnoreKeywords, config.DefaultIgnoreKeywords)

	// 重设超时时间
	taskConf.TabRunTimeout = time.Minute * 5

	// 企图覆盖自定义的时间, 不应该允许, 程序初始化时只能配置一次, 先由用户配置
	pkg.WithTabRunTimeout(time.Second * 5)(&taskConf)
	assert.NotEqual(t, taskConf.TabRunTimeout, time.Second*5)
	assert.NotEqual(t, taskConf.TabRunTimeout, config.TabRunTimeout)
}
