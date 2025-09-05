package request

type UpdateProfileReq struct {
	Name          string `json:"name"`
	Avatar        string `json:"avatar"`
	Gender        string `json:"gender"`
	BirthDate     string `json:"birth_date"`
	SlugName      string `json:"slug_name"`
	AboutMe       string `json:"about_me"`
	InstitutionID string `json:"institution_id"`
}
