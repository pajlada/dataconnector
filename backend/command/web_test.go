package command

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/chromedp/chromedp"
)

func TestWebValid(t *testing.T) {
	for _, c := range []struct {
		name    string
		web     *Web
		wantErr error
	}{
		{
			name: "validates",
			web: &Web{
				URL: "https://example.com",
			},
			wantErr: nil,
		},
		{
			name: "invalid url should return error",
			web: &Web{
				URL: "not a url",
			},
			wantErr: errInvalidURL,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			gotErr := c.web.Valid()
			if !errors.Is(gotErr, c.wantErr) {
				t.Fatalf("got %v; want %v", gotErr, c.wantErr)
			}
		})
	}
}

func TestWebDeParameterize(t *testing.T) {
	for _, c := range []struct {
		name    string
		web     *Web
		params  []string
		want    *Web
		wantErr error
	}{
		{
			name: "deparameterizes correctly",
			web: &Web{
				URL: "https://example.com/+++1+++?a=+++2+++",
				Headers: []Header{
					{Key: "+++3+++", Value: "+++4+++"},
				},
			},
			params: []string{"path", "b", "header_key", "header_value"},
			want: &Web{
				URL: "https://example.com/path?a=b",
				Headers: []Header{
					{Key: "header_key", Value: "header_value"},
				},
			},
			wantErr: nil,
		},
		{
			name: "returns on unhandled url parameters",
			web: &Web{
				URL: "https://example.com/+++1+++?a=+++2+++",
			},
			params: []string{"path"},
			want: &Web{
				URL: "https://example.com/path?a=+++2+++",
			},
			wantErr: errUnhandledParams,
		},
		{
			name: "returns on unhandled header parameters",
			web: &Web{
				URL: "https://example.com",
				Headers: []Header{
					{Key: "+++1+++", Value: "+++2+++"},
				},
			},
			params: []string{"header_key"},
			want: &Web{
				URL: "https://example.com",
				Headers: []Header{
					{Key: "header_key", Value: "+++2+++"},
				},
			},
			wantErr: errUnhandledParams,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			gotErr := c.web.DeParameterize(c.params)
			if !errors.Is(gotErr, c.wantErr) {
				t.Fatalf("got %v; want %v", gotErr, c.wantErr)
			}

			if !reflect.DeepEqual(c.web, c.want) {
				t.Fatalf("got %+v; want %+v", c.web, c.want)
			}
		})
	}
}

func TestWebRun(t *testing.T) {
	// we could also test this with a file:// but seems unnecessary
	for _, c := range []struct {
		name    string
		web     *Web
		want    []byte
		wantErr error
	}{
		{
			name: "runs successfully",
			web: &Web{
				URL: "https://example.com",
			},
			want:    []byte("<html><body>hello world!</body></html>"),
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			c.web.runner = c.web.mockChromedpRunner
			got, gotErr := c.web.Run()
			if !errors.Is(gotErr, c.wantErr) {
				t.Fatalf("got %v; want %v", gotErr, c.wantErr)
			}

			if !bytes.Equal(got, c.want) {
				t.Fatalf("got %s; want %s", got, c.want)
			}
		})
	}
}

func (w *Web) mockChromedpRunner(ctx context.Context, actions ...chromedp.Action) error {
	w.result = "<html><body>hello world!</body></html>"
	return nil
}
