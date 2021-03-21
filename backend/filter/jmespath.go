package filter

import (
	"encoding/json"
	"fmt"

	"github.com/jmespath/go-jmespath"
)

var (
	jmesPathSearcher       = jmespath.Search
	errInputIsNotValidJSON = fmt.Errorf(`unable to apply JMESPath filter to non-JSON data`)
	errUnknownDataType     = fmt.Errorf(`unknown data type`)
)

// JMESPath is a JMESPath filter
type JMESPath struct {
	Expression string `json:"expression"`
}

// StripUnsafe removes unsafe characters to prevent cli injection
func (j *JMESPath) StripUnsafe() error {
	// not applicable to JMESPath
	return nil
}

// Run applies a JMESPath filter
func (j *JMESPath) Run(bdy []byte) (out interface{}, err error) {
	if j.Expression == "" {
		// TODO: make this a method on our interface or something to make it apply to all filters (e.g. return on empty filter expression)
		out = string(bdy)
		return
	}

	var data interface{}
	if err = json.Unmarshal(bdy, &data); err != nil {
		err = errInputIsNotValidJSON
		return
	}

	out, err = jmesPathSearcher(j.Expression, data)
	if err != nil {
		return
	}

	return
}
