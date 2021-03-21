package filter

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

var (
	// Runner lets us mock exec.Command for testing
	Runner = exec.Command
	// ErrUnknownFilter indicates the type of filter is not recognized
	ErrUnknownFilter = fmt.Errorf("unknown filter type")

	jmesPathType = "jmespath"
	jqType       = "jq"
	jsonPathType = "jsonpath"
	pupType      = "pup"
	xPathType    = "xpath"
	// TODO: more FilterTypes
)

// Filter is an individual filter
type Filter struct {
	Type     string `json:"type"`
	Filterer `json:"-"`
}

// Filterer outlines methods to filter data
type Filterer interface {
	StripUnsafe() error
	Run(bdy []byte) (out interface{}, err error)
}

// MarshalJSON encodes our Filterer as `filter`
func (f *Filter) MarshalJSON() ([]byte, error) {
	type Alias Filter
	return json.Marshal(&struct {
		*Alias
		Filter interface{} `json:"filter"`
	}{
		Alias:  (*Alias)(f),
		Filter: f.Filterer,
	})
}

// UnmarshalJSON unmarshals our Filterer interface to a struct
func (f *Filter) UnmarshalJSON(b []byte) error {
	type Alias Filter
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(f),
	}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	var m map[string]*json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	filterer, ok := m["filter"]
	if !ok || filterer == nil {
		return fmt.Errorf("missing filter")
	}

	switch f.Type {
	case jmesPathType:
		jmes := &JMESPath{}
		if err := json.Unmarshal(*filterer, &jmes); err != nil {
			return err
		}
		f.Filterer = jmes
	case jqType: // TODO
	case jsonPathType: // TODO
	case pupType: // TODO
	case xPathType: // TODO
	default:
		return ErrUnknownFilter
	}

	return nil
}
