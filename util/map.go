package util

import (
	"encoding/json"
	"net/url"
	"sort"
	"strings"
)

/*StringAble StringAble */
type StringAble interface {
	String() string
}

/*String String */
type String string

/*String String */
func (s String) String() string {
	return string(s)
}

/*ToString ToString */
func ToString(s string) String {
	return String(s)
}

/*Map Map */
type Map map[string]interface{}

/*String transfer map to JSON string */
func (m Map) String() string {
	return string(m.ToJSON())
}

/*MapNilMake if m is nil result a nil map */
func MapNilMake(m Map) Map {
	if m == nil {
		return make(Map)
	}
	return m
}

/*MapFromMap transfer to Map from map[string]interface{} */
func MapFromMap(p map[string]interface{}) Map {
	return Map(p)
}

/*Set set interface */
func (m Map) Set(key string, v interface{}) Map {
	return m.SetPath(strings.Split(key, "."), v)
}

// SetPathWithComment is the same as SetPath, but allows you to provide comment
// information to the key, that will be reused by Marshal().
func (m Map) SetPath(keys []string, v interface{}) Map {
	subtree := m
	for _, intermediateKey := range keys[:len(keys)-1] {
		nextTree, exists := subtree[intermediateKey]
		if !exists {
			nextTree = make(Map)
			subtree[intermediateKey] = nextTree // add new element here
		}
		switch node := nextTree.(type) {
		case Map:
			subtree = node
		case []Map:
			// go to most recent element
			if len(node) == 0 {
				// create element if it does not exist
				subtree[intermediateKey] = append(node, make(Map))
			}
			subtree = node[len(node)-1]
		}
	}
	subtree[keys[len(keys)-1]] = v
	return m
}

/*SetNil set interface if key is not exist */
func (m Map) SetNil(s string, v interface{}) Map {
	if !m.Has(s) {
		m.Set(s, v)
	}
	return m
}

/*SetHas set interface if key is exist */
func (m Map) SetHas(s string, v interface{}) Map {
	if m.Has(s) {
		m.Set(s, v)
	}
	return m
}

/*SetGet set value from map if key is exist */
func (m Map) SetGet(s string, v Map) Map {
	if v.Has(s) {
		m.Set(s, v[s])
	}
	return m
}

/*Get get interface from map with out default */
func (m Map) Get(key string) interface{} {
	if key == "" {
		return m
	}
	return m.GetPath(strings.Split(key, "."))
}

/*GetD get interface from map with default */
func (m Map) GetD(s string, d interface{}) interface{} {
	if v, b := (m)[s]; b {
		return v
	}
	return d
}

/*GetMap get map from map with out default */
func (m Map) GetMap(s string) Map {
	if v, b := m.Get(s).(map[string]interface{}); b {
		return v
	}

	// if v, b := m.Get(s).(Map); b {
	// 	return v
	// }
	return nil
}

/*GetMapD get map from map with default */
func (m Map) GetMapD(s string, d Map) Map {
	if v := m.GetMap(s); v != nil {
		return v
	}
	return d
}

/*GetBool get bool from map with out default */
func (m Map) GetBool(s string) bool {
	return m.GetBoolD(s, false)
}

/*GetBoolD get bool from map with default */
func (m Map) GetBoolD(s string, b bool) bool {
	if v, b := m.Get(s).(bool); b {
		return v
	}
	return b
}

/*GetNumber get float64 from map with out default */
func (m Map) GetNumber(s string) float64 {
	return m.GetNumberD(s, 0)
}

/*GetNumberD get float64 from map with default */
func (m Map) GetNumberD(s string, i float64) float64 {
	n, b := ParseNumber(m.Get(s))
	if b {
		return n
	}
	return i
}

/*GetInt64 get int64 from map with out default */
func (m Map) GetInt64(s string) int64 {
	i, b := ParseInt(m.Get(s))
	if b {
		return i
	}
	return 0
}

/*GetString get string from map with out default */
func (m Map) GetString(s string) string {
	if v, b := m.Get(s).(string); b {
		return v
	}
	return ""
}

