// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package main

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson89aae3efDecodeGithubComDiscHighloadcup(in *jlexer.Lexer, out *Visit) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = uint(in.Uint())
		case "location":
			out.Location = uint(in.Uint())
		case "user":
			out.User = uint(in.Uint())
		case "visited_at":
			out.Visited_at = int(in.Int())
		case "mark":
			out.Mark = uint(in.Uint())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson89aae3efEncodeGithubComDiscHighloadcup(out *jwriter.Writer, in Visit) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"id\":")
	out.Uint(uint(in.Id))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"location\":")
	out.Uint(uint(in.Location))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"user\":")
	out.Uint(uint(in.User))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"visited_at\":")
	out.Int(int(in.Visited_at))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"mark\":")
	out.Uint(uint(in.Mark))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Visit) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson89aae3efEncodeGithubComDiscHighloadcup(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Visit) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson89aae3efEncodeGithubComDiscHighloadcup(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Visit) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson89aae3efDecodeGithubComDiscHighloadcup(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Visit) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson89aae3efDecodeGithubComDiscHighloadcup(l, v)
}
func easyjson89aae3efDecodeGithubComDiscHighloadcup1(in *jlexer.Lexer, out *UserVisitsFilter) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson89aae3efEncodeGithubComDiscHighloadcup1(out *jwriter.Writer, in UserVisitsFilter) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserVisitsFilter) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson89aae3efEncodeGithubComDiscHighloadcup1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserVisitsFilter) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson89aae3efEncodeGithubComDiscHighloadcup1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserVisitsFilter) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson89aae3efDecodeGithubComDiscHighloadcup1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserVisitsFilter) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson89aae3efDecodeGithubComDiscHighloadcup1(l, v)
}
func easyjson89aae3efDecodeGithubComDiscHighloadcup2(in *jlexer.Lexer, out *UserVisits) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "visits":
			if in.IsNull() {
				in.Skip()
				out.Visits = nil
			} else {
				in.Delim('[')
				if out.Visits == nil {
					if !in.IsDelim(']') {
						out.Visits = make([]UserVisit, 0, 2)
					} else {
						out.Visits = []UserVisit{}
					}
				} else {
					out.Visits = (out.Visits)[:0]
				}
				for !in.IsDelim(']') {
					var v1 UserVisit
					(v1).UnmarshalEasyJSON(in)
					out.Visits = append(out.Visits, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson89aae3efEncodeGithubComDiscHighloadcup2(out *jwriter.Writer, in UserVisits) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"visits\":")
	if in.Visits == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in.Visits {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserVisits) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson89aae3efEncodeGithubComDiscHighloadcup2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserVisits) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson89aae3efEncodeGithubComDiscHighloadcup2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserVisits) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson89aae3efDecodeGithubComDiscHighloadcup2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserVisits) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson89aae3efDecodeGithubComDiscHighloadcup2(l, v)
}
func easyjson89aae3efDecodeGithubComDiscHighloadcup3(in *jlexer.Lexer, out *UserVisit) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "mark":
			out.Mark = uint(in.Uint())
		case "visited_at":
			out.Visited_at = int(in.Int())
		case "place":
			out.Place = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson89aae3efEncodeGithubComDiscHighloadcup3(out *jwriter.Writer, in UserVisit) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"mark\":")
	out.Uint(uint(in.Mark))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"visited_at\":")
	out.Int(int(in.Visited_at))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"place\":")
	out.String(string(in.Place))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserVisit) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson89aae3efEncodeGithubComDiscHighloadcup3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserVisit) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson89aae3efEncodeGithubComDiscHighloadcup3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserVisit) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson89aae3efDecodeGithubComDiscHighloadcup3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserVisit) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson89aae3efDecodeGithubComDiscHighloadcup3(l, v)
}
func easyjson89aae3efDecodeGithubComDiscHighloadcup4(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = uint(in.Uint())
		case "email":
			out.Email = string(in.String())
		case "first_name":
			out.First_name = string(in.String())
		case "last_name":
			out.Last_name = string(in.String())
		case "gender":
			out.Gender = string(in.String())
		case "birth_date":
			out.Birth_date = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson89aae3efEncodeGithubComDiscHighloadcup4(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"id\":")
	out.Uint(uint(in.Id))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"email\":")
	out.String(string(in.Email))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"first_name\":")
	out.String(string(in.First_name))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"last_name\":")
	out.String(string(in.Last_name))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"gender\":")
	out.String(string(in.Gender))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"birth_date\":")
	out.Int(int(in.Birth_date))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson89aae3efEncodeGithubComDiscHighloadcup4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson89aae3efEncodeGithubComDiscHighloadcup4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson89aae3efDecodeGithubComDiscHighloadcup4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson89aae3efDecodeGithubComDiscHighloadcup4(l, v)
}
func easyjson89aae3efDecodeGithubComDiscHighloadcup5(in *jlexer.Lexer, out *LocationAvgFilter) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson89aae3efEncodeGithubComDiscHighloadcup5(out *jwriter.Writer, in LocationAvgFilter) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LocationAvgFilter) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson89aae3efEncodeGithubComDiscHighloadcup5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LocationAvgFilter) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson89aae3efEncodeGithubComDiscHighloadcup5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LocationAvgFilter) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson89aae3efDecodeGithubComDiscHighloadcup5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LocationAvgFilter) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson89aae3efDecodeGithubComDiscHighloadcup5(l, v)
}
func easyjson89aae3efDecodeGithubComDiscHighloadcup6(in *jlexer.Lexer, out *LocationAvg) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "avg":
			out.Avg = float64(in.Float64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson89aae3efEncodeGithubComDiscHighloadcup6(out *jwriter.Writer, in LocationAvg) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"avg\":")
	out.Float64(float64(in.Avg))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LocationAvg) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson89aae3efEncodeGithubComDiscHighloadcup6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LocationAvg) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson89aae3efEncodeGithubComDiscHighloadcup6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LocationAvg) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson89aae3efDecodeGithubComDiscHighloadcup6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LocationAvg) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson89aae3efDecodeGithubComDiscHighloadcup6(l, v)
}
func easyjson89aae3efDecodeGithubComDiscHighloadcup7(in *jlexer.Lexer, out *Location) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = uint(in.Uint())
		case "place":
			out.Place = string(in.String())
		case "country":
			out.Country = string(in.String())
		case "city":
			out.City = string(in.String())
		case "distance":
			out.Distance = uint(in.Uint())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson89aae3efEncodeGithubComDiscHighloadcup7(out *jwriter.Writer, in Location) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"id\":")
	out.Uint(uint(in.Id))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"place\":")
	out.String(string(in.Place))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"country\":")
	out.String(string(in.Country))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"city\":")
	out.String(string(in.City))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"distance\":")
	out.Uint(uint(in.Distance))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Location) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson89aae3efEncodeGithubComDiscHighloadcup7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Location) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson89aae3efEncodeGithubComDiscHighloadcup7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Location) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson89aae3efDecodeGithubComDiscHighloadcup7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Location) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson89aae3efDecodeGithubComDiscHighloadcup7(l, v)
}
