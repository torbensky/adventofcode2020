package main

import (
	"strings"
	"testing"
)

func TestIntersect(t *testing.T) {

	for _, c := range []struct {
		a fieldSet
		b fieldSet
		l int
	}{
		{
			a: fieldSet{"foo": struct{}{}, "bar": struct{}{}},
			b: fieldSet{"bar": struct{}{}},
			l: 1,
		},
		{
			a: fieldSet{"foo": struct{}{}, "bar": struct{}{}},
			b: fieldSet{"foo": struct{}{}, "bar": struct{}{}},
			l: 2,
		},
		{
			a: fieldSet{"foo": struct{}{}, "bar": struct{}{}, "baz": struct{}{}},
			b: fieldSet{"bar": struct{}{}, "foo": struct{}{}},
			l: 2,
		},
		{
			a: fieldSet{"foo": struct{}{}, "bar": struct{}{}, "baz": struct{}{}},
			b: fieldSet{},
			l: 0,
		},
		{
			a: fieldSet{},
			b: fieldSet{},
			l: 0,
		},
	} {
		result := intersect(c.a, c.b)
		got := len(result)
		if c.l != got {
			t.Fatalf("%v | %v failed: wanted %d got %d\n", c.a, c.b, c.l, len(result))
		}
	}
}

func TestColumnRules(t *testing.T) {
	_, tickets, _, _ := loadTicketData(strings.NewReader(example2))

	cr := newColumnRules(tickets)

	want := 1
	if count := cr.countFields(0); want != count {
		t.Fatalf("wanted %d got %d\n", want, count)
	}

	want = 2
	if count := cr.countFields(1); want != count {
		t.Fatalf("wanted %d got %d\n", want, count)
	}

	want = 3
	if count := cr.countFields(2); want != count {
		t.Fatalf("wanted %d got %d\n", want, count)
	}

	col, field := cr.identifyNext()
	if field != "row" || col != 0 {
		t.Fatalf("wanted [0,row] got [%d,%s]\n", col, field)
	}

	cr.deleteField("row")

	want = 1
	if count := cr.countFields(1); want != count {
		t.Fatalf("wanted %d got %d\n", want, count)
	}

	want = 2
	if count := cr.countFields(2); want != count {
		t.Fatalf("wanted %d got %d\n", want, count)
	}

	col, field = cr.identifyNext()
	if field != "class" || col != 1 {
		t.Fatalf("wanted [1,class] got [%d,%s]\n", col, field)
	}

	cr.deleteField("class")

	want = 1
	if count := cr.countFields(2); want != count {
		t.Fatalf("wanted %d got %d\n", want, count)
	}

	col, field = cr.identifyNext()
	if field != "seat" || col != 2 {
		t.Fatalf("wanted [2,seat] got [%d,%s]\n", col, field)
	}
}

func TestFindFields(t *testing.T) {
	schema, tickets, _, _ := loadTicketData(strings.NewReader(example1))

	for _, c := range []struct {
		col    int
		fields []string
	}{
		{col: 0, fields: []string{"class", "row"}},
		{col: 1, fields: []string{"class"}},
		{col: 2, fields: []string{"seat"}},
	} {
		fields := schema.findFields(tickets[0], c.col)
		for _, want := range c.fields {
			found := false
			for _, got := range fields {
				if want == got {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("column %d is missing field %s\n", c.col, want)
			}
		}
	}
}

func TestLoadData(t *testing.T) {
	schema, tickets, yt, errors := loadTicketData(strings.NewReader(example1))

	if errors != 71 {
		t.Fatalf("wanted 71 got %d\n", errors)
	}

	for i, want := range []int{7, 1, 14} {
		got := yt.values[i]
		if got != want {
			t.Fatalf("wanted %d got %d\n", want, got)
		}
	}

	if len(schema) != 3 {
		t.Fatalf("want 3 got %d\n", len(schema))
	}

	for _, c := range []struct {
		field  string
		ranges [2]fieldRange
	}{
		{field: "class", ranges: [2]fieldRange{{1, 3}, {5, 7}}},
		{field: "row", ranges: [2]fieldRange{{6, 11}, {33, 44}}},
		{field: "seat", ranges: [2]fieldRange{{13, 40}, {45, 50}}},
	} {
		if _, ok := schema[c.field]; !ok {
			t.Fatalf("schema is missing field %s\n", c.field)
		}
		for i, r := range c.ranges {
			startMatches := r.start == schema[c.field][i].start
			endMatches := r.end == schema[c.field][i].end
			if !startMatches || !endMatches {
				t.Fatalf(
					"schema error on field %s: [%d,%d] != [%d,%d]\n",
					c.field,
					c.ranges[i].start,
					c.ranges[i].end,
					schema[c.field][i].start,
					schema[c.field][i].end,
				)
			}
		}
	}

	if len(tickets) != 1 {
		t.Fatalf("wanted 1 valid ticket, got %d\n", len(tickets))
	}
}

const example1 = `class: 1-3 or 5-7
row: 6-11 or 33-44
seat: 13-40 or 45-50

your ticket:
7,1,14

nearby tickets:
7,3,47
40,4,50
55,2,20
38,6,12`

const example2 = `class: 0-1 or 4-19
row: 0-5 or 8-19
seat: 0-13 or 16-19

your ticket:
11,12,13

nearby tickets:
3,9,18
15,1,5
5,14,9`
