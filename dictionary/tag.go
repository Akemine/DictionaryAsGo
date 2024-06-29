package dictionary

// Tag is a custom type representing different tags.
type Tag string

// Constants representing different tags.
const (
	Tech  Tag = "tech"
	Philo Tag = "philo"
	Food  Tag = "food"
)

// Global list of all tags.
var allTags = []Tag{Tech, Philo, Food}

// String method returns the string representation of a Tag.
func (t Tag) String() string {
	return string(t)
}

// IsValidTag checks if a given tag is valid (exists in the global list "allTags").
func IsValidTag(tag Tag) bool {
	for _, t := range allTags {
		if t == tag {
			return true
		}
	}
	return false
}
