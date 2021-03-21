package command

/*
func TestCurlValid(t *testing.T) {
	for _, c := range []struct {
		name    string
		command string
		want    error
	}{
		{
			name:    `must start with "curl"`,
			command: "rm -r *",
			want:    errInvalidCurlCommand,
		},
		{
			name:    `cannot contain "|"`,
			command: "curl 'http://www.example.com' | rm *",
			want:    errInvalidCurlCommand,
		},
		{
			name:    `cannot contain "&"`,
			command: "curl 'http://www.example.com' & rm *",
			want:    errInvalidCurlCommand,
		},
		{
			name:    `cannot contain "|"`,
			command: "curl 'http://www.example.com'; rm *",
			want:    errInvalidCurlCommand,
		},
		{
			name:    "url is valid",
			command: "curl 'https://www.example.com/my-path?this=that'",
			want:    nil,
		},
		{
			name:    "scheme is invalid",
			command: "curl '+++1+++://example.com",
			want:    errInvalidURL,
		},
		{
			name:    "host is invalid",
			command: "curl 'https://+++2+++/",
			want:    errInvalidURL,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := &Curl{
				Command: c.command,
			}

			got := cmd.Valid()
			if !errors.Is(got, c.want) {
				t.Fatalf("got %v; want %v", got, c.want)
			}
		})
	}
}

func TestCurlDeParameterize(t *testing.T) {
	for _, c := range []struct {
		name    string
		command string
		params  []string
		want    string
		wantErr error
	}{
		{
			name:    "parameterize path",
			command: "curl 'https://www.example.com/+++1+++?a=b'",
			params:  []string{"another-path"},
			want:    "curl 'https://www.example.com/another-path?a=b'",
			wantErr: nil,
		},
		{
			name:    "parameterize query",
			command: "curl 'https://www.example.com?this=+++1+++'",
			params:  []string{"that"},
			want:    "curl 'https://www.example.com?this=that'",
			wantErr: nil,
		},
		{
			name:    "parameterize headers",
			command: `curl 'https://www.example.com' -H "API Key: +++1+++" -H "Content-Type: application/json"`,
			params:  []string{"my-key"},
			want:    `curl 'https://www.example.com' -H "API Key: my-key" -H "Content-Type: application/json"`,
			wantErr: nil,
		},
		{
			name:    "parameterize multiple parameters",
			command: "curl https://www.example.com/+++2+++/?a=+++1+++&c=+++3+++",
			params:  []string{"b", "mypath", "d"},
			want:    "curl https://www.example.com/mypath/?a=b&c=d",
			wantErr: nil,
		},
		{
			name:    "unhandled parameters",
			command: "curl https://www.example.com/+++2+++/?a=+++1+++&c=+++3+++&e=+++4+++",
			params:  []string{"b", "mypath", "d"},
			want:    "",
			wantErr: errUnhandledParams,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := &Curl{
				Command: c.command,
			}

			gotErr := cmd.DeParameterize(c.params)
			if !errors.Is(gotErr, c.wantErr) {
				t.Fatalf("got %v; want %v", gotErr, c.wantErr)
			}

			if !reflect.DeepEqual(cmd.Command, c.want) {
				t.Fatalf("got %q; want %q", cmd.Command, c.want)
			}
		})
	}
}

func TestCurlRun(t *testing.T) {
	Runner = mockRunner
	defer func() { Runner = exec.Command }()

	for _, c := range []struct {
		name    string
		command string
		params  []string
		want    []byte
		wantErr error
	}{
		{
			name:    "test returned data",
			command: "curl 'https://www.example.com'",
			want:    []byte(`{"result": "12"}`),
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := &Curl{
				Command: c.command,
			}

			got, gotErr := cmd.Run()
			if !errors.Is(gotErr, c.wantErr) {
				t.Fatalf("got %v; want %v", gotErr, c.wantErr)
			}

			if !bytes.Equal(got, c.want) {
				t.Fatalf("got %s; want %s", got, c.want)
			}
		})
	}
}

func mockRunner(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, `{"result": "12"}`)
	os.Exit(0)
}
*/
