package constants

type DmType string

const (
	Part     DmType = "part"
	Category DmType = "category"
	Tag      DmType = "tag"
)

// 直接识别：doc/docx、ppt/pptx、md、txt
// OCR识别：jpg/jpeg、png、pdf

type DetByteType string

const (
	//Doc  DetByteType = "doc"
	//Ppt  DetByteType = "ppt"

	Docx DetByteType = "docx"
	Pptx DetByteType = "pptx"
	Xlsx DetByteType = "xlsx"
	Md   DetByteType = "md"
	Txt  DetByteType = "txt"
)

var DetByteTypes = []DetByteType{Docx, Xlsx, Pptx, Md, Txt}

type DetOcrType string

const (
	Jpg  DetOcrType = "jpg"
	Jpeg DetOcrType = "jpeg"
	Png  DetOcrType = "png"
	Pdf  DetOcrType = "pdf"
)

var DetOcrTypes = []DetOcrType{Jpg, Jpeg, Png, Pdf}

var (
	NotUploadSearchUid = []uint{12}
)

const (
	LikeNum     = "like_num"
	ScanNum     = "scan_num"
	DownloadNum = "download_num"
)

var SortNums = []string{LikeNum, ScanNum, DownloadNum}
var SortSearchNums = []string{"_score", "upload_date", LikeNum, ScanNum, DownloadNum}

const (
	ClassDocument = "class_document"
)
