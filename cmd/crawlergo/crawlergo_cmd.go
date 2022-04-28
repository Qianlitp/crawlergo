package main

import (
	"crawlergo/pkg"
	"crawlergo/pkg/config"
	"crawlergo/pkg/logger"
	model2 "crawlergo/pkg/model"
	"crawlergo/pkg/tools"
	"crawlergo/pkg/tools/requests"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
)

/**
命令行调用适配器

用于生成开源的二进制程序
*/

type Result struct {
	ReqList       []Request `json:"req_list"`
	AllReqList    []Request `json:"all_req_list"`
	AllDomainList []string  `json:"all_domain_list"`
	SubDomainList []string  `json:"sub_domain_list"`
}

type Request struct {
	Url     string                 `json:"url"`
	Method  string                 `json:"method"`
	Headers map[string]interface{} `json:"headers"`
	Data    string                 `json:"data"`
	Source  string                 `json:"source"`
}

type ProxyTask struct {
	req       *model2.Request
	pushProxy string
}

const DefaultMaxPushProxyPoolMax = 10
const DefaultLogLevel = "Info"

var taskConfig pkg.TaskConfig
var outputMode string
var postData string
var signalChan chan os.Signal
var ignoreKeywords *cli.StringSlice
var customFormTypeValues *cli.StringSlice
var customFormKeywordValues *cli.StringSlice
var pushAddress string
var pushProxyPoolMax int
var pushProxyWG sync.WaitGroup
var outputJsonPath string
var logLevel string

