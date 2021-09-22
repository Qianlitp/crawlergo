package filter

import (
	"crawlergo/pkg/config"
	"crawlergo/pkg/model"
	"github.com/deckarep/golang-set"
	"strings"
)

type SimpleFilter struct {
	UniqueSet mapset.Set
	HostLimit string
}

/**
需要过滤则返回 true
*/
func (s *SimpleFilter) DoFilter(req *model.Request) bool {
	if s.UniqueSet == nil {
		s.UniqueSet = mapset.NewSet()
	}
	// 首先判断是否需要过滤域名
	if s.HostLimit != "" && s.DomainFilter(req) {
		return true
	}
	// 去重
	if s.UniqueFilter(req) {
		return true
	}
	// 过滤静态资源
	if s.StaticFilter(req) {
		return true
	}
	return false
}

/**
请求去重
*/
func (s *SimpleFilter) UniqueFilter(req *model.Request) bool {
	if s.UniqueSet == nil {
		s.UniqueSet = mapset.NewSet()
	}
	if s.UniqueSet.Contains(req.UniqueId()) {
		return true
	} else {
		s.UniqueSet.Add(req.UniqueId())
		return false
	}
}

/**
静态资源过滤
*/
func (s *SimpleFilter) StaticFilter(req *model.Request) bool {
	if s.UniqueSet == nil {
		s.UniqueSet = mapset.NewSet()
	}
	// 首先将slice转换成map
	extMap := map[string]int{}
	staticSuffix := append(config.StaticSuffix, "js", "css", "json")
	for _, suffix := range staticSuffix {
		extMap[suffix] = 1
	}

	if req.URL.FileExt() == "" {
		return false
	}
	if _, ok := extMap[req.URL.FileExt()]; ok {
		return true
	}
	return false
}

/**
只保留指定域名的链接
*/
func (s *SimpleFilter) DomainFilter(req *model.Request) bool {
	if s.UniqueSet == nil {
		s.UniqueSet = mapset.NewSet()
	}
	if req.URL.Host == s.HostLimit || req.URL.Hostname() == s.HostLimit {
		return false
	}
	if strings.HasSuffix(s.HostLimit, ":80") && req.URL.Port() == "" && req.URL.Scheme == "http" {
		if req.URL.Hostname()+":80" == s.HostLimit {
			return false
		}
	}
	if strings.HasSuffix(s.HostLimit, ":443") && req.URL.Port() == "" && req.URL.Scheme == "https" {
		if req.URL.Hostname()+":443" == s.HostLimit {
			return false
		}
	}
	return true
}
