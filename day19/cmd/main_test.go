package main

import (
	"fmt"
	"strings"
	"testing"
)

func loadRules(data string, replace bool) map[int]rule {
	rules := make(map[int]rule)
	for _, rs := range strings.Split(data, "\n") {
		rs = strings.TrimSpace(rs)
		fmt.Println(replace, rs)
		if replace {
			if strings.HasPrefix(rs, "8: ") {
				fmt.Println("Replacing 8!")
				rs = "8: " + replace8(10)
			}

			if strings.HasPrefix(rs, "11: ") {
				fmt.Println("Replacing 11!")
				rs = "11: " + replace11(10)
			}
		}
		rn, r := parseRuleLine(rs)
		rules[rn] = r
	}

	return rules
}

func TestBasic(t *testing.T) {
	ruleStr := `0: 1 2
1: "a"
2: 1 3 | 3 1
3: "b"`
	rules := loadRules(ruleStr, false)

	for _, line := range []string{"aab", "aba"} {
		if !matchRule(line, rules, 0) {
			t.Errorf("rule does not match '%s'\n", line)
		}
	}

	ruleStr = `0: 4 1 5
	1: 2 3 | 3 2
	2: 4 4 | 5 5
	3: 4 5 | 5 4
	4: "a"
	5: "b"`
	rules = loadRules(ruleStr, false)
	for _, c := range []struct {
		line  string
		match bool
	}{
		{"ababbb", true},
		{"abbbab", true},
		{"bababa", false},
		{"aaabbb", false},
		{"aaaabbb", false},
	} {
		matched := matchRule(c.line, rules, 0)
		if matched != c.match {
			t.Errorf("wanted %t got %t for %q\n", c.match, matched, c.line)
		}
	}
}

func TestManyGroups2(t *testing.T) {
	rString := `0: 1
	1: 2 | 2 2 | 2 2 2
	2: "a"`

	expects := []ruleExpect{
		{"a", true},
		{"aa", true},
		{"aaa", true},
		{"aaaa", false},
	}

	testRules(t, false, rString, false, expects)
}

func TestManyGroups(t *testing.T) {
	ruleStr := `0: 1 2 3
	1: 2 | 2 2 | 2 2 2 | 2 2 2 2
	2: "a"
	3: "b"`

	expects := []ruleExpect{
		{"aab", true},
		{"aa", false},
		{"bb", false},
		{"aaab", true},
		{"aaa", false},
		{"bbb", false},
		{"aaaab", true},
		{"aaaa", false},
		{"bbbb", false},
		{"aaaaab", true},
		{"aaaaa", false},
		{"bbbbb", false},
		{"aaaaaab", false},
		{"aaaaaa", false},
		{"bbbbbb", false},
		{"aaaaaaa", false},
		{"bbbbbbb", false},
		{"aaaaaaa", false},
		{"bbbbbbb", false},
	}

	testRules(t, false, ruleStr, false, expects)
}

func TestPart1(t *testing.T) {
	ruleStr := `42: 9 14 | 10 1
	9: 14 27 | 1 26
	10: 23 14 | 28 1
	1: "a"
	11: 42 31
	5: 1 14 | 15 1
	19: 14 1 | 14 14
	12: 24 14 | 19 1
	16: 15 1 | 14 14
	31: 14 17 | 1 13
	6: 14 14 | 1 14
	2: 1 24 | 14 4
	0: 8 11
	13: 14 3 | 1 12
	15: 1 | 14
	17: 14 2 | 1 7
	23: 25 1 | 22 14
	28: 16 1
	4: 1 1
	20: 14 14 | 1 15
	3: 5 14 | 16 1
	27: 1 6 | 14 18
	14: "b"
	21: 14 1 | 1 14
	25: 1 1 | 1 14
	22: 14 14
	8: 42
	26: 14 22 | 1 20
	18: 15 15
	7: 14 5 | 1 21
	24: 14 1`
	expects := []ruleExpect{
		{"bbabbbbaabaabba", true},
		{"ababaaaaaabaaab", true},
		{"ababaaaaabbbaba", true},

		{"abbbbbabbbaaaababbaabbbbabababbbabbbbbbabaaaa", false},
		{"babbbbaabbbbbabbbbbbaabaaabaaa", false},
		{"aaabbbbbbaaaabaababaabababbabaaabbababababaaa", false},
		{"bbbbbbbaaaabbbbaaabbabaaa", false},
		{"bbbababbbbaaaaaaaabbababaaababaabab", false},
		{"baabbaaaabbaaaababbaababb", false},
		{"abbbbabbbbaaaababbbbbbaaaababb", false},
		{"aaaaabbaabaaaaababaa", false},
		{"aaaabbaaaabbaaa", false},
		{"aaaabbaabbaaaaaaabbbabbbaaabbaabaaa", false},
		{"babaaabbbaaabaababbaabababaaab", false},
		{"aabbbbbaabbbaaaaaabbbbbababaaaaabbaaabba", false},
	}
	testRules(t, false, ruleStr, false, expects)
}

