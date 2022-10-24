package core

const (
	DateRangeFormat = "2006-01-02 15:04:05 MST"
)

// This will be set by -ldflags at build time
var Version = "v0.0.0"
var exists = struct{}{}

func IdMapToArray(idMap map[string]struct{}) (ids []string) {
	for id := range idMap {
		ids = append(ids, id)
	}
	return
}

func ConsolidateIdMaps(idMaps ...map[string]struct{}) (idMap map[string]struct{}) {
	idMap = make(map[string]struct{})

	for _, m := range idMaps {
		if m != nil {
			for k, v := range m {
				idMap[k] = v
			}
		}
	}

	return idMap
}

func SuppressIdMap(idMap map[string]struct{}, ids []string) {
	for _, id := range ids {
		if _, ok := idMap[id]; ok {
			delete(idMap, id)
		}
	}
}

func IdArrayToMap(ids []string) (idMap map[string]struct{}) {
	idMap = make(map[string]struct{})

	for _, id := range ids {
		idMap[id] = exists
	}

	return
}
