package requests

type ProfileFilter struct {
	ProfileGender             *string
	ProfileGenderPreference   *string
	CandidateGender           *string
	CandidateGenderPreference *string
	ProfileAge                *int64
	CandidateAge              *int64
	ProfileMinAgePreference   *int64
	ProfileMaxAgePreference   *int64
	CandidateMinAgePreference *int64
	CandidateMaxAgePreference *int64
	ExcludeUUIDs              []string // TODO change this to IDs
	UserUUID                  *string
	MaxDistanceMeters         *int64
	ProfileLat                *float64
	ProfileLng                *float64
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

type DiscoverProfile struct {
	Zipcode          string
	Gender           string
	GenderPreference string
	Age              int64
	MinAgePref       int64
	MaxAgePref       int64
	UserUUID         string
	Name             string
	CurrentLat       float64
	CurrentLng       float64
}

type MatchesForUserResult struct {
	CandidatesMatchingPrefs []*DiscoverProfile
	AbortCode               string
}

type TrackedLike struct {
	UUID       string
	UserUUID   string
	TargetUUID string
	Liked      bool
}
