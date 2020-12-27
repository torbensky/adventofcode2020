package main

import (
	"fmt"
	"io"
	"strings"

	common "github.com/torbensky/adventofcode-common"
)

func main() {

	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %d\n", part2(common.OpenInputFile()))
}

func part1(reader io.Reader) int {
	sum := 0
	common.ScanLines(reader, func(line string) {
		sum += evaluateExpr(line)
	})
	return sum
}

func part2(reader io.Reader) int {
	sum := 0
	common.ScanLines(reader, func(line string) {
		sum += evaluate2(line)
	})
	return sum
}

func isNum(b byte) bool {
	return b >= '0' && b <= '9'
}

func eval(a int, op byte, b int) int {
	var result int
	if op == '+' {
		result = a + b
	} else {
		result = a * b
	}

	return result
}

func doEvaluate1(line string) (result, i int) {
	// default to addition mode
	var op byte = '+'
	var total int // running total
	for i < len(line) {
		switch line[i] {
		case ' ':
			// ignore whitespace
			i++
		case '(':
			i++
			v, newI := doEvaluate1(line[i:])
			total = eval(total, op, v)
			i += newI
		case ')':
			// return on closing parenthesis
			i++
			return total, i
		case '+', '*':
			// change op type
			op = line[i]
			i++
		default:

			// read full digit
			numEnd := i + 1
			for ; numEnd < len(line); numEnd++ {
				if !isNum(line[numEnd]) {
					break
				}
			}
			val := common.Atoi(line[i:numEnd])

			i = numEnd
			total = eval(total, op, val)
		}
	}

	return total, i
}

func evaluateExpr(line string) int {
	result, _ := doEvaluate1(line)
	return result
}

func operatorPriority(operator Token) int {
	switch operator.Kind() {
	case add:
		return 1
	case prod:
		return 0
	default:
		// nothing else is a priority that we care about
		return -1
	}
}

type expression []Token

func (s *expression) push(tkn Token) {
	*s = append(*s, tkn) // Simply append the new value to the end of the stack
}

func (s *expression) pop() Token {
	index := len(*s) - 1
	tkn := (*s)[index]
	*s = (*s)[:index]
	return tkn
}

func evaluateSub(tokens []Token) int {

	findSubRange := func(start int) int {
		count := 1
		for i := start; i < len(tokens); i++ {
			t := tokens[i]
			switch t.Kind() {
			case openParen:
				count++
			case closeParen:
				count--
			}

			if count == 0 {
				return i
			}
		}

		panic("unclosed parentheses encountered!")
	}

	// Do all the add operations first
	var remaining expression
	pos := 0
	for pos < len(tokens) {

		current := tokens[pos]
		pos++

		switch current.Kind() {
		case add:
			// handle addition immediately
			left := remaining.pop()
			right := tokens[pos]
			pos++

			// check for nested sub
			if right.Kind() == openParen {
				end := findSubRange(pos)
				result := evaluateSub(tokens[pos:end])
				pos = end + 1
				right = token{kind: number, val: &result}
			}

			result := left.MustValue() + right.MustValue()
			remaining.push(token{kind: number, val: &result})
		case openParen:
			// handle nested sub-expressions
			end := findSubRange(pos)
			result := evaluateSub(tokens[pos:end])
			pos = end + 1
			remaining.push(token{kind: number, val: &result})
		default:
			remaining.push(current)
		}
	}

	// Finish with the product operations
	result := 1
	// We go left-to-right over the remaining values
	for i := 0; i < len(remaining); i += 2 {
		result *= remaining[i].MustValue()
	}

	return result
}

func evaluate2(line string) int {
	line = strings.ReplaceAll(line, " ", "")
	lexer := newLexer(line)
	tokens, err := lexer.ReadAll()
	if err != nil {
		panic(err)
	}
	return evaluateSub(tokens)
}

type tokenKind int

const (
	add tokenKind = iota
	prod
	openParen
	closeParen
	number
)

type token struct {
	kind tokenKind
	val  *int
}

func (t token) String() string {
	switch t.kind {
	case number:
		return fmt.Sprintf("%d", t.MustValue())
	case add:
		return "+"
	case prod:
		return "*"
	case openParen:
		return "("
	case closeParen:
		return ")"
	default:
		return fmt.Sprintf("kind:%d", t.kind)
	}
}

func (t token) Kind() tokenKind {
	return t.kind
}

func (t token) Value() (int, error) {
	switch t.kind {
	case number:
		return *t.val, nil
	default:
		return -1, fmt.Errorf("cannot get value on ")
	}
}

func (t token) MustValue() int {
	val, err := t.Value()
	if err != nil {
		panic(err.Error())
	}

	return val
}

type Token interface {
	Kind() tokenKind
	Value() (int, error)
	MustValue() int
	String() string
}

type lexer struct {
	pos  int
	line string
}

func newLexer(line string) Lexer {
	return &lexer{pos: 0, line: line}
}

var endOfTokensError = fmt.Errorf("no more tokens")

func (l *lexer) ReadAll() ([]Token, error) {
	var all []Token
	for {
		current, err := l.NextToken()
		if err != nil {
			if err != endOfTokensError {
				return nil, err
			}
			break
		}
		all = append(all, current)
	}

	return all, nil
}

func (l *lexer) NextToken() (Token, error) {
	// End of expression
	if l.pos >= len(l.line) {
		return nil, endOfTokensError
	}

	switch c := l.line[l.pos]; c {
	case '+':
		l.pos++
		return token{kind: add}, nil
	case '*':
		l.pos++
		return token{kind: prod}, nil
	case '(':
		l.pos++
		return token{kind: openParen}, nil
	case ')':
		l.pos++
		return token{kind: closeParen}, nil
	case ' ':
		l.pos++
		return l.NextToken()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// number
		// read full digit
		numEnd := l.pos + 1
		for ; numEnd < len(l.line); numEnd++ {
			if !isNum(l.line[numEnd]) {
				break
			}
		}
		val := common.Atoi(l.line[l.pos:numEnd])
		l.pos = numEnd
		return token{kind: number, val: &val}, nil
	default:
		return nil, fmt.Errorf("unrecognized character type %q at position %d", c, l.pos)
	}
}

type Lexer interface {
	NextToken() (Token, error)
	ReadAll() ([]Token, error)
}
