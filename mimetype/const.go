package mimetype

const (
	PrefixImage = "image"
	PrefixVideo = "video"
	PrefixAudio = "audio"
	PrefixFile  = "file"
)

var (
	// TODO galan formatlary goymaly
	Images = []string{
		JPEG, PNG, HEIF, HEIC, GIF, BMP, TIFF, SVG, ICO,
	}

	Videos = []string{
		MP4, MOV,
	}

	Audios = []string{
		MP3, WAV, OGG, FLAC, AAC, M4A, XM4A, WMA, AIFF,
	}

	Files = []string{
		TTF, OTF, WOFF, WOFF2,
		TXT, HTML, CSS, JS, JSON, XML, CSV,
		PDF, DOC, DOCX, XLS, XLSX, PPT, PPTX, ODT, ODS, ODP, RTF,
		ZIP, RAR, TAR, GZ, _7Z, ISO,
		EXE, DLL, APP,
		YAML, MARKDOWN,
	}
)

const (
	JPEG    = "image/jpeg"
	PNG     = "image/png"
	GIF     = "image/gif"
	BMP     = "image/bmp"
	TIFF    = "image/tiff"
	HEIF    = "image/heif"
	HEIFSeq = "image/heif-sequence"
	HEIC    = "image/heic"
	HEICSeq = "image/heic-sequence"
	WEBP    = "image/webp"
	ICO     = "image/x-icon"
	SVG     = "image/svg+xml"
	JP2     = "image/jp2"
	PCX     = "image/x-pcx"
	MICO    = "image/vnd.microsoft.icon"
)

const (
	MP4  = "video/mp4"
	MKV  = "video/x-matroska"
	AVI  = "video/x-msvideo"
	MOV  = "video/quicktime"
	WMV  = "video/x-ms-wmv"
	FLV  = "video/x-flv"
	WEBM = "video/webm"
	_3GP = "video/3gpp"
	OGV  = "video/ogg"
)

const (
	MP3  = "audio/mpeg"
	WAV  = "audio/x-wav"
	OGG  = "audio/ogg"
	FLAC = "audio/flac"
	AAC  = "audio/aac"
	M4A  = "audio/mp4"
	XM4A = "audio/x-m4a"
	WMA  = "audio/x-ms-wma"
	AIFF = "audio/aiff"
	OPUS = "audio/opus"
)

const (
	TXT  = "text/plain"
	HTML = "text/html"
	CSS  = "text/css"
	JS   = "application/javascript"
	JSON = "application/json"
	XML  = "application/xml"
	CSV  = "text/csv"
)

const (
	PDF  = "application/pdf"
	DOC  = "application/vnd.microsoft.portable-executable"
	DOCX = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	XLS  = "application/x-ole-storage"
	XLSX = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	PPT  = "application/vnd.ms-powerpoint"
	PPTX = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	ODT  = "application/vnd.oasis.opendocument.text"
	ODS  = "application/vnd.oasis.opendocument.spreadsheet"
	ODP  = "application/vnd.oasis.opendocument.presentation"
	RTF  = "application/rtf"
)

const (
	ZIP = "application/zip"
	RAR = "application/x-rar-compressed"
	TAR = "application/x-tar"
	GZ  = "application/gzip"
	_7Z = "application/x-7z-compressed"
	ISO = "application/x-iso9660-image"
)

const (
	EXE = "application/vnd.microsoft.portable-executable"
	DLL = "application/x-msdownload"
	APP = "application/x-apple-diskimage"
)

const (
	TTF   = "font/ttf"
	OTF   = "font/otf"
	WOFF  = "font/woff"
	WOFF2 = "font/woff2"
)

const (
	YAML     = "application/x-yaml"
	MARKDOWN = "text/markdown"
)
