package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"cal.ope.cool/pkg/nhl"
	"github.com/stretchr/testify/assert"
)

func Test_Tmpl(t *testing.T) {
	fixture, err := ioutil.ReadFile("./kraken.json")
	assert.NoError(t, err)

	var sched nhl.Schedule
	err = json.Unmarshal(fixture, &sched)
	assert.NoError(t, err)

	for _, v := range *sched.Dates {
		for _, g := range *v.Games {
			fmt.Println(CreateDescription(&g))
			assert.Equal(t, "", CreateDescription(&g))
			break
		}
	}
}
