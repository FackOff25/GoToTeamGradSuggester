package prefupdater

const (
	MaxPreferenceValue = 20
	MinPreferenceValue = 0.1
	LikeUpdateValue    = 0.4
	VisitedUpdateValue = 0.1
	RefuseUpdateValue  = 0.3
)

func DefaultUpdateFunc(preference float32) float32 {
	return preference
}

func LikeUpdateFunc(preference float32) float32 {
	if preference > MaxPreferenceValue {
		return MaxPreferenceValue
	}

	difference := MaxPreferenceValue - preference

	difference *= LikeUpdateValue

	preference += difference

	return preference
}

func UnlikeUpdateFunc(preference float32) float32 {
	if preference > MaxPreferenceValue {
		return MaxPreferenceValue
	}

	difference := MaxPreferenceValue - preference

	difference /= LikeUpdateValue

	preference -= difference

	if preference < MinPreferenceValue {
		preference = MinPreferenceValue
	}

	return preference
}

func VisitedUpdateFunc(preference float32) float32 {
	if preference > MaxPreferenceValue {
		return MaxPreferenceValue
	}

	difference := MaxPreferenceValue - preference

	difference *= VisitedUpdateValue

	preference += difference

	return preference
}

func UnvisitedUpdateFunc(preference float32) float32 {
	if preference > MaxPreferenceValue {
		return MaxPreferenceValue
	}

	difference := MaxPreferenceValue - preference

	difference /= VisitedUpdateValue

	preference -= difference

	if preference < MinPreferenceValue {
		preference = MinPreferenceValue
	}

	return preference
}

func RefuseUpdateFunc(preference float32) float32 {
	if preference > MaxPreferenceValue {
		return MaxPreferenceValue
	}

	difference := MaxPreferenceValue - preference

	difference /= RefuseUpdateValue

	preference -= difference

	if preference < MinPreferenceValue {
		preference = MinPreferenceValue
	}

	return preference
}

func UnrefuseUpdateFunc(preference float32) float32 {
	if preference > MaxPreferenceValue {
		return MaxPreferenceValue
	}

	difference := MaxPreferenceValue - preference

	difference *= RefuseUpdateValue

	preference += difference

	return preference
}
