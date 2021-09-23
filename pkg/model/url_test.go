package model

import (
	"net/url"
	"testing"

	"golang.org/x/net/publicsuffix"
)

var (
	rootDomainTestCases = []struct {
		domain     string
		rootDomain string
		wantICANN  bool
	}{
		{"www.amazon.co.uk", "amazon.co.uk", true},
		{"www.baidu.com", "baidu.com", true},
		{"www.baidu.com.cn", "baidu.com.cn", true},
		{"www.pku.edu.cn", "pku.edu.cn", true},
		{"www.example1.debian.org", "debian.org", true},
		{"www.golang.dev", "golang.dev", true},
		// 以下都是一些特殊的 case，主要包括一些特殊的域名和私有域名，一般情况遇不到
		// error domains
		{"com.cn", "", true},
		// not an icann domain
		{"www.example0.debian.net", "", false},
		{"s3.cn-north-1.amazonaws.com.cn", "", false},
		{"www.0emm.com", "", false},
		{"there.is.no.such-tld", "", false},
	}
)

func TestRootDomain(t *testing.T) {
	for _, tc := range rootDomainTestCases {
		u := &URL{url.URL{Host: tc.domain}}
		rootDomain := u.RootDomain()
		_, icann := publicsuffix.PublicSuffix(u.Hostname())
		if rootDomain != tc.rootDomain {
			t.Errorf("%s parse root domain failed", tc.domain)
		}
		if icann != tc.wantICANN {
			t.Errorf("%s not an icann domain", tc.domain)
		}
	}
}
