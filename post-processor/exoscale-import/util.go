package exoscaleimport

func nonEmptyStringPtr(s string) *string {
	if s != "" {
		return &s
	}

	return nil
}
