package filter

import "github.com/Qianlitp/crawlergo/pkg/model"

type FilterHandler interface {
	DoFilter(req *model.Request) bool
}
