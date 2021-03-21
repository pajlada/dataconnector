package filter

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jmespath/go-jmespath"
)

func TestJMESPathStripUnsafe(t *testing.T) {
	for _, c := range []struct {
		name     string
		jmespath *JMESPath
		bdy      []byte
		want     interface{}
		wantErr  error
	}{
		{
			name: "returns successfully",
			jmespath: &JMESPath{
				Expression: "myexpression",
			},
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			gotErr := c.jmespath.StripUnsafe()
			if !errors.Is(gotErr, c.wantErr) {
				t.Fatalf("got %v; want %v", gotErr, c.wantErr)
			}
		})
	}
}

func TestJMESPathRun(t *testing.T) {
	jmesPathSearcher = mockJMESPathSearch
	defer func() { jmesPathSearcher = jmespath.Search }()

	for _, c := range []struct {
		name     string
		jmespath *JMESPath
		bdy      []byte
		want     interface{}
		wantErr  error
	}{
		{
			name: "returns successfully",
			jmespath: &JMESPath{
				Expression: "myexpression",
			},
			bdy:     []byte(`{"mydata": [["bob","fred"]]}`),
			want:    []string{"bob", "fred"},
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, gotErr := c.jmespath.Run(c.bdy)
			if !errors.Is(gotErr, c.wantErr) {
				t.Fatalf("got %v; want %v", gotErr, c.wantErr)
			}

			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %s; want %s", got, c.want)
			}
		})
	}
}

func mockJMESPathSearch(expression string, data interface{}) (interface{}, error) {
	return []string{"bob", "fred"}, nil
}
