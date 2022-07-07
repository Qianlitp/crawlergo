package config_test

import (
	"testing"

	"github.com/Qianlitp/crawlergo/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestStaticSuffix(t *testing.T) {
	assert.Equal(t, true, config.StaticSuffixSet.Contains("png"))
	assert.Equal(t, false, config.StaticSuffixSet.Contains("demo"))

	assert.Equal(t, true, config.ScriptSuffixSet.Contains("asp"))
	assert.Equal(t, false, config.ScriptSuffixSet.Contains("demo"))
}