func main() {
	author := cli.Author{
		Name:  "9ian1i",
		Email: "9ian1itp@gmail.com",
	}

	ignoreKeywords = cli.NewStringSlice(config.DefaultIgnoreKeywords...)
	customFormTypeValues = cli.NewStringSlice()
	customFormKeywordValues = cli.NewStringSlice()

	app := &cli.App{
		Name:        "crawlergo",
		Usage:       "A powerful browser crawler for web vulnerability scanners",
		UsageText:   "crawlergo [global options] url1 url2 url3 ... (must be same host)",
		Version:     "v0.4.2",
		Authors:     []*cli.Author{&author},
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:        "chromium-path",
				Aliases:     []string{"c"},
				Usage:       "`Path` of chromium executable. Such as \"/home/test/chrome-linux/chrome\"",
				Required:    true,
				Destination: &taskConfig.ChromiumPath,
				EnvVars:     []string{"CRAWLERGO_CHROMIUM_PATH"},
			},
			&cli.StringFlag{
				Name:        "custom-headers",
				Usage:       "add additional `Headers` to each request. The input string will be called json.Unmarshal",
				Value:       fmt.Sprintf(`{"Spider-Name": "crawlergo", "User-Agent": "%s"}`, config.DefaultUA),
				Destination: &taskConfig.ExtraHeadersString,
			},
			&cli.StringFlag{
				Name:        "post-data",
				Aliases:     []string{"d"},
				Usage:       "set `PostData` to target and use POST method.",
				Destination: &postData,
			},
			&cli.IntFlag{
				Name:        "max-crawled-count",
				Aliases:     []string{"m"},
				Value:       config.MaxCrawlCount,
				Usage:       "the maximum `Number` of URLs visited by the crawler in this task.",
				Destination: &taskConfig.MaxCrawlCount,
			},
			&cli.StringFlag{
				Name:        "filter-mode",
				Aliases:     []string{"f"},
				Value:       "smart",
				Usage:       "filtering `Mode` used for collected requests. Allowed mode:\"simple\", \"smart\" or \"strict\".",
				Destination: &taskConfig.FilterMode,
			},
			&cli.StringFlag{
				Name:        "output-mode",
				Aliases:     []string{"o"},
				Value:       "console",
				Usage:       "console print or serialize output. Allowed mode:\"console\" ,\"json\" or \"none\".",
				Destination: &outputMode,
			},
			&cli.StringFlag{
				Name:        "output-json",
				Usage:       "write output to a json file.Such as result_www_crawlergo_com.json",
				Destination: &outputJsonPath,
			},
			&cli.BoolFlag{
				Name:        "incognito-context",
				Aliases:     []string{"i"},
				Value:       true,
				Usage:       "whether the browser is launched in incognito mode.",
				Destination: &taskConfig.IncognitoContext,
			},
			&cli.IntFlag{
				Name:        "max-tab-count",
				Aliases:     []string{"t"},
				Value:       8,
				Usage:       "maximum `Number` of tabs allowed.",
				Destination: &taskConfig.MaxTabsCount,
			},
			&cli.BoolFlag{
				Name:        "fuzz-path",
				Value:       false,
				Usage:       "whether to fuzz the target with common paths.",
				Destination: &taskConfig.PathByFuzz,
			},
			&cli.PathFlag{
				Name:        "fuzz-path-dict",
				Usage:       "`Path` of fuzz dict. Such as \"/home/test/fuzz_path.txt\"",
				Destination: &taskConfig.FuzzDictPath,
			},
			&cli.BoolFlag{
				Name:        "robots-path",
				Value:       false,
				Usage:       "whether to resolve paths from /robots.txt.",
				Destination: &taskConfig.PathFromRobots,
			},
			&cli.StringFlag{
				Name:        "request-proxy",
				Usage:       "all requests connect through defined proxy server.",
				Destination: &taskConfig.Proxy,
			},
			//&cli.BoolFlag{
			//	Name:        "bypass",
			//	Value:       false,
			//	Usage:       "whether to encode url with detected charset.",
			//	Destination: &taskConfig.EncodeURLWithCharset,
			//},
			&cli.BoolFlag{
				Name:        "encode-url",
				Value:       false,
				Usage:       "whether to encode url with detected charset.",
				Destination: &taskConfig.EncodeURLWithCharset,
			},
			&cli.DurationFlag{
				Name:        "tab-run-timeout",
				Value:       config.TabRunTimeout,
				Usage:       "the `Timeout` of a single tab task.",
				Destination: &taskConfig.TabRunTimeout,
			},
			&cli.DurationFlag{
				Name:        "wait-dom-content-loaded-timeout",
				Value:       config.DomContentLoadedTimeout,
				Usage:       "the `Timeout` of waiting for a page dom ready.",
				Destination: &taskConfig.DomContentLoadedTimeout,
			},
			&cli.StringFlag{
				Name:        "event-trigger-mode",
				Value:       config.EventTriggerAsync,
				Usage:       "this `Value` determines how the crawler automatically triggers events.Allowed mode:\"async\" or \"sync\".",
				Destination: &taskConfig.EventTriggerMode,
			},
			&cli.DurationFlag{
				Name:        "event-trigger-interval",
				Value:       config.EventTriggerInterval,
				Usage:       "the `Interval` of triggering each event.",
				Destination: &taskConfig.EventTriggerInterval,
			},
			&cli.DurationFlag{
				Name:        "before-exit-delay",
				Value:       config.BeforeExitDelay,
				Usage:       "the `Time` of waiting before crawler exit.",
				Destination: &taskConfig.BeforeExitDelay,
			},
			&cli.StringSliceFlag{
				Name:        "ignore-url-keywords",
				Aliases:     []string{"iuk"},
				Value:       ignoreKeywords,
				Usage:       "crawlergo will not crawl these URLs matched by `Keywords`. e.g.: -iuk logout -iuk quit -iuk exit",
				DefaultText: "Default [logout quit exit]",
			},
			&cli.StringSliceFlag{
				Name:    "form-values",
				Aliases: []string{"fv"},
				Value:   customFormTypeValues,
				Usage:   "custom filling text for each form type. e.g.: -fv username=crawlergo_nice -fv password=admin123",
			},
			// 根据关键词自行选择填充文本
			&cli.StringSliceFlag{
				Name:    "form-keyword-values",
				Aliases: []string{"fkv"},
				Value:   customFormKeywordValues,
				Usage:   "custom filling text, fuzzy matched by keyword. e.g.: -fkv user=crawlergo_nice -fkv pass=admin123",
			},
			&cli.StringFlag{
				Name:        "push-to-proxy",
				Usage:       "every request in 'req_list' will be pushed to the proxy `Address`. Such as \"http://127.0.0.1:8080/\"",
				Destination: &pushAddress,
			},
			&cli.IntFlag{
				Name:        "push-pool-max",
				Usage:       "maximum `Number` of concurrency when pushing results to proxy.",
				Value:       DefaultMaxPushProxyPoolMax,
				Destination: &pushProxyPoolMax,
			},
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "log print `Level`, options include debug, info, warn, error and fatal.",
				Value:       DefaultLogLevel,
				Destination: &logLevel,
			},
			&cli.BoolFlag{
				Name:        "no-headless",
				Value:       false,
				Usage:       "no headless mode",
				Destination: &taskConfig.NoHeadless,
			},
		},
		Action: run,
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Logger.Fatal(err)
	}
}

