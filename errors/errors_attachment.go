package errors

const (
	TypeAttachmentNotFound     = "ATTACHMENT_NOT_FOUND"
	TypeAttachmentUploadFailed = "ATTACHMENT_UPLOAD_FAILED"
	TypeAttachmentSizeExceeded = "ATTACHMENT_SIZE_EXCEEDED"
	TypeAttachmentTypeInvalid  = "ATTACHMENT_TYPE_INVALID"
	TypeAttachmentDeleteFailed = "ATTACHMENT_DELETE_FAILED"
)

var (
	ErrAttachmentNotFound     = New(TypeAttachmentNotFound, "附件不存在")
	ErrAttachmentUploadFailed = New(TypeAttachmentUploadFailed, "附件上传失败")
	ErrAttachmentSizeExceeded = New(TypeAttachmentSizeExceeded, "附件大小超过限制")
	ErrAttachmentTypeInvalid  = New(TypeAttachmentTypeInvalid, "附件类型不支持")
	ErrAttachmentDeleteFailed = New(TypeAttachmentDeleteFailed, "附件删除失败")
)
