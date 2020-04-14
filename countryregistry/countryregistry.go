package countryregistry

var country = map[string]string{
	"KR" : "Asia",
	"JP" : "Asia",
	"IR" : "Asia",
	"IL" : "Asia",
	"LA" : "Asia",
	"AF" : "Asia",
	"AM" : "Asia",
	"AI" : "North America",
	"AG" : "North America",
	"BB" : "North America",
	"VG" : "North America",
	"KY" : "North America",
	"AL" : "Europe",
	"AD" : "Europe",
	"CZ" : "Europe",
	"DK" : "Europe",
	"FO" : "Europe",
	"DZ" : "Africa",
	"AO" : "Africa",
	"BW" : "Africa",
	"BF" : "Africa",
	"CD" : "Africa",
	"EG" : "Africa",
	"AR" : "South America",
	"BO" : "South America",
	"BR" : "South America",
	"CL" : "South America",
	"GY" : "South America",
	"PY" : "South America",
	"AS" : "Australia",
	"AU" : "Australia",
	"NZ" : "Australia",
	"CK" : "Australia",
	"AQ" : "Antarctica",
}

func Continent() map[string]string {
	return country
}
