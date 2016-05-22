package xmljsonproxy

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"strconv"
	"strings"
)

func coerce(s string) interface{} {
	var err error

	i, err := strconv.ParseInt(s, 10, 64)

	if err == nil {
		return i
	}

	fl, err := strconv.ParseFloat(s, 64)

	if err == nil {
		return fl
	}

	b, err := strconv.ParseBool(s)

	if err == nil {
		return b
	}

	return s
}

func transformObject(d *xml.Decoder, st xml.StartElement) (interface{}, error) {
	var mapOutput map[string]interface{} = make(map[string]interface{})
	var token xml.Token

	for _, e := range st.Attr {
		mapOutput[e.Name.Local] = coerce(e.Value)
	}

	var stringOutput string

Loop:
	for true {
		token, _ = d.Token()

		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == "rowset" {
				name, rowset, err := transformRowset(d, token)
				if err != nil {
					return nil, err
				}
				mapOutput[name] = rowset
			} else {
				field, err := transformObject(d, token)
				if err != nil {
					return nil, err
				}
				mapOutput[token.Name.Local] = field
			}
			break
		case xml.EndElement:
			if token.Name == st.Name {
				break Loop
			}
		case xml.CharData:
			r := string(token)
			if strings.TrimSpace(r) != "" {
				stringOutput = r
			}

		}
	}

	var output interface{}

	if stringOutput != "" {
		output = stringOutput
	} else {
		output = mapOutput
	}

	return output, nil
}

func transformRowset(d *xml.Decoder, st xml.StartElement) (string, []interface{}, error) {
	var output []interface{}
	var token xml.Token

	var name string
	for _, e := range st.Attr {
		if e.Name.Local == "name" {
			name = e.Value
			break
		}
	}

Loop:
	for true {
		token, _ = d.Token()

		switch token := token.(type) {
		case xml.StartElement:
			row, err := transformObject(d, token)
			if err != nil {
				return name, nil, err
			}
			output = append(output, row)
		case xml.EndElement:
			if token.Name.Local == "rowset" {
				break Loop
			}
		}
	}

	return name, output, nil
}

func Transform(body io.Reader) ([]byte, error) {
	d := xml.NewDecoder(body)
	var token xml.Token
	var err error

	var output map[string]interface{} = make(map[string]interface{})

	for true {
		token, err = d.Token()
		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			switch token.Name.Local {
			case "eveapi":
				for _, e := range token.Attr {
					output[e.Name.Local] = coerce(e.Value)
				}
			default:
				field, err := transformObject(d, token)

				if err != nil {
					return nil, err
				}
				output[token.Name.Local] = field
			}
		}

	}

	return json.Marshal(output)
}

func TransformString(s string) (string, error) {
	buffer := bytes.NewBufferString(s)
	b, err := Transform(buffer)

	return string(b), err
}
