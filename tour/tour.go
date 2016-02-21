package tour

// Must be kept in sync with /tours.json of the NLgids site.
// TODO(miek): parse tours.json and populate this on startup?

const nonExists = "<niet bestaand>"

// Tour holds the conversion from tour ID to the Dutch name.
var Tour = map[string]string{
	"walks/koninklijke": "Van Koninklijke Huize",
	"walks/romeinen":    "Romeinen en Bankiers",
	"walks/dutch":       "Nederlandse in de City",
	"walks/custom":      "Wandeling Op Maat",

	"cycling/london": "London Tour",
	"cycling/secret": "Secret Tour",
	"cycling/custom": "Fietstocht Op Maat",

	"specials/bus":    "Afternoon Tea Bustour",
	"specials/gin":    "Gin and Cocktails",
	"specials/happen": "Happen en Stappen",
}

// TourOrNonExists returns either the Dutch name of the tour keyed
// by key or the string "<niet bestaand>".
func NameOrNonExists(key string) string {
	if v, ok := Tour[key]; ok {
		return v
	}
	return nonExists
}
