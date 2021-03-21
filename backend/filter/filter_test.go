package filter

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMarshalAndUnmarshalJSON(t *testing.T) {
	for _, c := range []struct {
		name   string
		filter *Filter
	}{
		{
			name: "marshals and unmarshals a  filter commander",
			filter: &Filter{
				Type: jmesPathType,
				Filterer: &JMESPath{
					Expression: "[[value]]",
				},
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			gotBytes, err := json.Marshal(c.filter)
			if err != nil {
				t.Fatal(err)
			}

			var got *Filter
			if err := json.Unmarshal(gotBytes, &got); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, c.filter) {
				t.Fatalf("got %+v; want %+v", got, c.filter)
			}
		})
	}
}
