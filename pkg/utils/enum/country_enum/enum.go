package country_enum

type Enum struct {
	ShortCurrencyCode string
	Data              string
	ProfileIDTest     string
	ProfileIDReal     string
	Domain            string
	FullCountryCode   string
	CurrencyRatio     float64
}

func None() Enum {
	return Enum{}
}

// US	Seller, vendor	$
// CA	Seller, vendor	$
// UK	Seller, vendor	£
// DE	Seller, vendor	€
// FR	Seller, vendor	€
// IT	Seller, vendor	€
// ES	Seller, vendor	€
// IN	Seller, vendor	₹
// JP	Seller, vendor	¥
// CN	Seller, vendor ¥
// AU	Seller, vendor	$
// MX	Seller, vendor	$
// AE	Seller, vendor	"د.إ"
// SA	Seller, vendor	"ر.س"
// BR	Seller, vendor	R$
// NL	Seller, vendor	€
// SG	Seller, vendor	$
// PL	Seller, vendor	zł
// SE	Seller, vendor	"kr"
// EG	Seller, vendor	£
// BE	Seller, vendor	€

func America() Enum {
	return Enum{
		ShortCurrencyCode: "$",
		Data:              "US",
		ProfileIDTest:     "4267160273121855",
		ProfileIDReal:     "2371148571106983",
		Domain:            "https://www.amazon.com",
		FullCountryCode:   "USA",
		CurrencyRatio:     100,
	}
}

func Canada() Enum {
	return Enum{
		ShortCurrencyCode: "CA$",
		Data:              "CA",
		ProfileIDTest:     "2409642068750812",
		ProfileIDReal:     "3893935798156229",
		Domain:            "https://www.amazon.ca",
		FullCountryCode:   "CAN",
		CurrencyRatio:     100,
	}
}

func UnitedKingdom() Enum {
	return Enum{
		ShortCurrencyCode: "£",
		Data:              "UK",
		ProfileIDTest:     "626034983712523",
		ProfileIDReal:     "31533957822184",
		Domain:            "https://www.amazon.co.uk",
		FullCountryCode:   "GBR",
		CurrencyRatio:     100,
	}
}

func Germany() Enum {
	return Enum{
		ShortCurrencyCode: "€",
		Data:              "DE",
		ProfileIDTest:     "526684214721623",
		ProfileIDReal:     "3649371837893097",
		Domain:            "https://www.amazon.de",
		FullCountryCode:   "DEU",
		CurrencyRatio:     100,
	}
}

func France() Enum {
	return Enum{
		ShortCurrencyCode: "€",
		Data:              "FR",
		ProfileIDTest:     "3487758683009456",
		ProfileIDReal:     "4174329355592618",
		Domain:            "https://www.amazon.fr",
		FullCountryCode:   "FRA",
		CurrencyRatio:     100,
	}
}

func Italy() Enum {
	return Enum{
		ShortCurrencyCode: "€",
		Data:              "IT",
		ProfileIDTest:     "2914147275591181",
		ProfileIDReal:     "3930800112315614",
		Domain:            "https://www.amazon.it",
		FullCountryCode:   "ITA",
		CurrencyRatio:     100,
	}
}

func Spain() Enum {
	return Enum{
		ShortCurrencyCode: "€",
		Data:              "ES",
		ProfileIDTest:     "3505205831405222",
		ProfileIDReal:     "565188787186877",
		Domain:            "https://www.amazon.es",
		FullCountryCode:   "ESP",
		CurrencyRatio:     100,
	}
}

func India() Enum {
	return Enum{
		ShortCurrencyCode: "₹",
		Data:              "IN",
		ProfileIDTest:     "4267160273121855",
		ProfileIDReal:     "",
		Domain:            "https://www.amazon.in",
		FullCountryCode:   "IND",
		CurrencyRatio:     100,
	}
}

func Japan() Enum {
	return Enum{
		ShortCurrencyCode: "¥",
		Data:              "JP",
		ProfileIDTest:     "1616025616640137",
		ProfileIDReal:     "1891729059514722",
		Domain:            "https://www.amazon.co.jp",
		FullCountryCode:   "JPN",
		CurrencyRatio:     1,
	}
}

func China() Enum {
	return Enum{
		ShortCurrencyCode: "¥",
		Data:              "CN",
		ProfileIDTest:     "4267160273121855",
		ProfileIDReal:     "",
		Domain:            "https://www.amazon.cn",
		FullCountryCode:   "CHN",
		CurrencyRatio:     100,
	}
}

