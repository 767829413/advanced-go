package chrome

type Pass struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Url          string `json:"url"`
	TargetSuffix string `json:"targetSuffix"`
}