func run(c *cli.Context) error {
	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	if c.Args().Len() == 0 {
		logger.Logger.Error("url must be set")
		return errors.New("url must be set")
	}

	// 设置日志输出级别
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	logger.Logger.SetLevel(level)

	var targets []*model2.Request
	for _, _url := range c.Args().Slice() {
		var req model2.Request
		url, err := model2.GetUrl(_url)
		if err != nil {
			logger.Logger.Error("parse url failed, ", err)
			continue
		}
		if postData != "" {
			req = model2.GetRequest(config.POST, url, getOption())
		} else {
			req = model2.GetRequest(config.GET, url, getOption())
		}
		req.Proxy = taskConfig.Proxy
		targets = append(targets, &req)
	}
	taskConfig.IgnoreKeywords = ignoreKeywords.Value()
	if taskConfig.Proxy != "" {
		logger.Logger.Info("request with proxy: ", taskConfig.Proxy)
	}

	if len(targets) == 0 {
		logger.Logger.Fatal("no validate target.")
	}

	// 检查自定义的表单参数配置
	taskConfig.CustomFormValues, err = parseCustomFormValues(customFormTypeValues.Value())
	if err != nil {
		logger.Logger.Fatal(err)
	}
	taskConfig.CustomFormKeywordValues, err = keywordStringToMap(customFormKeywordValues.Value())
	if err != nil {
		logger.Logger.Fatal(err)
	}

	// 开始爬虫任务
	task, err := pkg.NewCrawlerTask(targets, taskConfig)
	if err != nil {
		logger.Logger.Error("create crawler task failed.")
		os.Exit(-1)
	}
	if len(targets) != 0 {
		logger.Logger.Info(fmt.Sprintf("Init crawler task, host: %s, max tab count: %d, max crawl count: %d.",
			targets[0].URL.Host, taskConfig.MaxTabsCount, taskConfig.MaxCrawlCount))
		logger.Logger.Info("filter mode: ", taskConfig.FilterMode)
	}

	// 提示自定义表单填充参数
	if len(taskConfig.CustomFormValues) > 0 {
		logger.Logger.Info("Custom form values, " + tools.MapStringFormat(taskConfig.CustomFormValues))
	}
	// 提示自定义表单填充参数
	if len(taskConfig.CustomFormKeywordValues) > 0 {
		logger.Logger.Info("Custom form keyword values, " + tools.MapStringFormat(taskConfig.CustomFormKeywordValues))
	}
	if _, ok := taskConfig.CustomFormValues["default"]; !ok {
		logger.Logger.Info("If no matches, default form input text: " + config.DefaultInputText)
		taskConfig.CustomFormValues["default"] = config.DefaultInputText
	}

	go handleExit(task)
	logger.Logger.Info("Start crawling.")
	task.Run()
	result := task.Result

	logger.Logger.Info(fmt.Sprintf("Task finished, %d results, %d requests, %d subdomains, %d domains found.",
		len(result.ReqList), len(result.AllReqList), len(result.SubDomainList), len(result.AllDomainList)))

	// 内置请求代理
	if pushAddress != "" {
		logger.Logger.Info("pushing results to ", pushAddress, ", max pool number:", pushProxyPoolMax)
		Push2Proxy(result.ReqList)
	}

	// 输出结果
	outputResult(result)

	return nil
}

