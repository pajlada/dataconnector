package filter

// Xidel is a filter that is used by XPath
type Xidel struct {
	Expression string `json:"expression"`
}

/*
// StripUnsafe removes potentially harmful characters from a Xidel command
func (x *Xidel) StripUnsafe() error {
	// https://wiki.owasp.org/index.php/Testing_for_Command_Injection_(OTG-INPVAL-013)
	// { }  ( ) < > & * ‘ | = ? ; [ ]  $ – # ~ ! . ” %  / \ : + , `
	return nil
}

// Run applies a filter using Xidel
func (x *Xidel) Run(bdy []byte) (out interface{}, err error) {
	// NOTE: in tests, test that the statement below starts with "XIDEL" since we don't check for it above...also, change the curl way to this....e.g. don't require "curl", put it there for them
	bdy = []byte(`<?xml version="1.0" encoding="UTF-8"?><bookstore><book><title lang="en">Harry Potter</title><author>J K. Rowling</author><year>2005</year><price>29.99</price></book></bookstore>`)

	cmd := Runner("bash", "-c", fmt.Sprintf(`xidel -s '%s' -e '%s' --output-format=json-wrapped`, string(bdy), x.Expression))
	out, err = cmd.CombinedOutput()
	return
}
*/
