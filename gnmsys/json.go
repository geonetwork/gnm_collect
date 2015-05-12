package gnmsys

type Json struct {
	Data map[string]interface{}
}

func (obj Json) Obj(key string) Json {
	switch v := obj.Data[key].(type) {
		case map[string]interface{}:
		return Json{v}
		default:
		return Json{map[string]interface{}}
	}
}

func (obj Json) String(key string) string {
	switch v := obj.Data[key].(type) {
		case string:
		return v
		default:
		return ""
	}
}
func (obj Json) Float(key string) float64 {
	return obj.Data[key].(float64)
}
func (obj Json) resolve(path ...string) interface {} {
	seg := path[0]
	remaining := path[1:]

	if len(path) == 1 {
		return obj.Data[seg]
	} else {
		return obj.Obj(seg).resolve(remaining...)
	}
}
func (obj Json) resolveString(path ...string) string {
	return obj.resolve(path...).(string)
}
func (obj Json) resolveFloat(path ...string) float64 {
	switch v := obj.resolve(path...).(type) {
		case float64:
		return v
		case float32:
		return float64(v)
		case int:
		return float64(v)
		case int8:
		return float64(v)
		case int16:
		return float64(v)
		case int32:
		return float64(v)
		case int64:
		return float64(v)
		case string:
		return -1
		default:
		return v.(float64)
	}
}
