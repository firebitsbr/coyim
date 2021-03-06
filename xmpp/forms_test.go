package xmpp

import (
	"encoding/xml"
	"errors"

	. "gopkg.in/check.v1"
)

type FormsXmppSuite struct{}

var _ = Suite(&FormsXmppSuite{})

func (s *FormsXmppSuite) Test_processForm_returnsErrorFromCallback(c *C) {
	e := errors.New("some kind of error")
	f := &Form{}
	_, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return e
	})

	c.Assert(err, Equals, e)
}

func (s *FormsXmppSuite) Test_processForm_returnsEmptySubmitFormForEmptyForm(c *C) {
	f := &Form{}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, Form{Type: "submit"})
}

func (s *FormsXmppSuite) Test_processForm_returnsFixedFields(c *C) {
	f := &Form{}
	f.Fields = []formField{
		formField{
			Label:  "hello",
			Type:   "fixed",
			Values: []string{"Something"},
		},
		formField{
			Label: "hello2",
			Type:  "fixed",
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields:       nil})
}

func (s *FormsXmppSuite) Test_processForm_returnsBooleanFields(c *C) {
	f := &Form{}
	f.Fields = []formField{
		formField{
			Label: "hello3",
			Type:  "boolean",
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []formField{
			formField{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*formFieldRequired)(nil),
				Values:   []string{"false"},
				Options:  []formFieldOption(nil),
				Media:    []formFieldMedia(nil)}}})
}

func (s *FormsXmppSuite) Test_processForm_returnsMultiFields(c *C) {
	f := &Form{}
	f.Fields = []formField{
		formField{
			Label: "hello4",
			Type:  "jid-multi",
		},
		formField{
			Label: "hello5",
			Type:  "text-multi",
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []formField{
			formField{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*formFieldRequired)(nil),
				Values:   []string(nil),
				Options:  []formFieldOption(nil),
				Media:    []formFieldMedia(nil)},
			formField{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*formFieldRequired)(nil),
				Values:   []string(nil),
				Options:  []formFieldOption(nil),
				Media:    []formFieldMedia(nil)}}})
}

func (s *FormsXmppSuite) Test_processForm_returnsListSingle(c *C) {
	f := &Form{}
	f.Fields = []formField{
		formField{
			Label: "hello7",
			Type:  "list-single",
			Options: []formFieldOption{
				formFieldOption{Label: "One", Value: "Two"},
				formFieldOption{Label: "Three", Value: "Four"},
			},
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []formField{
			formField{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*formFieldRequired)(nil),
				Values:   []string{"Two"},
				Options:  []formFieldOption(nil), Media: []formFieldMedia(nil)}}})
}

func (s *FormsXmppSuite) Test_processForm_returnsListMulti(c *C) {
	f := &Form{}
	f.Fields = []formField{
		formField{
			Label: "hello1o7",
			Type:  "list-multi",
			Options: []formFieldOption{
				formFieldOption{Label: "One", Value: "Two"},
				formFieldOption{Label: "Three", Value: "Four"},
			},
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []formField{
			formField{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*formFieldRequired)(nil),
				Values:   []string(nil),
				Options:  []formFieldOption(nil), Media: []formFieldMedia(nil)}}})
}

func (s *FormsXmppSuite) Test_processForm_returnsHidden(c *C) {
	f := &Form{}
	f.Fields = []formField{
		formField{
			Label: "hello1o71",
			Type:  "hidden",
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []formField{
			formField{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*formFieldRequired)(nil),
				Values:   []string(nil),
				Options:  []formFieldOption(nil), Media: []formFieldMedia(nil)}}})
}

func (s *FormsXmppSuite) Test_processForm_returnsUnknown(c *C) {
	f := &Form{}
	f.Fields = []formField{
		formField{
			Label: "hello1o71",
			Type:  "another-fancy-type",
		},
		formField{
			Label:  "hello1o73",
			Type:   "another-fancy-type",
			Values: []string{"another one"},
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []formField{
			formField{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*formFieldRequired)(nil),
				Values:   []string{""},
				Options:  []formFieldOption(nil),
				Media:    []formFieldMedia(nil)},
			formField{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*formFieldRequired)(nil),
				Values:   []string{""},
				Options:  []formFieldOption(nil),
				Media:    []formFieldMedia(nil)}}})
}

type testOtherFormType struct{}

func (s *FormsXmppSuite) Test_processForm_panicsWhenGivenAWeirdFormType(c *C) {
	f := &Form{}
	f.Fields = []formField{
		formField{
			Label: "hello1o71",
			Type:  "another-fancy-type",
		},
	}
	c.Assert(func() {
		processForm(f, nil, func(title, instructions string, fields []interface{}) error {
			fields[0] = testOtherFormType{}
			return nil
		})
	}, PanicMatches, "unknown field type in result from callback: xmpp.testOtherFormType")
}

func (s *FormsXmppSuite) Test_processForm_setsAValidBooleanReturnValue(c *C) {
	f := &Form{}
	f.Fields = []formField{
		formField{
			Label: "hello1o71",
			Type:  "boolean",
		},
	}
	f2, _ := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		fields[0].(*BooleanFormField).Result = true
		return nil
	})
	c.Assert(*f2, DeepEquals, Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []formField{
			formField{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*formFieldRequired)(nil),
				Values:   []string{"true"},
				Options:  []formFieldOption(nil),
				Media:    []formFieldMedia(nil)}}})
}

func (s *FormsXmppSuite) Test_processForm_returnsListMultiWithResults(c *C) {
	f := &Form{}
	f.Fields = []formField{
		formField{
			Label: "hello1o7",
			Type:  "list-multi",
			Options: []formFieldOption{
				formFieldOption{Label: "One", Value: "Two"},
				formFieldOption{Label: "Three", Value: "Four"},
			},
		},
	}
	f2, err := processForm(f, nil, func(title, instructions string, fields []interface{}) error {
		fields[0].(*MultiSelectionFormField).Results = []int{1}
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []formField{
			formField{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: (*formFieldRequired)(nil),
				Values:   []string{"Four"},
				Options:  []formFieldOption(nil), Media: []formFieldMedia(nil)}}})
}

func (s *FormsXmppSuite) Test_processForm_dealsWithMediaCorrectly(c *C) {
	f := &Form{}
	datas := []bobData{
		bobData{
			CID:    "foobax",
			Base64: ".....",
		},
		bobData{
			CID:    "foobar",
			Base64: "aGVsbG8=",
		},
	}
	f.Fields = []formField{
		formField{
			Label: "hello1o7",
			Type:  "hidden",
			Media: []formFieldMedia{
				formFieldMedia{
					URIs: []mediaURI{
						mediaURI{
							MIMEType: "",
							URI:      "",
						},
						mediaURI{
							MIMEType: "",
							URI:      "hello:world",
						},
						mediaURI{
							MIMEType: "",
							URI:      "cid:foobar",
						},
						mediaURI{
							MIMEType: "",
							URI:      "cid:foobax",
						},
					},
				},
			},
		},
	}
	f2, err := processForm(f, datas, func(title, instructions string, fields []interface{}) error {
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(*f2, DeepEquals, Form{
		XMLName:      xml.Name{Space: "", Local: ""},
		Type:         "submit",
		Title:        "",
		Instructions: "",
		Fields: []formField{
			formField{
				XMLName:  xml.Name{Space: "", Local: ""},
				Desc:     "",
				Var:      "",
				Type:     "",
				Label:    "",
				Required: nil,
				Values:   nil,
				Options:  nil,
				Media:    nil}}})
}
