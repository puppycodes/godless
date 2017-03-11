package godless

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleQuery
	ruleSelect
	ruleSelectKey
	ruleLimit
	ruleWhere
	ruleWhereClause
	ruleAndClause
	ruleOrClause
	rulePredicateClause
	rulePredicate
	rulePredicateValue
	rulePredicateRowKey
	rulePredicateKey
	rulePredicateLiteralValue
	rulePositiveInteger
	ruleKey
	ruleEscape
	ruleMustSpacing
	ruleSpacing
	ruleAction0
	rulePegText
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
)

var rul3s = [...]string{
	"Unknown",
	"Query",
	"Select",
	"SelectKey",
	"Limit",
	"Where",
	"WhereClause",
	"AndClause",
	"OrClause",
	"PredicateClause",
	"Predicate",
	"PredicateValue",
	"PredicateRowKey",
	"PredicateKey",
	"PredicateLiteralValue",
	"PositiveInteger",
	"Key",
	"Escape",
	"MustSpacing",
	"Spacing",
	"Action0",
	"PegText",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Printf("%v %v\n", rule, quote)
			} else {
				fmt.Printf("\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(buffer string) {
	node.print(false, buffer)
}

func (node *node32) PrettyPrint(buffer string) {
	node.print(true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type QueryParser struct {
	QueryAST

	Buffer string
	buffer []rune
	rules  [32]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *QueryParser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *QueryParser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *QueryParser
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *QueryParser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *QueryParser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.AddSelect()
		case ruleAction1:
			p.SetTableName(buffer[begin:end])
		case ruleAction2:
			p.SetLimit(buffer[begin:end])
		case ruleAction3:
			p.InitWhere()
		case ruleAction4:
			p.PushWhere()
		case ruleAction5:
			p.PopWhere()
		case ruleAction6:
			p.InitPredicate()
		case ruleAction7:
			p.SetPredicateCommand(buffer[begin:end])
		case ruleAction8:
			p.UsePredicateRowKey()
		case ruleAction9:
			p.AddPredicateKey(buffer[begin:end])
		case ruleAction10:
			p.AddPredicateLiteral(buffer[begin:end])

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *QueryParser) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Query <- <(Spacing Select Action0 !.)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[ruleSpacing]() {
					goto l0
				}
				{
					position2 := position
					if buffer[position] != rune('s') {
						goto l0
					}
					position++
					if buffer[position] != rune('e') {
						goto l0
					}
					position++
					if buffer[position] != rune('l') {
						goto l0
					}
					position++
					if buffer[position] != rune('e') {
						goto l0
					}
					position++
					if buffer[position] != rune('c') {
						goto l0
					}
					position++
					if buffer[position] != rune('t') {
						goto l0
					}
					position++
					if !_rules[ruleMustSpacing]() {
						goto l0
					}
					{
						position3 := position
						{
							position4 := position
							if !_rules[ruleKey]() {
								goto l0
							}
							add(rulePegText, position4)
						}
						{
							add(ruleAction1, position)
						}
						add(ruleSelectKey, position3)
					}
					{
						position6, tokenIndex6 := position, tokenIndex
						if !_rules[ruleMustSpacing]() {
							goto l6
						}
						{
							position8 := position
							if buffer[position] != rune('w') {
								goto l6
							}
							position++
							if buffer[position] != rune('h') {
								goto l6
							}
							position++
							if buffer[position] != rune('e') {
								goto l6
							}
							position++
							if buffer[position] != rune('r') {
								goto l6
							}
							position++
							if buffer[position] != rune('e') {
								goto l6
							}
							position++
							{
								add(ruleAction3, position)
							}
							if !_rules[ruleMustSpacing]() {
								goto l6
							}
							if !_rules[ruleWhereClause]() {
								goto l6
							}
							add(ruleWhere, position8)
						}
						goto l7
					l6:
						position, tokenIndex = position6, tokenIndex6
					}
				l7:
					{
						position10, tokenIndex10 := position, tokenIndex
						if !_rules[ruleMustSpacing]() {
							goto l10
						}
						{
							position12 := position
							if buffer[position] != rune('l') {
								goto l10
							}
							position++
							if buffer[position] != rune('i') {
								goto l10
							}
							position++
							if buffer[position] != rune('m') {
								goto l10
							}
							position++
							if buffer[position] != rune('i') {
								goto l10
							}
							position++
							if buffer[position] != rune('t') {
								goto l10
							}
							position++
							if !_rules[ruleMustSpacing]() {
								goto l10
							}
							{
								position13 := position
								{
									position14 := position
									if c := buffer[position]; c < rune('1') || c > rune('9') {
										goto l10
									}
									position++
								l15:
									{
										position16, tokenIndex16 := position, tokenIndex
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l16
										}
										position++
										goto l15
									l16:
										position, tokenIndex = position16, tokenIndex16
									}
									add(rulePositiveInteger, position14)
								}
								add(rulePegText, position13)
							}
							{
								add(ruleAction2, position)
							}
							add(ruleLimit, position12)
						}
						goto l11
					l10:
						position, tokenIndex = position10, tokenIndex10
					}
				l11:
					add(ruleSelect, position2)
				}
				{
					add(ruleAction0, position)
				}
				{
					position19, tokenIndex19 := position, tokenIndex
					if !matchDot() {
						goto l19
					}
					goto l0
				l19:
					position, tokenIndex = position19, tokenIndex19
				}
				add(ruleQuery, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Select <- <('s' 'e' 'l' 'e' 'c' 't' MustSpacing SelectKey (MustSpacing Where)? (MustSpacing Limit)?)> */
		nil,
		/* 2 SelectKey <- <(<Key> Action1)> */
		nil,
		/* 3 Limit <- <('l' 'i' 'm' 'i' 't' MustSpacing <PositiveInteger> Action2)> */
		nil,
		/* 4 Where <- <('w' 'h' 'e' 'r' 'e' Action3 MustSpacing WhereClause)> */
		nil,
		/* 5 WhereClause <- <(Action4 ((&('s') PredicateClause) | (&('o') OrClause) | (&('a') AndClause)) Action5)> */
		func() bool {
			position24, tokenIndex24 := position, tokenIndex
			{
				position25 := position
				{
					add(ruleAction4, position)
				}
				{
					switch buffer[position] {
					case 's':
						{
							position28 := position
							{
								add(ruleAction6, position)
							}
							{
								position30 := position
								{
									position31 := position
									{
										position32, tokenIndex32 := position, tokenIndex
										if buffer[position] != rune('s') {
											goto l33
										}
										position++
										if buffer[position] != rune('t') {
											goto l33
										}
										position++
										if buffer[position] != rune('r') {
											goto l33
										}
										position++
										if buffer[position] != rune('_') {
											goto l33
										}
										position++
										if buffer[position] != rune('e') {
											goto l33
										}
										position++
										if buffer[position] != rune('q') {
											goto l33
										}
										position++
										goto l32
									l33:
										position, tokenIndex = position32, tokenIndex32
										if buffer[position] != rune('s') {
											goto l24
										}
										position++
										if buffer[position] != rune('t') {
											goto l24
										}
										position++
										if buffer[position] != rune('r') {
											goto l24
										}
										position++
										if buffer[position] != rune('_') {
											goto l24
										}
										position++
										if buffer[position] != rune('n') {
											goto l24
										}
										position++
										if buffer[position] != rune('e') {
											goto l24
										}
										position++
										if buffer[position] != rune('q') {
											goto l24
										}
										position++
									}
								l32:
									add(rulePegText, position31)
								}
								{
									add(ruleAction7, position)
								}
								add(rulePredicate, position30)
							}
							if !_rules[ruleSpacing]() {
								goto l24
							}
							if buffer[position] != rune('(') {
								goto l24
							}
							position++
							if !_rules[ruleSpacing]() {
								goto l24
							}
							if !_rules[rulePredicateValue]() {
								goto l24
							}
						l35:
							{
								position36, tokenIndex36 := position, tokenIndex
								if buffer[position] != rune(',') {
									goto l36
								}
								position++
								if !_rules[ruleSpacing]() {
									goto l36
								}
								if !_rules[rulePredicateValue]() {
									goto l36
								}
								if !_rules[ruleSpacing]() {
									goto l36
								}
								goto l35
							l36:
								position, tokenIndex = position36, tokenIndex36
							}
							if buffer[position] != rune(')') {
								goto l24
							}
							position++
							add(rulePredicateClause, position28)
						}
						break
					case 'o':
						{
							position37 := position
							if buffer[position] != rune('o') {
								goto l24
							}
							position++
							if buffer[position] != rune('r') {
								goto l24
							}
							position++
							if !_rules[ruleSpacing]() {
								goto l24
							}
							if buffer[position] != rune('(') {
								goto l24
							}
							position++
							if !_rules[ruleSpacing]() {
								goto l24
							}
							if !_rules[ruleWhereClause]() {
								goto l24
							}
							if !_rules[ruleSpacing]() {
								goto l24
							}
						l38:
							{
								position39, tokenIndex39 := position, tokenIndex
								if buffer[position] != rune(',') {
									goto l39
								}
								position++
								if !_rules[ruleSpacing]() {
									goto l39
								}
								if !_rules[ruleWhereClause]() {
									goto l39
								}
								if !_rules[ruleSpacing]() {
									goto l39
								}
								goto l38
							l39:
								position, tokenIndex = position39, tokenIndex39
							}
							if buffer[position] != rune(')') {
								goto l24
							}
							position++
							add(ruleOrClause, position37)
						}
						break
					default:
						{
							position40 := position
							if buffer[position] != rune('a') {
								goto l24
							}
							position++
							if buffer[position] != rune('n') {
								goto l24
							}
							position++
							if buffer[position] != rune('d') {
								goto l24
							}
							position++
							if !_rules[ruleSpacing]() {
								goto l24
							}
							if buffer[position] != rune('(') {
								goto l24
							}
							position++
							if !_rules[ruleSpacing]() {
								goto l24
							}
							if !_rules[ruleWhereClause]() {
								goto l24
							}
							if !_rules[ruleSpacing]() {
								goto l24
							}
						l41:
							{
								position42, tokenIndex42 := position, tokenIndex
								if buffer[position] != rune(',') {
									goto l42
								}
								position++
								if !_rules[ruleSpacing]() {
									goto l42
								}
								if !_rules[ruleWhereClause]() {
									goto l42
								}
								if !_rules[ruleSpacing]() {
									goto l42
								}
								goto l41
							l42:
								position, tokenIndex = position42, tokenIndex42
							}
							if buffer[position] != rune(')') {
								goto l24
							}
							position++
							add(ruleAndClause, position40)
						}
						break
					}
				}

				{
					add(ruleAction5, position)
				}
				add(ruleWhereClause, position25)
			}
			return true
		l24:
			position, tokenIndex = position24, tokenIndex24
			return false
		},
		/* 6 AndClause <- <('a' 'n' 'd' Spacing '(' Spacing WhereClause Spacing (',' Spacing WhereClause Spacing)* ')')> */
		nil,
		/* 7 OrClause <- <('o' 'r' Spacing '(' Spacing WhereClause Spacing (',' Spacing WhereClause Spacing)* ')')> */
		nil,
		/* 8 PredicateClause <- <(Action6 Predicate Spacing '(' Spacing PredicateValue (',' Spacing PredicateValue Spacing)* ')')> */
		nil,
		/* 9 Predicate <- <(<(('s' 't' 'r' '_' 'e' 'q') / ('s' 't' 'r' '_' 'n' 'e' 'q'))> Action7)> */
		nil,
		/* 10 PredicateValue <- <((&('\'') PredicateLiteralValue) | (&('@') PredicateRowKey) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '\\' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') PredicateKey))> */
		func() bool {
			position48, tokenIndex48 := position, tokenIndex
			{
				position49 := position
				{
					switch buffer[position] {
					case '\'':
						{
							position51 := position
							if buffer[position] != rune('\'') {
								goto l48
							}
							position++
							{
								position52 := position
							l53:
								{
									position54, tokenIndex54 := position, tokenIndex
									{
										position55, tokenIndex55 := position, tokenIndex
										if buffer[position] != rune('\'') {
											goto l55
										}
										position++
										goto l54
									l55:
										position, tokenIndex = position55, tokenIndex55
									}
									if !matchDot() {
										goto l54
									}
									goto l53
								l54:
									position, tokenIndex = position54, tokenIndex54
								}
								add(rulePegText, position52)
							}
							if buffer[position] != rune('\'') {
								goto l48
							}
							position++
							{
								add(ruleAction10, position)
							}
							add(rulePredicateLiteralValue, position51)
						}
						break
					case '@':
						{
							position57 := position
							if buffer[position] != rune('@') {
								goto l48
							}
							position++
							if buffer[position] != rune('k') {
								goto l48
							}
							position++
							if buffer[position] != rune('e') {
								goto l48
							}
							position++
							if buffer[position] != rune('y') {
								goto l48
							}
							position++
							{
								add(ruleAction8, position)
							}
							add(rulePredicateRowKey, position57)
						}
						break
					default:
						{
							position59 := position
							{
								position60 := position
								if !_rules[ruleKey]() {
									goto l48
								}
								add(rulePegText, position60)
							}
							{
								add(ruleAction9, position)
							}
							add(rulePredicateKey, position59)
						}
						break
					}
				}

				add(rulePredicateValue, position49)
			}
			return true
		l48:
			position, tokenIndex = position48, tokenIndex48
			return false
		},
		/* 11 PredicateRowKey <- <('@' 'k' 'e' 'y' Action8)> */
		nil,
		/* 12 PredicateKey <- <(<Key> Action9)> */
		nil,
		/* 13 PredicateLiteralValue <- <('\'' <(!'\'' .)*> '\'' Action10)> */
		nil,
		/* 14 PositiveInteger <- <([1-9] [0-9]*)> */
		nil,
		/* 15 Key <- <(Escape / ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))+> */
		func() bool {
			position66, tokenIndex66 := position, tokenIndex
			{
				position67 := position
				{
					position70, tokenIndex70 := position, tokenIndex
					{
						position72 := position
						if buffer[position] != rune('\\') {
							goto l71
						}
						position++
						{
							switch buffer[position] {
							case 'v':
								if buffer[position] != rune('v') {
									goto l71
								}
								position++
								break
							case 't':
								if buffer[position] != rune('t') {
									goto l71
								}
								position++
								break
							case 'r':
								if buffer[position] != rune('r') {
									goto l71
								}
								position++
								break
							case 'n':
								if buffer[position] != rune('n') {
									goto l71
								}
								position++
								break
							case 'f':
								if buffer[position] != rune('f') {
									goto l71
								}
								position++
								break
							case 'b':
								if buffer[position] != rune('b') {
									goto l71
								}
								position++
								break
							case 'a':
								if buffer[position] != rune('a') {
									goto l71
								}
								position++
								break
							case '\\':
								if buffer[position] != rune('\\') {
									goto l71
								}
								position++
								break
							case '?':
								if buffer[position] != rune('?') {
									goto l71
								}
								position++
								break
							case '"':
								if buffer[position] != rune('"') {
									goto l71
								}
								position++
								break
							default:
								if buffer[position] != rune('\'') {
									goto l71
								}
								position++
								break
							}
						}

						add(ruleEscape, position72)
					}
					goto l70
				l71:
					position, tokenIndex = position70, tokenIndex70
					{
						switch buffer[position] {
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l66
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l66
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l66
							}
							position++
							break
						}
					}

				}
			l70:
			l68:
				{
					position69, tokenIndex69 := position, tokenIndex
					{
						position75, tokenIndex75 := position, tokenIndex
						{
							position77 := position
							if buffer[position] != rune('\\') {
								goto l76
							}
							position++
							{
								switch buffer[position] {
								case 'v':
									if buffer[position] != rune('v') {
										goto l76
									}
									position++
									break
								case 't':
									if buffer[position] != rune('t') {
										goto l76
									}
									position++
									break
								case 'r':
									if buffer[position] != rune('r') {
										goto l76
									}
									position++
									break
								case 'n':
									if buffer[position] != rune('n') {
										goto l76
									}
									position++
									break
								case 'f':
									if buffer[position] != rune('f') {
										goto l76
									}
									position++
									break
								case 'b':
									if buffer[position] != rune('b') {
										goto l76
									}
									position++
									break
								case 'a':
									if buffer[position] != rune('a') {
										goto l76
									}
									position++
									break
								case '\\':
									if buffer[position] != rune('\\') {
										goto l76
									}
									position++
									break
								case '?':
									if buffer[position] != rune('?') {
										goto l76
									}
									position++
									break
								case '"':
									if buffer[position] != rune('"') {
										goto l76
									}
									position++
									break
								default:
									if buffer[position] != rune('\'') {
										goto l76
									}
									position++
									break
								}
							}

							add(ruleEscape, position77)
						}
						goto l75
					l76:
						position, tokenIndex = position75, tokenIndex75
						{
							switch buffer[position] {
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l69
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l69
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l69
								}
								position++
								break
							}
						}

					}
				l75:
					goto l68
				l69:
					position, tokenIndex = position69, tokenIndex69
				}
				add(ruleKey, position67)
			}
			return true
		l66:
			position, tokenIndex = position66, tokenIndex66
			return false
		},
		/* 16 Escape <- <('\\' ((&('v') 'v') | (&('t') 't') | (&('r') 'r') | (&('n') 'n') | (&('f') 'f') | (&('b') 'b') | (&('a') 'a') | (&('\\') '\\') | (&('?') '?') | (&('"') '"') | (&('\'') '\'')))> */
		nil,
		/* 17 MustSpacing <- <((&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))+> */
		func() bool {
			position81, tokenIndex81 := position, tokenIndex
			{
				position82 := position
				{
					switch buffer[position] {
					case '\n':
						if buffer[position] != rune('\n') {
							goto l81
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l81
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l81
						}
						position++
						break
					}
				}

			l83:
				{
					position84, tokenIndex84 := position, tokenIndex
					{
						switch buffer[position] {
						case '\n':
							if buffer[position] != rune('\n') {
								goto l84
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l84
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l84
							}
							position++
							break
						}
					}

					goto l83
				l84:
					position, tokenIndex = position84, tokenIndex84
				}
				add(ruleMustSpacing, position82)
			}
			return true
		l81:
			position, tokenIndex = position81, tokenIndex81
			return false
		},
		/* 18 Spacing <- <((&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position88 := position
			l89:
				{
					position90, tokenIndex90 := position, tokenIndex
					{
						switch buffer[position] {
						case '\n':
							if buffer[position] != rune('\n') {
								goto l90
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l90
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l90
							}
							position++
							break
						}
					}

					goto l89
				l90:
					position, tokenIndex = position90, tokenIndex90
				}
				add(ruleSpacing, position88)
			}
			return true
		},
		/* 20 Action0 <- <{ p.AddSelect() }> */
		nil,
		nil,
		/* 22 Action1 <- <{ p.SetTableName(buffer[begin:end]) }> */
		nil,
		/* 23 Action2 <- <{ p.SetLimit(buffer[begin:end])}> */
		nil,
		/* 24 Action3 <- <{ p.InitWhere() }> */
		nil,
		/* 25 Action4 <- <{ p.PushWhere() }> */
		nil,
		/* 26 Action5 <- <{ p.PopWhere() }> */
		nil,
		/* 27 Action6 <- <{ p.InitPredicate() }> */
		nil,
		/* 28 Action7 <- <{ p.SetPredicateCommand(buffer[begin:end]) }> */
		nil,
		/* 29 Action8 <- <{ p.UsePredicateRowKey() }> */
		nil,
		/* 30 Action9 <- <{ p.AddPredicateKey(buffer[begin:end]) }> */
		nil,
		/* 31 Action10 <- <{ p.AddPredicateLiteral(buffer[begin:end])}> */
		nil,
	}
	p.rules = _rules
}
