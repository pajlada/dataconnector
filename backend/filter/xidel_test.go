package filter

/*
func TestXidelRun(t *testing.T) {
	Runner = mockRunner
	defer func() { Runner = exec.Command }()

	for _, c := range []struct {
		name       string
		bdy        []byte
		expression string
		want       []byte
		wantErr    error
	}{
		{
			name:       "XPath array of arrays query",
			bdy:        []byte(`<?xml version="1.0" encoding="UTF-8"?><bookstore><book><title lang="en">Harry Potter</title><price>29.99</price></book><book><title lang="en">Learning XML</title><price>39.95</price></book></bookstore>`),
			expression: "//title",
			want:       []byte(`[["Harry Potter", "Learning XML"]]`),
			wantErr:    nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			f := &Xidel{
				Expression: c.expression,
			}

			got, gotErr := f.Run(c.bdy)
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
	fmt.Fprintf(os.Stdout, `[["Harry Potter", "Learning XML"]]`)
	os.Exit(0)
}
*/
