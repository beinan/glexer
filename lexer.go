package glexer

// type Lexer struct {
// 	src  string
// 	s    *scanner.Scanner
// 	Tok  rune   //last rune
// 	Text string //the token text scanner just scanned
// }

// func InitLexer(text string) *Lexer {
// 	s := &scanner.Scanner{
// 		Mode: scanner.ScanIdents,
// 	}
// 	s.Init(strings.NewReader(text))
// 	return &Lexer{
// 		src: text,
// 		s:   s,
// 	}
// }

// func (l *Lexer) Next() bool {
// 	l.Tok = l.s.Scan()
// 	if l.Tok == scanner.EOF {
// 		return false
// 	}
// 	l.Text = l.s.TokenText()
// 	return true
// }

// func (l *Lexer) IsIdent() bool {
// 	return l.Tok == scanner.Ident
// }

// func (l *Lexer) TextAndMove() string {
// 	defer l.Next()
// 	return l.Text
// }

// func (l *Lexer) LineStr() string {
// 	line := ""
// 	for {
// 		next := l.s.Next()
// 		if next == '\n' || next == scanner.EOF {
// 			l.Next() //go to next line
// 			break
// 		}
// 		line += string(next)
// 	}
// 	return line
// }

// func (l *Lexer) Pos() Position {
// 	return Position{
// 		Offset: l.s.Offset,
// 		Line:   l.s.Line,
// 		Column: l.s.Column,
// 	}
// }

// type Result struct {
// 	Found bool
// 	Data  interface{}
// 	Pos   Position
// }

// var ResultNotFound = Result{Found: false}

// func ResultFound(data interface{}, pos Position) Result {
// 	return Result{
// 		Found: true,
// 		Data:  data,
// 		Pos:   pos,
// 	}
// }

// type Position struct {
// 	Offset int
// 	Line   int
// 	Column int
// }

// func ARune(ch rune) func(*Lexer) Result {
// 	return func(l *Lexer) Result {
// 		fmt.Println("expect:", string(ch), ch, "actual", string(l.Tok), l.Tok)
// 		if l.Tok != ch {
// 			return ResultNotFound
// 		}
// 		l.Next()                        //eat this rune and move next
// 		return ResultFound(ch, l.Pos()) //we don't really need a result
// 	}
// }

// func AnIdent(l *Lexer) Result {
// 	fmt.Println("An Ident curr token:", l.Text, scanner.TokenString(l.Tok))
// 	if l.Tok != scanner.Ident {
// 		return ResultNotFound
// 	}
// 	defer l.Next() //eat this Ident and then move next
// 	return ResultFound(l.Text, l.Pos())
// }

// func Enforce(aFunc func(*Lexer) Result, panicMsg string) func(*Lexer) Result {
// 	return func(l *Lexer) Result {
// 		result := aFunc(l)
// 		if !result.Found {
// 			panic(panicMsg)
// 		}
// 		return result
// 	}
// }

// func Many(aFunc func(*Lexer) Result) func(*Lexer) Result {
// 	return func(l *Lexer) Result {
// 		var results []Result
// 		for {
// 			aResult := aFunc(l)
// 			if !aResult.Found {
// 				break
// 			}
// 			results = append(results, aResult)
// 		}
// 		return ResultFound(results, l.Pos())
// 	}
// }

// func Try(aFunc func(*Lexer) Result) func(*Lexer) Result {
// 	return func(l *Lexer) (ret Result) {
// 		Defer func() {
// 			if err := recover(); err != nil {
// 				for i, line := range strings.Split(l.src, "\n") {
// 					fmt.Printf("%v : %v \n", i, line)
// 					if i == l.Pos().Line-1 {
// 						fmt.Printf("Sytax Parsing Error: %v ^^^^ \n", err)
// 					}
// 				}
// 				ret = ResultNotFound
// 			}
// 		}()
// 		return aFunc(l)
// 	}
// }