func getOption() model2.Options {
	var option model2.Options
	if postData != "" {
		option.PostData = postData
	}
	if taskConfig.ExtraHeadersString != "" {
		err := json.Unmarshal([]byte(taskConfig.ExtraHeadersString), &taskConfig.ExtraHeaders)
		if err != nil {
			logger.Logger.Fatal("custom headers can't be Unmarshal.")
			panic(err)
		}
		option.Headers = taskConfig.ExtraHeaders
	}
	return option
}

func parseCustomFormValues(customData []string) (map[string]string, error) {
	parsedData := map[string]string{}
	for _, item := range customData {
		keyValue := strings.Split(item, "=")
		if len(keyValue) < 2 {
			return nil, errors.New("invalid form item: " + item)
		}
		key := keyValue[0]
		if !tools.StringSliceContain(config.AllowedFormName, key) {
			return nil, errors.New("not allowed form key: " + key)
		}
		value := keyValue[1]
		parsedData[key] = value
	}
	return parsedData, nil
}

func keywordStringToMap(data []string) (map[string]string, error) {
	parsedData := map[string]string{}
	for _, item := range data {
		keyValue := strings.Split(item, "=")
		if len(keyValue) < 2 {
			return nil, errors.New("invalid keyword format: " + item)
		}
		key := keyValue[0]
		value := keyValue[1]
		parsedData[key] = value
	}
	return parsedData, nil
}

func outputResult(result *pkg.Result) {
	// 输出结果
	if outputMode == "json" {
		fmt.Println("--[Mission Complete]--")
		resBytes := getJsonSerialize(result)
		fmt.Println(string(resBytes))
	} else if outputMode == "console" {
		for _, req := range result.ReqList {
			req.FormatPrint()
		}
	}
	if len(outputJsonPath) != 0 {
		resBytes := getJsonSerialize(result)
		tools.WriteFile(outputJsonPath, resBytes)
	}
}

/**
原生被动代理推送支持
*/
func Push2Proxy(reqList []*model2.Request) {
	pool, _ := ants.NewPool(pushProxyPoolMax)
	defer pool.Release()
	for _, req := range reqList {
		task := ProxyTask{
			req:       req,
			pushProxy: pushAddress,
		}
		pushProxyWG.Add(1)
		go func() {
			err := pool.Submit(task.doRequest)
			if err != nil {
				logger.Logger.Error("add Push2Proxy task failed: ", err)
				pushProxyWG.Done()
			}
		}()
	}
	pushProxyWG.Wait()
}

/**
协程池请求的任务
*/
func (p *ProxyTask) doRequest() {
	defer pushProxyWG.Done()
	_, _ = requests.Request(p.req.Method, p.req.URL.String(), tools.ConvertHeaders(p.req.Headers), []byte(p.req.PostData),
		&requests.ReqOptions{Timeout: 1, AllowRedirect: false, Proxy: p.pushProxy})
}

func handleExit(t *pkg.CrawlerTask) {
	select {
	case <-signalChan:
		fmt.Println("exit ...")
		t.Pool.Tune(1)
		t.Pool.Release()
		t.Browser.Close()
		os.Exit(-1)
	}
}

func getJsonSerialize(result *pkg.Result) []byte {
	var res Result
	var reqList []Request
	var allReqList []Request
	for _, _req := range result.ReqList {
		var req Request
		req.Method = _req.Method
		req.Url = _req.URL.String()
		req.Source = _req.Source
		req.Data = _req.PostData
		req.Headers = _req.Headers
		reqList = append(reqList, req)
	}
	for _, _req := range result.AllReqList {
		var req Request
		req.Method = _req.Method
		req.Url = _req.URL.String()
		req.Source = _req.Source
		req.Data = _req.PostData
		req.Headers = _req.Headers
		allReqList = append(allReqList, req)
	}
	res.AllReqList = allReqList
	res.ReqList = reqList
	res.AllDomainList = result.AllDomainList
	res.SubDomainList = result.SubDomainList

	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Fatal("Marshal result error")
	}
	return resBytes
}
