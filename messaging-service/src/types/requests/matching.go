package requests

type MatchingFilter struct {
	Zipcodes            []string
	ProfileGenderIs     string
	ProfileGenderPrefIs string
	ProfileAgeIs        int64
	ProfileMinAgePrefIs int64
	ProfileMaxAgePrefIs int64
	ExcludeUUIDs        []string // TODO change this to IDs
	UserUUID            string
}

type MatchingPreferences struct {
	Zipcode          string
	Gender           string
	GenderPreference string
	Age              int64
	MinAgePref       int64
	MaxAgePref       int64
	UserUUID         string
}

// what the user sees
type Question struct {
	Text     string
	Index    int64
	Category string
	UUID     string
}

// after user has answered
type TrackedQuestion struct {
	UUID         string
	Text         string
	Index        int64
	Category     string
	UserUUID     string
	QuestionUUID string
	Liked        bool
}

// connect to matching preferences
type Profile struct {
}
