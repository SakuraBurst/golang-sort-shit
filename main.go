package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

type CollectionSorter interface {
	sort.Interface
	SetToWriter(w io.Writer, format string)
	SetSort(s string)
	Fields() []string
}

type SortableCollection struct {
	Headers []string
	Values  CollectionSorter
}

func (c SortableCollection) IsHeaderExist(h string) bool {
	for _, v := range c.Headers {
		if strings.EqualFold(h, v) {
			return true
		}
	}
	return false
}

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

func printCol(c SortableCollection) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	header := printHeader(c)
	fmt.Fprintf(tw, format, header["nameHeaders"]...)
	fmt.Fprintf(tw, format, header["dividers"]...)
	c.Values.SetToWriter(tw, format)
	tw.Flush() // calculate column widths and print table
}

func printHeader(c SortableCollection) map[string][]interface{} {
	hMap := map[string][]interface{}{}
	for _, v := range c.Headers {
		hMap["nameHeaders"] = append(hMap["nameHeaders"], v)
		re := regexp.MustCompile(".")
		hMap["dividers"] = append(hMap["dividers"], re.ReplaceAllLiteralString(v, "-"))
	}
	return hMap
}

//!-printTracks

type TrackSorter struct {
	items []*Track
	sorts []string
}

func (x *TrackSorter) Len() int { return len(x.items) }
func (x *TrackSorter) Less(i, j int) bool {
	sorts := []string{}
	for a := len(x.sorts) - 1; a >= 0; a-- {
		sorts = append(sorts, x.sorts[a])
	}
	for _, v := range sorts {
		switch strings.ToLower(v) {
		case "title":
			if x.items[i].Title != x.items[j].Title {
				return x.items[i].Title < x.items[j].Title
			}
		case "artist":
			if x.items[i].Artist != x.items[j].Artist {
				return x.items[i].Artist < x.items[j].Artist
			}
		case "album":
			if x.items[i].Album != x.items[j].Album {
				return x.items[i].Album < x.items[j].Album
			}
		case "year":
			if x.items[i].Year != x.items[j].Year {
				return x.items[i].Year < x.items[j].Year
			}
		case "length":
			if x.items[i].Length.Seconds() != x.items[j].Length.Seconds() {
				return x.items[i].Length.Seconds() < x.items[j].Length.Seconds()
			}
		}
	}
	return x.items[i].Artist < x.items[j].Artist
}
func (x *TrackSorter) Swap(i, j int) { x.items[i], x.items[j] = x.items[j], x.items[i] }
func (x *TrackSorter) SetSort(s string) {
	x.sorts = append(x.sorts, s)
}
func (x *TrackSorter) SetToWriter(w io.Writer, format string) {
	for _, t := range x.items {
		fmt.Fprintf(w, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
}
func (x *TrackSorter) Fields() []string {
	if x.items == nil || len(x.items) == 0 {
		return []string{}
	}
	fields := []string{}
	visFields := reflect.VisibleFields(reflect.TypeOf(*x.items[0]))
	for _, v := range visFields {
		fields = append(fields, v.Name)
	}
	return fields
}

func main() {
	col := SortableCollection{make([]string, 0), &TrackSorter{items: tracks, sorts: []string{}}}
	col.Headers = append(col.Headers, col.Values.Fields()...)
	col.Values.Fields()
	printCol(col)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if col.IsHeaderExist(s.Text()) {
			col.Values.SetSort(s.Text())
			sort.Sort(col.Values)
			printCol(col)
		} else {
			fmt.Println("there is no such header")
		}

	}
}
