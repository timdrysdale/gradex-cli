package csv

import (
	"errors"
	"io"
	"sort"
	"strings"

	"github.com/fvbommel/sortorder"
)

var (
	empty = "-"
)

type CheckSummary struct {
	Checked  string
	Original string
	What     string
	Who      string
	When     string
}

type Line struct {
	data  map[string]string
	store *Store
}

type Store struct {
	Lines          []*Line
	Header         []string
	FixedHeader    []string
	RequiredHeader []string
	HeaderString   string
}

//remove commas to avoid inadverent line breaks
func safe(text string) string {
	return strings.ReplaceAll(text, ",", "")
}

func New() *Store {
	return &Store{}
}

func (s *Store) SetFixedHeader(fixedHeader []string) {
	s.FixedHeader = fixedHeader
}

func (s *Store) SetRequiredHeader(reqdHeader []string) {
	s.RequiredHeader = reqdHeader
}

func (s *Store) Add() *Line {
	newLine := &Line{
		data:  make(map[string]string),
		store: s,
	}
	s.Lines = append(s.Lines, newLine)
	return newLine
}

// add, overwriting
func (line *Line) Add(key, value string) {

	key = safe(key)
	value = safe(value)
	line.data[key] = value

}

// add, without overwriting
func (line *Line) AddSafe(key, value string) error {

	key = safe(key)
	value = safe(value)

	if _, ok := line.data[key]; !ok {
		line.data[key] = value
		return nil
	} else {
		return errors.New("Already exists")
	}
}

func TokensToString(tokens []string) string {

	return strings.Join(tokens, ",")

}

func LinesToString(lines []string) string {

	return strings.Join(lines, "\n")

}

func (line *Line) ToString() string {

	var tokens []string

	for _, key := range line.store.Header {
		token := "-"
		if value, ok := line.data[key]; ok {
			token = value
		}
		tokens = append(tokens, token)
	}

	return TokensToString(tokens)

}

// set up headers (some fixed + naturally ordered)
func (s *Store) OrderHeader() {

	// Map all the headers in the store
	// to make a unique set (in a map)
	headerMap := make(map[string]bool)

	// Put any fixed headers in their fixed order first
	s.Header = []string{}

	for _, hdr := range s.FixedHeader {
		s.Header = append(s.Header, hdr)
		headerMap[hdr] = true
	}

	// add required headers to header map but leave to be sorted with the rest
	for _, hdr := range s.RequiredHeader {
		headerMap[hdr] = true
	}

	// identity additional unique headers by mapping
	for _, line := range s.Lines {
		for key, _ := range line.data {
			if _, ok := headerMap[key]; !ok { //not in dynamic (yet)
				headerMap[key] = true
			}
		}
	}

	// delete the fixedHeaders from map to avoid repeating them
	for _, hdr := range s.FixedHeader {
		delete(headerMap, hdr)
	}

	// sort the remaining headers
	var header []string

	for key, _ := range headerMap {
		header = append(header, key)
	}

	sort.Sort(sortorder.Natural(header))

	// add sorted headers to the Store's slice
	for _, hdr := range header {
		s.Header = append(s.Header, hdr)
	}

	s.HeaderString = TokensToString(s.Header)
}

func (s *Store) GetHeader() string {

	return s.HeaderString
}

func (s *Store) ToString() string {

	s.OrderHeader()

	lines := []string{s.GetHeader()}

	for _, line := range s.Lines {
		lines = append(lines, line.ToString())
	}

	return LinesToString(lines)
}

func (s *Store) WriteCSV(w io.Writer) (int, error) {
	return w.Write([]byte(s.ToString()))
}
