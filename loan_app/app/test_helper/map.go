package test_helper

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"

	"github.com/goProjects/loan_app/lib/utils"
)

func AssertEqualMap(t *testing.T, expected map[string]interface{}, actual map[string]interface{}, ignoreKeys []string) bool {

	var ignoreTopLevelEntries = func(key string, val interface{}) bool {
		return utils.StringContains(key, ignoreKeys, true)
	}

	cmpIgnoreOption := cmpopts.IgnoreMapEntries(ignoreTopLevelEntries)

	isEqual := assert.True(t, cmp.Equal(expected, actual, cmpIgnoreOption))

	if !isEqual {
		t.Logf("Diff: %s", cmp.Diff(expected, actual, cmpIgnoreOption))
	}

	return isEqual
}
