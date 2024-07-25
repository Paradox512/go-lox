package main

import (
	"errors"
	"fmt"
	"math/big"
	"os"
)

type Scanner struct {
	current int
	contents string
	tokens []Token
}

func (scanner *Scanner) AddToken(token_type TokenType, lexeme string, literal interface{}) {
	var new_token Token
	new_token.token_type = token_type
	new_token.lexeme = lexeme
	new_token.literal = literal
	scanner.tokens = append(scanner.tokens, new_token)
}

func (scanner *Scanner) Advance() byte {
	ch := scanner.contents[scanner.current]
	scanner.current++
	return ch
}

func (scanner *Scanner) Match(char byte) bool {
	if scanner.AtEnd() || scanner.contents[scanner.current] != char {
		return false
	}
	scanner.current++
	return true
}

func (scanner *Scanner) AtEnd() bool {
	return scanner.current == len(scanner.contents)
}

func (scanner *Scanner) Peek() byte {
	if scanner.AtEnd() {
		return 0
	}
	return scanner.contents[scanner.current]
}

func (scanner *Scanner) PeekNext() byte {
	if scanner.current + 1 >= len(scanner.contents) {
		return 0
	}
	return scanner.contents[scanner.current + 1]
}

func IsAlpha(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'z')
}

func (scanner *Scanner) Scan(lox_file_contents string) error {
	scanner.contents = lox_file_contents
	scanner.current = 0
	line := 1
	found_error := false
	for ; !scanner.AtEnd(); {
		start := scanner.current
		char := scanner.Advance()
		switch char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			var power, digit, ten, tmp, float_literal big.Float
			power.SetFloat64(1.0)
			ten.SetFloat64(10.0)
			float_literal.SetFloat64(float64(char - '0'))
			for {
				peek := scanner.Peek()
				if peek < '0' || peek > '9' {
					break
				}
				digit.SetFloat64(float64(peek - '0'))
				tmp.Mul(&float_literal, &ten)
				float_literal.Add(&tmp, &digit)
				scanner.Advance()
			}
			if scanner.Peek() == '.' && scanner.PeekNext() >= '0' && scanner.PeekNext() <= '9' {
				scanner.Advance()
				for {
					peek := scanner.Peek()
					if peek < '0' || peek > '9' {
						break
					}
					digit.SetFloat64(float64(peek - '0'))
					power.Quo(&power, &ten)
					tmp.Mul(&power, &digit)
					float_literal.Add(&float_literal, &tmp)
					scanner.Advance()
				}
			}
			scanner.AddToken(NUMBER, scanner.contents[start:scanner.current], float_literal)
			break
		case '(': scanner.AddToken(LEFT_PAREN, "(", nil); break;
		case ')': scanner.AddToken(RIGHT_PAREN, ")", nil); break;
		case '{': scanner.AddToken(LEFT_BRACE, "{", nil); break;
		case '}': scanner.AddToken(RIGHT_BRACE, "}", nil); break;
		case ',': scanner.AddToken(COMMA, ",", nil); break;
		case '.': scanner.AddToken(DOT, ".", nil); break;
		case '-': scanner.AddToken(MINUS, "-", nil); break;
		case '+': scanner.AddToken(PLUS, "+", nil); break;
		case ';': scanner.AddToken(SEMICOLON, ";", nil); break;
		case '*': scanner.AddToken(STAR, "*", nil); break;
		case '\n': line++; break;
		case '=':
			if scanner.Match('=') {
				scanner.AddToken(EQUAL_EQUAL, "==", nil)
			} else {
				scanner.AddToken(EQUAL, "=", nil)
			}
			break;
		case '!':
			if scanner.Match('=') {
				scanner.AddToken(BANG_EQUAL, "!=", nil)
			} else {
				scanner.AddToken(BANG, "!", nil)
			}
			break;
		case '<':
			if scanner.Match('=') {
				scanner.AddToken(LESS_EQUAL, "<=", nil)
			} else {
				scanner.AddToken(LESS, "<", nil)
			}
			break;
		case '>':
			if scanner.Match('=') {
				scanner.AddToken(GREATER_EQUAL, ">=", nil)
			} else {
				scanner.AddToken(GREATER, ">", nil)
			}
			break;
		case '/':
			if scanner.Peek() != '/' {
				scanner.AddToken(SLASH, "/", nil)
				break;
			}
			for ; scanner.Peek() != '\n' && scanner.Peek() != 0 ; {
				scanner.Advance()
			}
			break;
		case '"':
			string_literal := ""
			for {
				if scanner.AtEnd() || scanner.Peek() == '\n' {
					fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", line)
					found_error = true
					break
				}
				char := scanner.Advance()
				if char == '"'{
					scanner.AddToken(STRING, "\"" + string_literal + "\"", string_literal)
					break
				}
				string_literal += string(char)
			}
			
		case '\t':
		case ' ':
			break;
		default:
			if IsAlpha(char) {
				for {
					peek := scanner.Peek()
					if !IsAlpha(peek) { break }
					scanner.Advance()
				}
				scanner.AddToken(IDENTIFIER, scanner.contents[start:scanner.current], nil)
			} else {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, char)
				found_error = true
			}
		}
	}
	scanner.AddToken(EOF, "", nil)
	if found_error {
		return errors.New("Error scanning file contents")
	}
	return nil
}

func (scanner Scanner) StringifyTokens() string {
	ret := ""
	for _, token := range scanner.tokens {
		if len(ret) > 0 {
			ret += "\n"
		}
		ret += token.String()
	}
	return ret;
}