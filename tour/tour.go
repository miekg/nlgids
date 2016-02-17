package tour

// Tour holds the conversion from tour ID to the Dutch name.
var Tour = map[string]string{
	"walks/koninklijke": "Van Koninklijke Huize",
	"walks/romeinen":    "Romeinen en Bankiers",
	"walks/dutch":       "Nederlandse in de City",
	"walks/custom":      "Op Maat Wandeling",

	"cycling/london": "London Tour",
	"cycling/secret": "Secret Tour",
	"cycling/custom": "Op Maat Fietstocht",

	"specials/bus":    "Afternoon Tea Bustour",
	"specials/gin":    "Gin and Cocktails",
	"specials/happen": "Happen en Stappen",
}
