// Copyright 2026 kenny-not-dead. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
//
// Package gopptx provides functionality to create and manipulate PowerPoint
// (.pptx) files in Go, using the Office Open XML (ECMA-376) format.

package gopptx

import (
	"bytes"
	"encoding/xml"
	"io"
)

// themeReader provides a function to get the pointer to the xl/theme/theme1.xml
// structure after deserialization.
func (f *File) themeReader() (*decodeTheme, error) {
	if _, ok := f.Pkg.Load(defaultXMLPathTheme); !ok {
		return nil, nil
	}
	theme := decodeTheme{}
	data := namespaceStrictToTransitional(f.readXML(defaultXMLPathTheme))
	// capture the origin root element attributes, so the serialized theme
	// keeps its origin root element on save
	if _, ok := f.xmlAttr.Load(defaultXMLPathTheme); !ok {
		d := f.xmlNewDecoder(bytes.NewReader(data))
		for {
			token, err := d.Token()
			if err != nil {
				break
			}
			if startElement, ok := token.(xml.StartElement); ok {
				f.xmlAttr.Store(defaultXMLPathTheme, startElement.Attr)
				break
			}
		}
	}
	if err := f.xmlNewDecoder(bytes.NewReader(data)).
		Decode(&theme); err != nil && err != io.EOF {
		return &theme, err
	}
	return &theme, nil
}
