//go:build !darwin

package platform

func markEdited(string, bool) bool { return false }
