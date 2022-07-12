package did

func isNotDigit(char byte) bool {
	return char < '\x30' || char > '\x39'
}

func isNotHexDigit(char byte) bool {
	return isNotDigit(char) &&
		(char < '\x41' || char > '\x46') &&
		(char < '\x61' || char > '\x66')
}

func isNotUppercaseLetter(char byte) bool {
	return char < '\x41' || char > '\x5A'
}

func isNotLowercaseLetter(char byte) bool {
	return char < '\x61' || char > '\x7A'
}

func isNotAlpha(char byte) bool {
	return isNotLowercaseLetter(char) && isNotUppercaseLetter(char)
}

func isNotValidParamChar(char byte) bool {
	return isNotAlpha(char) && isNotDigit(char) &&
		char != '.' && char != '-' && char != '_' && char != ':'
}

func isNotValidIDChar(char byte) bool {
	return isNotAlpha(char) && isNotDigit(char) && char != '.' && char != '-'
}

func isNotValidPathChar(char byte) bool {
	return isNotUnreservedOrSubDelim(char) && char != ':' && char != '@'
}

func isNotUnreservedOrSubDelim(char byte) bool {
	switch char {
	case '-', '.', '_', '~', '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=':
		return false
	default:
		if isNotAlpha(char) && isNotDigit(char) {
			return true
		}
		return false
	}
}

func isNotValidQueryOrFragmentChar(char byte) bool {
	return isNotValidPathChar(char) && char != '/' && char != '?'
}
