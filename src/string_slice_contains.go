package src

// Contains returns whether the given string slice contains the given string
func Contains(haystack []string, needle string) bool {
    for _, straw := range haystack {
        if straw == needle {
            return true
        }
    }
    return false
}
