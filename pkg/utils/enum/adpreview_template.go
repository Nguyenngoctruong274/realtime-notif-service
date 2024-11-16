// nolint
package enum

type AdPreviewTemplate struct {
	TypeImage     bool
	TypeVideo     bool
	Key           string
	Name          string
	NameStringKey string
	Size          Size
	PreviewSize   Size
	PreviewUrl    string
}

type Size struct {
	Width  int
	Height int
}

func PRODUCT_LISTING_PAGE() AdPreviewTemplate {
	return AdPreviewTemplate{
		TypeImage:     true,
		Key:           "PRODUCT_LISTING_PAGE",
		Name:          "Amazon product page right",
		NameStringKey: "sd_template_amazon_product_listing_page",
		PreviewUrl:    "https://m.media-amazon.com/images/G/01/sponsored-display-where-will-my-ads-appear/Amazon_product_page_right.png",
		Size: Size{
			Width:  245,
			Height: 250,
		},
	}
}

func DETAIL_PAGE_MIDDLE() AdPreviewTemplate {
	return AdPreviewTemplate{
		TypeImage:     true,
		Key:           "DETAIL_PAGE_MIDDLE",
		Name:          "Amazon product page middle",
		NameStringKey: "sd_template_amazon_product_detail_page_middle",
		PreviewUrl:    "https://m.media-amazon.com/images/G/01/sponsored-display-where-will-my-ads-appear/Amazon_product_page_middle.png",
		Size: Size{
			Width:  650,
			Height: 130,
		},
	}
}

func MOBILE_PRODUCT_SEARCH_PAGE() AdPreviewTemplate {
	return AdPreviewTemplate{
		TypeImage:     true,
		Key:           "MOBILE_PRODUCT_SEARCH_PAGE",
		Name:          "Amazon mobile rectangle",
		NameStringKey: "sd_template_amazon_mobile_pages",
		PreviewUrl:    "https://m.media-amazon.com/images/G/01/sponsored-display-where-will-my-ads-appear/Amazon_mobile_rectangle.png",
		Size: Size{
			Width:  414,
			Height: 125,
		},
	}
}

func MOBILE_LEADERBOARD() AdPreviewTemplate {
	return AdPreviewTemplate{
		TypeImage:     true,
		Key:           "MOBILE_LEADERBOARD",
		Name:          "Mobile leaderboard",
		NameStringKey: "sd_template_mobile_leaderboard",
		PreviewUrl:    "https://m.media-amazon.com/images/G/01/sponsored-display-where-will-my-ads-appear/Mobile_leaderboard.png",
		Size: Size{
			Width:  320,
			Height: 50,
		},
	}
}

func WIDE_SKYSCRAPER() AdPreviewTemplate {
	return AdPreviewTemplate{
		TypeImage:     true,
		Key:           "WIDE_SKYSCRAPER",
		Name:          "Wide skyscraper",
		NameStringKey: "sd_template_wide_skyscraper",
		PreviewUrl:    "https://m.media-amazon.com/images/G/01/sponsored-display-where-will-my-ads-appear/Wide_skyscraper.png",
		Size: Size{
			Width:  160,
			Height: 600,
		},
	}
}

func MEDIUM_RECTANGLE() AdPreviewTemplate {
	return AdPreviewTemplate{
		TypeImage:     true,
		TypeVideo:     true,
		Key:           "MEDIUM_RECTANGLE",
		Name:          "Medium rectangle",
		NameStringKey: "sd_template_medium_rectangle",
		PreviewUrl:    "https://m.media-amazon.com/images/G/01/sponsored-display-where-will-my-ads-appear/Medium_rectangle.png",
		Size: Size{
			Width:  300,
			Height: 250,
		},
	}
}

func LEADER_BOARD() AdPreviewTemplate {
	return AdPreviewTemplate{
		TypeImage:     true,
		Key:           "LEADER_BOARD",
		Name:          "Leaderboard",
		NameStringKey: "sd_template_leaderboard",
		PreviewUrl:    "https://m.media-amazon.com/images/G/01/sponsored-display-where-will-my-ads-appear/Leaderboard.png",
		Size: Size{
			Width:  728,
			Height: 90,
		},
	}
}

func BILLBOARD() AdPreviewTemplate {
	return AdPreviewTemplate{
		TypeImage:     true,
		TypeVideo:     true,
		Key:           "BILLBOARD",
		Name:          "Billboard",
		NameStringKey: "sd_template_billboard",
		PreviewUrl:    "https://m.media-amazon.com/images/G/01/sponsored-display-where-will-my-ads-appear/Billboard.png",
		Size: Size{
			Width:  970,
			Height: 250,
		},
	}
}

func LARGE_RECTANGLE() AdPreviewTemplate {
	return AdPreviewTemplate{
		TypeImage:     true,
		Key:           "LARGE_RECTANGLE",
		Name:          "Large rectangle",
		NameStringKey: "sd_template_large_rectangle",
		PreviewUrl:    "https://m.media-amazon.com/images/G/01/sponsored-display-where-will-my-ads-appear/Large_rectangle.png",
		Size: Size{
			Width:  300,
			Height: 600,
		},
	}
}

func AMAZON_TOP_STRIPE() AdPreviewTemplate {
	return AdPreviewTemplate{
		TypeImage:     true,
		Key:           "AMAZON_TOP_STRIPE",
		Name:          "Amazon top stripe",
		NameStringKey: "sd_template_amazon_top_stripe",
		PreviewUrl:    "https://m.media-amazon.com/images/G/01/sponsored-display-where-will-my-ads-appear/Amazon_top_stripe.png",
		Size: Size{
			Width:  980,
			Height: 55,
		},
	}
}

func AMAZON_HOMEPAGE_MOBILE() AdPreviewTemplate {
	return AdPreviewTemplate{
		TypeVideo:     true,
		Key:           "AMAZON_HOMEPAGE_MOBILE",
		Name:          "Amazon mobile rectangle",
		NameStringKey: "sd_template_amazon_mobile_pages",
		PreviewUrl:    "https://m.media-amazon.com/images/G/01/sponsored-display-where-will-my-ads-appear/Amazon_homepage_mobile_video.png",
		Size: Size{
			Width:  414,
			Height: 125,
		},
		PreviewSize: Size{
			Width:  414,
			Height: 283,
		},
	}
}

func ListAdPreviewByType(value string) []AdPreviewTemplate {
	switch value {
	case "IMAGE":
		return ListAdPreviewImage()
	case "VIDEO":
		return ListAdPreviewVideo()
	default:
		return []AdPreviewTemplate{}
	}
}

func ListAdPreviewImage() []AdPreviewTemplate {
	return []AdPreviewTemplate{PRODUCT_LISTING_PAGE(), DETAIL_PAGE_MIDDLE(), MOBILE_PRODUCT_SEARCH_PAGE(),
		MOBILE_LEADERBOARD(), WIDE_SKYSCRAPER(), MEDIUM_RECTANGLE(), LEADER_BOARD(), BILLBOARD(), LARGE_RECTANGLE(),
		AMAZON_TOP_STRIPE()}
}

func ListAdPreviewVideo() []AdPreviewTemplate {
	return []AdPreviewTemplate{MEDIUM_RECTANGLE(), BILLBOARD(), AMAZON_HOMEPAGE_MOBILE()}
}
