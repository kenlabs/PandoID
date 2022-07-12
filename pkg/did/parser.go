package did

import (
	"fmt"
)

type parser struct {
	input        string
	currentIndex int
	out          *DID
	err          error
}

type parserStep func() parserStep

func (p *parser) parsePath() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	var indexIncrement int
	var next parserStep
	var percentEncoded bool

	for {
		if currentIndex == inputLength {
			next = nil
			break
		}

		char := input[currentIndex]

		if char == '/' {
			next = p.parsePath
			break
		}

		if char == '?' {
			next = p.parseQuery
			break
		}

		if char == '%' {
			if (currentIndex+2 >= inputLength) ||
				isNotHexDigit(input[currentIndex+1]) ||
				isNotHexDigit(input[currentIndex+2]) {
				return p.errorf(currentIndex, "%% is not followed by 2 hex digits")
			}

			percentEncoded = true
			indexIncrement = 3
		} else {
			percentEncoded = false
			indexIncrement = 1
		}

		if !percentEncoded && isNotValidPathChar(char) {
			return p.errorf(currentIndex, "character is not allowed in path")
		}

		currentIndex = currentIndex + indexIncrement
	}

	if currentIndex == startIndex && len(p.out.PathSegments) == 0 {
		return p.errorf(currentIndex, "size of the first path segment should not less than 1")
	}

	p.currentIndex = currentIndex
	p.out.PathSegments = append(p.out.PathSegments, input[startIndex:currentIndex])

	return next
}

func (p *parser) parseFragment() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	var indexIncrement int
	var percentEncoded bool

	for {
		if currentIndex == inputLength {
			break
		}

		char := input[currentIndex]

		if char == '%' {
			if (currentIndex+2 >= inputLength) ||
				isNotHexDigit(input[currentIndex+1]) ||
				isNotHexDigit(input[currentIndex+2]) {
				return p.errorf(currentIndex, "%% is not followed by 2 hex digits")
			}

			percentEncoded = true
			indexIncrement = 3
		} else {
			percentEncoded = false
			indexIncrement = 1
		}

		if !percentEncoded && isNotValidQueryOrFragmentChar(char) {
			return p.errorf(currentIndex, "character is not allowed in fragment: %c", char)
		}

		currentIndex = currentIndex + indexIncrement
	}

	p.currentIndex = currentIndex
	p.out.Fragment = input[startIndex:currentIndex]

	return nil
}

func (p *parser) parseQuery() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	var indexIncrement int
	var next parserStep
	var percentEncoded bool

	for {
		if currentIndex == inputLength {
			break
		}

		char := input[currentIndex]

		if char == '#' {
			next = p.parseFragment
			break
		}

		if char == '%' {
			if (currentIndex+2 >= inputLength) ||
				isNotHexDigit(input[currentIndex+1]) ||
				isNotHexDigit(input[currentIndex+2]) {
				return p.errorf(currentIndex, "%% is not followed by 2 hex digits")
			}

			percentEncoded = true
			indexIncrement = 3
		} else {
			percentEncoded = false
			indexIncrement = 1
		}

		if !percentEncoded && isNotValidQueryOrFragmentChar(char) {
			return p.errorf(currentIndex, "character is not allowed in query: %c", char)
		}

		currentIndex = currentIndex + indexIncrement
	}

	p.currentIndex = currentIndex
	p.out.Query = input[startIndex:currentIndex]

	return next
}

func (p *parser) parseParamName() parserStep {
	input := p.input
	startIndex := p.currentIndex + 1
	next := p.paramTransition()
	currentIndex := p.currentIndex

	if currentIndex == startIndex {
		p.errorf(currentIndex, "size of param name should not less than 1")
	}

	p.out.Params = append(p.out.Params, Param{Name: input[startIndex:currentIndex], Value: ""})

	return next
}

func (p *parser) parseParamValue() parserStep {
	input := p.input
	startIndex := p.currentIndex + 1
	next := p.paramTransition()
	currentIndex := p.currentIndex

	p.out.Params[len(p.out.Params)-1].Value = input[startIndex:currentIndex]

	return next
}

func (p *parser) parseID() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	var next parserStep

	for {
		if currentIndex == inputLength {
			next = nil
			break
		}

		char := input[currentIndex]
		if char == ':' {
			next = p.parseID
			break
		}
		if char == ';' {
			next = p.parseParamName
			break
		}
		if char == '/' {
			next = p.parsePath
			break
		}
		if char == '?' {
			next = p.parseQuery
			break
		}
		if char == '#' {
			next = p.parseFragment
			break
		}

		if isNotValidIDChar(char) {
			return p.errorf(currentIndex, "byte is not ALPHA or DIGIT or '.' or '-'")
		}

		currentIndex = currentIndex + 1
	}

	if currentIndex == startIndex {
		return p.errorf(currentIndex, "size of id should not less than 1")
	}

	p.currentIndex = currentIndex
	p.out.IDStrings = append(p.out.IDStrings, input[startIndex:currentIndex])

	return next
}

func (p *parser) parseMethod() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	for {
		if currentIndex == inputLength {
			return p.errorf(currentIndex, "input does not have a second `:` marking end of method name")
		}

		char := input[currentIndex]
		if char == ':' {
			if currentIndex == startIndex {
				return p.errorf(currentIndex, "method is empty")
			}

			break
		}

		if isNotDigit(char) && isNotLowercaseLetter(char) {
			return p.errorf(currentIndex, "character is not a-z OR 0-9")
		}

		currentIndex = currentIndex + 1
	}

	p.currentIndex = currentIndex
	p.out.Method = input[startIndex:currentIndex]

	return p.parseID
}

func (p *parser) parseScheme() parserStep {
	currentIndex := 3
	if p.input[:currentIndex+1] != "did:" {
		return p.errorf(currentIndex, "input does not begin with 'did:' prefix")
	}

	p.currentIndex = currentIndex
	return p.parseMethod
}

func (p *parser) checkLength() parserStep {
	inputLength := len(p.input)

	if inputLength < 7 {
		return p.errorf(inputLength, "input length is less than 7")
	}

	return p.parseScheme
}

func (p *parser) paramTransition() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1

	var indexIncrement int
	var next parserStep
	var percentEncoded bool

	for {
		if currentIndex == inputLength {
			next = nil
			break
		}

		char := input[currentIndex]

		if char == ';' {
			next = p.parseParamName
			break
		}

		if char == '=' {
			next = p.parseParamValue
			break
		}

		if char == '/' {
			next = p.parsePath
			break
		}

		if char == '?' {
			next = p.parseQuery
			break
		}

		if char == '#' {
			next = p.parseFragment
			break
		}

		if char == '%' {
			if (currentIndex+2 >= inputLength) ||
				isNotHexDigit(input[currentIndex+1]) ||
				isNotHexDigit(input[currentIndex+2]) {
				return p.errorf(currentIndex, "%% is not followed by 2 hex digits")
			}

			percentEncoded = true
			indexIncrement = 3
		} else {
			percentEncoded = false
			indexIncrement = 1
		}

		if !percentEncoded && isNotValidParamChar(char) {
			return p.errorf(currentIndex, "character is not allowed in [aram: %c", char)
		}

		currentIndex = currentIndex + indexIncrement
	}

	p.currentIndex = currentIndex

	return next
}

func (p *parser) errorf(index int, format string, args ...interface{}) parserStep {
	p.currentIndex = index
	p.err = fmt.Errorf(format, args...)
	return nil
}
