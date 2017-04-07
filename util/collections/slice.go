package collections

// Prepend is complement to builtin append.
func Prepend(c []interface{}, items ...interface{}) []interface{} {
	return append(items, c...)
}