/*GetBytes get bytes from map with default */
func (m Map) GetBytes(s string) []byte {
	if v, b := m.Get(s).([]byte); b {
		return v
	}
	return []byte(nil)
}

/*GetStringD get string from map with default */
func (m Map) GetStringD(s string, d string) string {
	if v, b := m.Get(s).(string); b {
		return v
	}
	return d
}

/*Delete delete if exist */
func (m Map) Delete(s string) {
	delete(m, s)
}

/*Has check if key exist */
func (m Map) Has(key string) bool {
	if key == "" {
		return false
	}
	return m.HasPath(strings.Split(key, "."))

}

// HasPath returns true if the given path of keys exists, false otherwise.
func (m Map) HasPath(keys []string) bool {
	return m.GetPath(keys) != nil
}

// GetPath returns the element in the tree indicated by 'keys'.
// If keys is of length zero, the current tree is returned.
func (m Map) GetPath(keys []string) interface{} {
	if len(keys) == 0 {
		return m
	}
	subtree := m
	for _, intermediateKey := range keys[:len(keys)-1] {
		value, exists := subtree[intermediateKey]
		if !exists {
			return nil
		}
		switch node := value.(type) {
		case Map:
			subtree = node
		case []Map:
			if len(node) == 0 {
				return nil
			}
			subtree = node[len(node)-1]
		default:
			return nil // cannot navigate through other node types
		}
	}
	// branch based on final node type
	return subtree[keys[len(keys)-1]]
}

/*SortKeys 排列key */
func (m Map) SortKeys() []string {
	var keys sort.StringSlice
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	return keys
}

/*ToXML transfer map to XML */
func (m Map) ToXML() string {
	if v, e := MapToXML(m); e == nil {
		return v
	}
	return ""

}

/*ParseXML parse XML bytes to map */
func (m Map) ParseXML(b []byte) {
	m.Join(XMLToMap(b))
}

/*ToJSON transfer map to JSON */
func (m Map) ToJSON() []byte {
	v, e := json.Marshal(m)
	if e != nil {
		return []byte(nil)
	}
	return v
}

/*ParseJSON parse JSON bytes to map */
func (m Map) ParseJSON(b []byte) Map {
	tmp := Map{}
	if e := json.Unmarshal(b, &tmp); e == nil {
		m.Join(tmp)
	}
	return m
}

/*URLEncode transfer map to url encode */
func (m Map) URLEncode() string {
	var buf strings.Builder
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := (m)[k]
		keyEscaped := url.QueryEscape(k)
		if v, b := vs.(string); b {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}

	}
	return buf.String()
}

func (m Map) join(source Map, replace bool) Map {
	for k, v := range source {
		if _, b := (m)[k]; replace || !b {
			(m)[k] = v
		}
	}
	return m
}

/*ReplaceJoin insert map s to m with replace */
func (m Map) ReplaceJoin(s Map) Map {
	return m.join(s, true)
}

/*Join insert map s to m with out replace */
func (m Map) Join(s Map) Map {
	return m.join(s, false)
}

//func (m Map) SaveAs(p string, f string) {
//
//}

/*Only get map with columns */
func (m Map) Only(columns []string) Map {
	p := Map{}
	for _, v := range columns {
		p[v] = (m)[v]
	}
	return p
}

/*Clone copy a map */
func (m Map) Clone() Map {
	m0 := make(Map)
	for k, v := range m {
		m0[k] = v
	}
	return m0
}

/*URLToSHA1 make sha1 from map */
func (m Map) URLToSHA1() string {
	return signatureSHA1(m)
}

//Range range all maps
func (m Map) Range(f func(key string, value interface{}) bool) {
	for k, v := range m {
		if !f(k, v) {
			return
		}
	}
}

//Check check all input keys
//return -1 if all is exist
//return index when not found
func (m Map) Check(s ...string) int {
	if s != nil {
		for idx, v := range s {
			if !m.Has(v) {
				return idx
			}
		}
	}
	return -1
}
