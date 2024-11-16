package constants

import "time"

const (
	SystemName     = "example-api"
	PromSystem     = "example_api"
	NameSpace      = "example"
	TrackIDHeader  = "track_id"
	SystemUser     = "ExampleSystem"
	Undefined      = "undefined"
	QCBudgetPrefix = "QC_TEST"
)

type ctxRequestIDKey int

const (
	RequestIDKey ctxRequestIDKey = 0
)

// Những biến chung
const (
	ContextTimeout10Sec          = 10 * time.Second
	ContextTimeout20Sec          = 20 * time.Second
	ContextTimeout60Sec          = 60 * time.Second
	ContextTimeout               = 3 // seconds
	TimeConfig                   = 2
	DefaultLimit                 = 20
	DecimalDot                   = "."
	WarnTime                     = 3000 // ms
	HourInDay                    = 24
	LastHourInDay                = 23
	TimeStampOndDay              = 24 * 60 * 60 // 24 hours * 60 minutes * 60 seconds
	HaflItemNumber               = 2
	PartOfTime                   = 10
	MaxCountModeration           = 100
	MaxStartIndexModeration      = 100
	MaxRetryConnectKafkaTimes    = 5
	MaxRetryConnectKafkaDuration = 5
)

// Type campaign
const (
	SponsoredProductsType = "Sponsored Products"
)

const (
	TeamIDIndex      = 0
	TeamNameIndex    = 0
	CateIDIndex      = 5
	CateNameIndex    = 6
	ProfileIDIndex   = 3
	MarketPlaceIndex = 1
)

// Message cho validate input
const (
	OperatorInvalid              = "Operator Invalid.!"
	FileNameInvalid              = "FileName Invalid.!"
	AsinInvalid                  = "Asin Invalid.!"
	CreativeTypeInvalid          = "Creative Type Invalid.!"
	AssetsIDInvalid              = "AssetsID Invalid.!"
	AssetsURLInvalid             = "AssetsURL Invalid.!"
	AssetTypeInvalid             = "AssetType Invalid.!"
	AssetSubTypeListInvalid      = "AssetSubTypeList Invalid.!"
	BrandEntityIDInvalid         = "BrandEntityID Invalid.!"
	TypeAdPreviewTemplateInvalid = "TypeAdPreviewTemplate Invalid.!"
	OnlyCreativesIDFilter        = "Please provide one and only one of creativeIdFilter or adGroupIdFilter."
	AdPreviewTemplateInvalid     = "Please provide Ad Preview Template."
	LanguageInvalid              = "Language Invalid.!"
	CountMorderationInvalid      = "Please provide 'count' in the inclusive range [1, 100]"
	StartIndexMorderationInvalid = "Please provide 'startIndex' in the inclusive range [0, 100]"
)

// type SB Template Preview
const (
	SBTypeVideo             = "VIDEO"
	SBTypeVideoStore        = "VIDEO_STORE"
	SBTypeProductCollection = "PRODUCT_COLLECTION"
	SBTypeStoreSpotlight    = "STORE_SPOTLIGHT"
)

var ListSBTemplateTypeImage = []string{
	SBTypeProductCollection,
	SBTypeStoreSpotlight,
}

var ListSBTemplateTypeAvailable = []string{
	SBTypeVideo,
	SBTypeProductCollection,
	SBTypeStoreSpotlight,
}

var ListSDLanguageAvailable = []string{
	"en-US", "es-MX", "zh-CN", "es-ES",
	"it-IT", "fr-FR", "fr-CA", "de-DE",
	"ja-JP", "ko-KR", "en-GB", "en-CA",
	"hi-IN", "en-IN", "en-DE", "en-ES",
	"en-FR", "en-IT", "en-JP", "en-AE",
	"ar-AE",
}
