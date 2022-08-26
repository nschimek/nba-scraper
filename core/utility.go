package core

func IdMapToArray(idMap map[string]bool) (ids []string) {
	for id, keep := range idMap {
		if keep {
			ids = append(ids, id)
		}
	}
	return
}

func ConsolidateIdMaps(idMaps ...map[string]bool) (idMap map[string]bool) {
	for _, m := range idMaps {
		for k, v := range m {
			idMap[k] = v
		}
	}

	return idMap
}

func SuppressIdMap(idMap map[string]bool, ids []string) {
	for _, id := range ids {
		if _, ok := idMap[id]; ok {
			delete(idMap, id)
		}
	}
}
