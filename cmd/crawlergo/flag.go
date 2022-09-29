package main

import (
	"fmt"

	"github.com/Qianlitp/crawlergo/pkg/config"
	"github.com/urfave/cli/v2"
)

var cliFlags = []cli.Flag{
	SetIgnorePatterns(),
	SetCrawlTimeout(),
	SetChromePath(),
	SetCustomHeaders(),
	SetPostData(),
	SetMaxCrawledCount(),
	SetFilterMod(),
	SetOutputMode(),
	SetOutputJSON(),
	//SetIgcognitoContext(),
	SetMaxTabCount(),
	SetFuzzPath(),
	SetFuzzPathDict(),
	SetRobotsPath(),
	SetRequestProxy(),
	SetEncodeURL(),
	SetTabRunTTL(),
	SetWaitDomContentLoadedTTL(),
	SetEventTriggerMode(),
	SetEventTriggerInterval(),
	SetBeforeExitDelay(),
	SetIgnoreUrlKeywords(),
	SetFormValues(),
	SetFormKeywordValue(),
	SetPushToProxy(),
	SetPushPoolMax(),
	SetLogLevel(),
	SetNoHeadless(),
}

func SetChromePath() *cli.PathFlag {
	return &cli.PathFlag{
		Name:        "chromium-path",
		Aliases:     []string{"c"},
		Usage:       "`Path` of chromium executable. Such as \"/home/test/chrome-linux/chrome\"",
		Destination: &taskConfig.ChromiumPath,
		EnvVars:     []string{"CRAWLERGO_CHROMIUM_PATH"},
	}
}

func SetCustomHeaders() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "custom-headers",
		Usage:       "add additional `Headers` to each request. The input string will be called json.Unmarshal",
		Value:       fmt.Sprintf(`{"Spider-Name": "crawlergo", "User-Agent": "%s"}`, config.DefaultUA),
		Destination: &taskConfig.ExtraHeadersString,
	}
}

func SetPostData() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "post-data",
		Aliases:     []string{"d"},
		Usage:       "set `PostData` to target and use POST method.",
		Destination: &postData,
	}
}

func SetMaxCrawledCount() *cli.IntFlag {
	return &cli.IntFlag{
		Name:        "max-crawled-count",
		Aliases:     []string{"m"},
		Value:       config.MaxCrawlCount,
		Usage:       "the maximum `Number` of URLs visited by the crawler in this task.",
		Destination: &taskConfig.MaxCrawlCount,
	}
}

func SetFilterMod() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "filter-mode",
		Aliases:     []string{"f"},
		Value:       "smart",
		Usage:       "filtering `Mode` used for collected requests. Allowed mode:\"simple\", \"smart\" or \"strict\".",
		Destination: &taskConfig.FilterMode,
	}
}

func SetOutputMode() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "output-mode",
		Aliases:     []string{"o"},
		Value:       "console",
		Usage:       "console print or serialize output. Allowed mode:\"console\" ,\"json\" or \"none\".",
		Destination: &outputMode,
	}
}

func SetOutputJSON() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "output-json",
		Usage:       "write output to a json file.Such as result_www_test_com.json",
		Destination: &outputJsonPath,
	}
}

func SetIgnorePatterns() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:    "ignore-patterns",
		Aliases: []string{"ipm"},
		Usage:   "crawlergo will not crawl these URLs matched by regular expressions. e.g.: -irm https?://foo.bar.com/[0-9].html http://.*.xxx.com",
	}
}

func SetCrawlTimeout() *cli.DurationFlag {
	return &cli.DurationFlag{
		Name:    "crawl-timeout",
		Aliases: []string{"ct"},
		Usage:   "crawlergo will not crawl these URLs matched by regular expressions. e.g.: --ct 10m",
	}
}

func SetMaxTabCount() *cli.IntFlag {
	return &cli.IntFlag{
		Name:        "max-tab-count",
		Aliases:     []string{"t"},
		Value:       8,
		Usage:       "maximum `Number` of tabs allowed.",
		Destination: &taskConfig.MaxTabsCount,
	}
}

func SetFuzzPath() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "fuzz-path",
		Value:       false,
		Usage:       "whether to fuzz the target with common paths.",
		Destination: &taskConfig.PathByFuzz,
	}
}

