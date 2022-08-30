package core

func IdMapToArray(idMap map[string]struct{}) (ids []string) {
	for id := range idMap {
		ids = append(ids, id)
	}
	return
}

func ConsolidateIdMaps(idMaps ...map[string]struct{}) (idMap map[string]struct{}) {
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
