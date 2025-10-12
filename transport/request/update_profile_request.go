package request

import "mime/multipart"

type UpdateProfileReq struct {
	Name          string                `form:"name"`
	Avatar        *multipart.FileHeader `form:"avatar"`
	Gender        string                `form:"gender"`
	BirthDate     string                `form:"birth_date"`
	SlugName      string                `form:"slug_name"`
	AboutMe       string                `form:"about_me"`
	InstitutionID string                `form:"institution_id"`
	Linkedin      string                `form:"linkedin"`

	AvatarDelete      *string `form:"avatar_delete"`
	AboutMeDelete     *string `form:"about_me_delete"`
	LinkedinDelete    *string `form:"linkedin_delete"`
	InstitutionDelete *string `form:"institution_delete"`
}