func Australia() Enum {
	return Enum{
		ShortCurrencyCode: "AU$",
		Data:              "AU",
		ProfileIDTest:     "1536032599447910",
		ProfileIDReal:     "47130830578699",
		Domain:            "https://www.amazon.com.au",
		FullCountryCode:   "AUS",
		CurrencyRatio:     1,
	}
}

func Mexico() Enum {
	return Enum{
		ShortCurrencyCode: "MX$",
		Data:              "MX",
		ProfileIDTest:     "1028855068287583",
		ProfileIDReal:     "1373034676703281",
		Domain:            "https://www.amazon.com.mx",
		FullCountryCode:   "MEX",
		CurrencyRatio:     100,
	}
}

func UnitedArabEmirates() Enum {
	return Enum{
		ShortCurrencyCode: "AED",
		Data:              "AE",
		ProfileIDTest:     "1345891740082852",
		ProfileIDReal:     "2685939952086058",
		Domain:            "https://www.amazon.ae",
		FullCountryCode:   "ARE",
		CurrencyRatio:     1,
	}
}

func SaudiArabia() Enum {
	return Enum{
		ShortCurrencyCode: "SAD",
		Data:              "SA",
		ProfileIDTest:     "4267160273121855",
		ProfileIDReal:     "",
		Domain:            "https://www.amazon.sa",
		FullCountryCode:   "SAU",
		CurrencyRatio:     100,
	}
}

func Brazil() Enum {
	return Enum{
		ShortCurrencyCode: "R$",
		Data:              "BR",
		ProfileIDTest:     "4267160273121855",
		ProfileIDReal:     "",
		Domain:            "https://www.amazon.com.br",
		FullCountryCode:   "BRA",
		CurrencyRatio:     100,
	}
}

func Netherlands() Enum {
	return Enum{
		ShortCurrencyCode: "€",
		Data:              "NL",
		ProfileIDTest:     "4267160273121855",
		ProfileIDReal:     "",
		Domain:            "https://www.amazon.nl",
		FullCountryCode:   "NLD",
		CurrencyRatio:     100,
	}
}

func Singapore() Enum {
	return Enum{
		ShortCurrencyCode: "SGD",
		Data:              "SG",
		ProfileIDTest:     "1739824897597105",
		ProfileIDReal:     "1045790351039774",
		Domain:            "https://www.amazon.sg",
		FullCountryCode:   "SGP",
		CurrencyRatio:     1,
	}
}

func Poland() Enum {
	return Enum{
		ShortCurrencyCode: "zł",
		Data:              "PL",
		ProfileIDTest:     "4267160273121855",
		ProfileIDReal:     "",
		Domain:            "https://www.amazon.pl",
		FullCountryCode:   "POL",
		CurrencyRatio:     100,
	}
}

func Sweden() Enum {
	return Enum{
		ShortCurrencyCode: "kr",
		Data:              "SE",
		ProfileIDTest:     "4267160273121855",
		ProfileIDReal:     "",
		Domain:            "https://www.amazon.se",
		FullCountryCode:   "SWE",
		CurrencyRatio:     100,
	}
}

func Egypt() Enum {
	return Enum{
		ShortCurrencyCode: "£",
		Data:              "EG",
		ProfileIDTest:     "4267160273121855",
		ProfileIDReal:     "",
		Domain:            "https://www.amazon.eg",
		FullCountryCode:   "EGY",
		CurrencyRatio:     100,
	}
}

func Belgium() Enum {
	return Enum{
		ShortCurrencyCode: "€",
		Data:              "BE",
		ProfileIDTest:     "4267160273121855",
		ProfileIDReal:     "",
		Domain:            "https://www.amazon.com.be",
		FullCountryCode:   "BEL",
		CurrencyRatio:     100,
	}
}

func GetEnumByData(data string) Enum {
	list := map[string]Enum{
		America().Data:            America(),
		Canada().Data:             Canada(),
		UnitedKingdom().Data:      UnitedKingdom(),
		Germany().Data:            Germany(),
		France().Data:             France(),
		Italy().Data:              Italy(),
		Spain().Data:              Spain(),
		India().Data:              India(),
		Japan().Data:              Japan(),
		China().Data:              China(),
		Australia().Data:          Australia(),
		Mexico().Data:             Mexico(),
		UnitedArabEmirates().Data: UnitedArabEmirates(),
		SaudiArabia().Data:        SaudiArabia(),
		Brazil().Data:             Brazil(),
		Netherlands().Data:        Netherlands(),
		Singapore().Data:          Singapore(),
		Poland().Data:             Poland(),
		Sweden().Data:             Sweden(),
		Egypt().Data:              Egypt(),
		Belgium().Data:            Belgium(),
	}
	if value, found := list[data]; found {
		return value
	}
	return None()
}
