// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"reflect"
)

var xmlSpecial = map[byte]string{
	'<':  "&lt;",
	'>':  "&gt;",
	'"':  "&quot;",
	'\'': "&apos;",
	'&':  "&amp;",
}

func xmlEscape(s string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); i++ {
		c := s[i]
		if s, ok := xmlSpecial[c]; ok {
			b.WriteString(s)
		} else {
			b.WriteByte(c)
		}
	}
	return b.String()
}

// Scan XML token stream to finc next Element (start or end)
func nextElement(p *xml.Decoder) (xml.Token, error) {
	for {
		t, err := p.Token()
		if err != nil {
			return xml.StartElement{}, err
		}

		switch elem := t.(type) {
		case xml.StartElement, xml.EndElement:
			return t, nil
		case xml.CharData:
			// https://xmpp.org/rfcs/rfc6120.html#xml-whitespace
			if string(elem) == " " {
				log.Println("xmpp: received whitespace ping")
			}
		default:
			log.Printf("xmpp: received unhandled element: %#v\n", elem)
		}
	}
}

// Scan XML token stream to find next StartElement.
func nextStart(p *xml.Decoder) (xml.StartElement, error) {
	for {
		t, err := nextElement(p)
		if err != nil {
			return xml.StartElement{}, err
		}

		if start, ok := t.(xml.StartElement); ok {
			return start, nil
		}
	}
}

// Scan XML token stream for next element and save into val.
// If val == nil, allocate new element based on proto map.
// Either way, return val.
func next(c *Conn) (xml.Name, interface{}, error) {
	elem, err := nextElement(c.in)
	if err != nil {
		return xml.Name{}, nil, err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	switch el := elem.(type) {
	case xml.StartElement:
		return decodeStartElement(c, el)
	case xml.EndElement:
		return decodeEndElement(el)
	}

	return xml.Name{}, nil, fmt.Errorf("unexpected element %s", elem)
}

func decodeStartElement(c *Conn, se xml.StartElement) (xml.Name, interface{}, error) {

	// Put it in an interface and allocate one.
	var nv interface{}
	if t, e := c.customStorage[se.Name]; e {
		nv = reflect.New(t).Interface()
	} else if t, e := defaultStorage[se.Name]; e {
		nv = reflect.New(t).Interface()
	} else {
		return xml.Name{}, nil, errors.New("unexpected XMPP message " +
			se.Name.Space + " <" + se.Name.Local + "/>")
	}

	// Unmarshal into that storage.
	if err := c.in.DecodeElement(nv, &se); err != nil {
		return xml.Name{}, nil, err
	}

	return se.Name, nv, nil
}

func decodeEndElement(ce xml.EndElement) (xml.Name, interface{}, error) {
	switch ce.Name {
	case xml.Name{NsStream, "stream"}:
		return ce.Name, &StreamClose{}, nil
	}

	return ce.Name, nil, nil
}

var defaultStorage = map[xml.Name]reflect.Type{
	xml.Name{Space: NsStream, Local: "features"}: reflect.TypeOf(streamFeatures{}),
	xml.Name{Space: NsStream, Local: "error"}:    reflect.TypeOf(StreamError{}),
	xml.Name{Space: NsTLS, Local: "starttls"}:    reflect.TypeOf(tlsStartTLS{}),
	xml.Name{Space: NsTLS, Local: "proceed"}:     reflect.TypeOf(tlsProceed{}),
	xml.Name{Space: NsTLS, Local: "failure"}:     reflect.TypeOf(tlsFailure{}),
	xml.Name{Space: NsSASL, Local: "mechanisms"}: reflect.TypeOf(saslMechanisms{}),
	xml.Name{Space: NsSASL, Local: "challenge"}:  reflect.TypeOf(""),
	xml.Name{Space: NsSASL, Local: "response"}:   reflect.TypeOf(""),
	xml.Name{Space: NsSASL, Local: "abort"}:      reflect.TypeOf(saslAbort{}),
	xml.Name{Space: NsSASL, Local: "success"}:    reflect.TypeOf(saslSuccess{}),
	xml.Name{Space: NsSASL, Local: "failure"}:    reflect.TypeOf(saslFailure{}),
	xml.Name{Space: NsBind, Local: "bind"}:       reflect.TypeOf(bindBind{}),
	xml.Name{Space: NsClient, Local: "message"}:  reflect.TypeOf(ClientMessage{}),
	xml.Name{Space: NsClient, Local: "presence"}: reflect.TypeOf(ClientPresence{}),
	xml.Name{Space: NsClient, Local: "iq"}:       reflect.TypeOf(ClientIQ{}),
	xml.Name{Space: NsClient, Local: "error"}:    reflect.TypeOf(ClientError{}),
}
