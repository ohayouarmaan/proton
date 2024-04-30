package helpers

func IsAlphaNumeric(s string) bool {
	for _, charVariable := range s {
		if (charVariable < 'a' || charVariable > 'z') && (charVariable < 'A' || charVariable > 'Z') && (charVariable != '_') {
			return false
		}
	}

	return true
}
