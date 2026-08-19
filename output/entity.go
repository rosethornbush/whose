package output

import "github.com/rosethornbush/whose/rdap"

func entityName(entities []rdap.Entity) string {
	for _, entity := range entities {
		if !entity.HasRole("registrant") {
			continue
		}

		if name := entity.Organization(); name != "" {
			return name
		}

		if name := entity.Name(); name != "" {
			return name
		}
	}

	return ""
}
