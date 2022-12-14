package tools

type sortable [][]string

func (s sortable) Len() int {
	return len(s)
}
func (s sortable) Less(i, j int) bool {
	if s[i][0] == s[j][0] {
		return s[i][1] <= s[j][1]
	}
	return s[i][0] <= s[j][0]
}
func (s sortable) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func GetKeys(kv map[string]interface{}) []string {
	i := 0
	//keys := reflect.ValueOf(pkg.SVCS).MapKeys()
	keys := make([]string, len(kv))
	for k := range kv {
		keys[i] = k
		i++
	}
	return keys
}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