func SetFuzzPathDict() *cli.PathFlag {
	return &cli.PathFlag{
		Name:        "fuzz-path-dict",
		Usage:       "`Path` of fuzz dict. Such as \"/home/test/fuzz_path.txt\"",
		Destination: &taskConfig.FuzzDictPath,
	}
}

func SetRobotsPath() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "robots-path",
		Value:       false,
		Usage:       "whether to resolve paths from /robots.txt.",
		Destination: &taskConfig.PathFromRobots,
	}
}

func SetRequestProxy() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "request-proxy",
		Usage:       "all requests connect through defined proxy server.",
		Destination: &taskConfig.Proxy,
	}
}

// return &cli.BoolFlag{
//	Name:        "bypass",
//	Value:       false,
//	Usage:       "whether to encode url with detected charset.",
//	Destination: &taskConfig.EncodeURLWithCharset,
//},
func SetEncodeURL() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "encode-url",
		Value:       false,
		Usage:       "whether to encode url with detected charset.",
		Destination: &taskConfig.EncodeURLWithCharset,
	}
}

func SetTabRunTTL() *cli.DurationFlag {

	return &cli.DurationFlag{
		Name:        "tab-run-timeout",
		Value:       config.TabRunTimeout,
		Usage:       "the `Timeout` of a single tab task.",
		Destination: &taskConfig.TabRunTimeout,
	}
}

func SetWaitDomContentLoadedTTL() *cli.DurationFlag {
	return &cli.DurationFlag{
		Name:        "wait-dom-content-loaded-timeout",
		Value:       config.DomContentLoadedTimeout,
		Usage:       "the `Timeout` of waiting for a page dom ready.",
		Destination: &taskConfig.DomContentLoadedTimeout,
	}
}

func SetEventTriggerMode() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "event-trigger-mode",
		Value:       config.EventTriggerAsync,
		Usage:       "this `Value` determines how the crawler automatically triggers events.Allowed mode:\"async\" or \"sync\".",
		Destination: &taskConfig.EventTriggerMode,
	}
}

func SetEventTriggerInterval() *cli.DurationFlag {
	return &cli.DurationFlag{
		Name:        "event-trigger-interval",
		Value:       config.EventTriggerInterval,
		Usage:       "the `Interval` of triggering each event.",
		Destination: &taskConfig.EventTriggerInterval,
	}
}

func SetBeforeExitDelay() *cli.DurationFlag {
	return &cli.DurationFlag{
		Name:        "before-exit-delay",
		Value:       config.BeforeExitDelay,
		Usage:       "the `Time` of waiting before crawler exit.",
		Destination: &taskConfig.BeforeExitDelay,
	}
}

func SetIgnoreUrlKeywords() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:        "ignore-url-keywords",
		Aliases:     []string{"iuk"},
		Value:       ignoreKeywords,
		Usage:       "crawlergo will not crawl these URLs matched by `Keywords`. e.g.: -iuk logout -iuk quit -iuk exit",
		DefaultText: "Default [logout quit exit]",
	}
}

func SetFormValues() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:    "form-values",
		Aliases: []string{"fv"},
		Value:   customFormTypeValues,
		Usage:   "custom filling text for each form type. e.g.: -fv username=crawlergo_nice -fv password=admin123",
	}
}

// 根据关键词自行选择填充文本
func SetFormKeywordValue() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:    "form-keyword-values",
		Aliases: []string{"fkv"},
		Value:   customFormKeywordValues,
		Usage:   "custom filling text, fuzzy matched by keyword. e.g.: -fkv user=crawlergo_nice -fkv pass=admin123",
	}
}

func SetPushToProxy() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "push-to-proxy",
		Usage:       "every request in 'req_list' will be pushed to the proxy `Address`. Such as \"http://127.0.0.1:8080/\"",
		Destination: &pushAddress,
	}
}

func SetPushPoolMax() *cli.IntFlag {
	return &cli.IntFlag{
		Name:        "push-pool-max",
		Usage:       "maximum `Number` of concurrency when pushing results to proxy.",
		Value:       DefaultMaxPushProxyPoolMax,
		Destination: &pushProxyPoolMax,
	}
}

func SetLogLevel() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "log-level",
		Usage:       "log print `Level`, options include debug, info, warn, error and fatal.",
		Value:       DefaultLogLevel,
		Destination: &logLevel,
	}
}

func SetNoHeadless() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "no-headless",
		Value:       false,
		Usage:       "no headless mode",
		Destination: &taskConfig.NoHeadless,
	}
}