func TestReplace(t *testing.T) {
	wanted := "42 "
	got := replace8(1)
	if got != wanted {
		t.Errorf("wanted %q got %q\n", wanted, got)
	}

	wanted = wanted + "| 42 42 "
	got = replace8(2)
	if got != wanted {
		t.Errorf("wanted %q got %q\n", wanted, got)
	}

	wanted = wanted + "| 42 42 42 "
	got = replace8(3)
	if got != wanted {
		t.Errorf("wanted %q got %q\n", wanted, got)
	}

	wanted = "42 31 "
	got = replace11(1)
	if got != wanted {
		t.Errorf("wanted %q got %q\n", wanted, got)
	}

	wanted = "42 31 | 42 42 31 31 "
	got = replace11(2)
	if got != wanted {
		t.Errorf("wanted %q got %q\n", wanted, got)
	}

	wanted = "42 31 | 42 42 31 31 | 42 42 42 31 31 31 "
	got = replace11(3)
	if got != wanted {
		t.Errorf("wanted %q got %q\n", wanted, got)
	}
}

func TestPart2(t *testing.T) {
	ruleStr := `42: 9 14 | 10 1
	9: 14 27 | 1 26
	10: 23 14 | 28 1
	1: "a"
	11: 42 31
	5: 1 14 | 15 1
	19: 14 1 | 14 14
	12: 24 14 | 19 1
	16: 15 1 | 14 14
	31: 14 17 | 1 13
	6: 14 14 | 1 14
	2: 1 24 | 14 4
	0: 8 11
	13: 14 3 | 1 12
	15: 1 | 14
	17: 14 2 | 1 7
	23: 25 1 | 22 14
	28: 16 1
	4: 1 1
	20: 14 14 | 1 15
	3: 5 14 | 16 1
	27: 1 6 | 14 18
	14: "b"
	21: 14 1 | 1 14
	25: 1 1 | 1 14
	22: 14 14
	8: 42
	26: 14 22 | 1 20
	18: 15 15
	7: 14 5 | 1 21
	24: 14 1`

	debug = false
	expects := []ruleExpect{
		{"bbabbbbaabaabba", true},
		{"babbbbaabbbbbabbbbbbaabaaabaaa", true},
		{"aaabbbbbbaaaabaababaabababbabaaabbababababaaa", true},
		{"bbbbbbbaaaabbbbaaabbabaaa", true},
		{"bbbababbbbaaaaaaaabbababaaababaabab", true},
		{"ababaaaaaabaaab", true},
		{"ababaaaaabbbaba", true},
		{"baabbaaaabbaaaababbaababb", true},
		{"abbbbabbbbaaaababbbbbbaaaababb", true},
		{"aaaaabbaabaaaaababaa", true},
		{"aaaabbaabbaaaaaaabbbabbbaaabbaabaaa", true},
		{"aabbbbbaabbbaaaaaabbbbbababaaaaabbaaabba", true},

		{"abbbbbabbbaaaababbaabbbbabababbbabbbbbbabaaaa", false},
		{"aaaabbaaaabbaaa", false},
		{"babaaabbbaaabaababbaabababaaab", false},
	}
	testRules(t, false, ruleStr, true, expects)
}

type ruleExpect struct {
	line  string
	match bool
}

func testRules(t *testing.T, fatal bool, rulesStr string, replace bool, expects []ruleExpect) {
	rules := loadRules(rulesStr, replace)
	fmt.Println(rules[8])
	fmt.Println(rules[11])
	for _, re := range expects {
		// _, matched := matchRule(re.line, rules, 0, 0, false)
		matched := matchRule(re.line, rules, 0)
		if matched != re.match {
			err := fmt.Sprintf("wanted %t got %t for %q\n", re.match, matched, re.line)
			if fatal {
				t.Fatal(err)
			} else {
				t.Error(err)
			}
		}
	}
}
