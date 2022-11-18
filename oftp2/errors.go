package oftp2

import "fmt"

func NewInvalidLengthError(expected, is int) InvalidLengthError {
	return InvalidLengthError{
		expected: expected,
		is:       is,
	}
}

type InvalidLengthError struct {
	expected int
	is       int
}

func (i InvalidLengthError) Error() string {
	return fmt.Sprintf("expected the length of %d, but got %d", i.expected, i.is)
}

func NewNoCrSuffixError(is string) InvalidSuffixError {
	return NewInvalidSuffixError("carriage return", is)
}

func NewInvalidSuffixError(expected, is string) InvalidSuffixError {
	return InvalidSuffixError{
		expectedSuffix: expected,
		isPrefix:       is,
	}
}

type InvalidSuffixError struct {
	expectedSuffix string
	isPrefix       string
}

func (i InvalidSuffixError) Error() string {
	return fmt.Sprintf("does not end on %v, but on %v", i.expectedSuffix, i.isPrefix)
}

func NewInvalidPrefixError(expected, is string) InvalidPrefixError {
	return InvalidPrefixError{
		expectedPrefix: expected,
		isPrefix:       is,
	}
}

type InvalidPrefixError struct {
	expectedPrefix string
	isPrefix       string
}

func (i InvalidPrefixError) Error() string {
	return fmt.Sprintf("does not start with %v, but with %v", i.expectedPrefix, i.isPrefix)
}
