package pb

import "mime/multipart"

//=======================================================
// name: UploadImage
// url : /api/upload_image
// type: form-data
//=======================================================

type UploadImageRequest struct {
	File      *multipart.FileHeader `form:"file"`
	OwnerName string                `form:"owner_name"`
	Type      string                `form:"type"`
}

type UploadImageResponse struct {
	ID string `json:"id"`
}

//=======================================================
// name: UploadPresignedImage
// url : /upload_presigned_image
// type: form-data
//=======================================================

type UploadPresignedImageRequest struct {
	File *multipart.FileHeader `form:"file"`
}

type UploadPresignedImageResponse struct {
	ID string `json:"id"`
}

//=======================================================
// name: DeleteImage
// url : /api/delete_image
//=======================================================

type DeleteImageRequest struct {
	ID string `json:"id"`
}

type DeleteImageResponse struct {
}

//=======================================================
// name: ValidateImage
// url : /api/validate_image
//=======================================================

type ValidateImageRequest struct {
	ID        string `json:"id"`
	OwnerName string `json:"owner_name"`
	Type      string `json:"type"`
}

type ValidateImageResponse struct {
	IsValid bool `json:"is_valid"`
}

//=======================================================
// name: CreateImagePointer
// url : /api/create_image_pointer
//=======================================================

type CreateImagePointerRequest struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type CreateImagePointerResponse struct {
}

//=======================================================
// name: DeleteImagePointer
// url : /api/delete_image_pointer
//=======================================================

type DeleteImagePointerRequest struct {
	ID string `json:"id"`
}

type DeleteImagePointerResponse struct {
}

//=======================================================
// name: CreatePresignedURL
// url : /api/create_presigned_url
//=======================================================

type CreatePresignedURLRequest struct {
	Type      string `json:"type"`
	OwnerName string `json:"owner_name"`
	Expires   int    `json:"expires"`
	Filesize  int    `json:"filesizes"`
}

type CreatePresignedURLResponse struct {
	URL string `json:"url"`
}
