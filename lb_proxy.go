package lb_proxy

import (
	"errors"
	"fmt"
	"github.com/HanJaeseung/LoadBalancing/clusterregistry"
	"log"
	//"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/HanJaeseung/LoadBalancing/ingressregistry"

	"github.com/oschwald/geoip2-golang"
	"github.com/umahmood/haversine"
)

var (
	ErrInvalidService = errors.New("invalid service/version")
)

var ExtractPath = extractPath
var LoadBalance = loadBalance
var ExtractIP = extractIP


func extractPath(target *url.URL) (string, error) {
	fmt.Println("----Extract Path----")
	path := target.Path
	if len(path) > 1 && path[0] == '/' {
		path = path[1:]
	}
	if path == "favicon.ico" {
		return "", fmt.Errorf("Invalid path")
	}
	fmt.Println("Path : " + path)
	return path, nil
}

func extractIP(target string) (string, error) {
	fmt.Println("----Extract IP----")
	tmp := strings.Split(target, ":")
	ip, _ := tmp[0], tmp[1]
	fmt.Println("IP : " + ip)
	return ip, nil
}


func extractGeo(cip string) (string, float64, float64){
	fmt.Println("----Extract Geo----")
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP("8.8.8.8")

	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
	if len(record.Subdivisions) > 0 {
		fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
	}
	fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
	//fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	//fmt.Printf("Time zone: %v\n", record.Location.TimeZone)

	fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
	return record.Country.IsoCode, record.Location.Latitude, record.Location.Longitude
}


func calcDistance(tlat, tlon, clat, clon float64) float64 {
	fmt.Println("----Calc Distance----")
	ip := haversine.Coord{Lat: tlat, Lon: tlon}
	cluster := haversine.Coord{Lat: clat, Lon: clon}
	mi, km := haversine.Distance(ip, cluster)
	fmt.Println("Miles: ", mi, "Kilometers: ", km)
	return km
}


func distanceScore(clusters []string, tcountry string, tlat, tlon float64, creg clusterregistry.Registry) map[string]float64 {
	fmt.Println("----Distance Score----")
	score := map[string]float64{}

	var policyDistance = []float64{10.0, 100.0, 1000.0, 1000000}

	for _,cluster := range clusters {
		//ccountry,_ := creg.Country(cluster)
		//ccontinent,_ := creg.Continent(cluster)
		clat,_ := creg.Latitude(cluster)
		clon,_ := creg.Longitude(cluster)
		distance := calcDistance(tlat, tlon, clat, clon)

		score[cluster] = 100.0
		for i := range policyDistance {

			if distance >= policyDistance[i] {
				score[cluster] = score[cluster] - (100.0 / float64(len(policyDistance)))
			}
		}
	}
	return score
}

func resourceScore(clusters []string, creg clusterregistry.Registry) map[string]float64 {
	fmt.Println("----Resource Score----")
	score := map[string]float64{}
	for _, cluster := range clusters {
		cScore,_ := creg.ResourceScore(cluster)
		score[cluster] = cScore
	}
	return score
}

func scoring(clusters []string, tcountry string, tlat, tlon float64, creg clusterregistry.Registry) string {
	fmt.Println("----Scoring----")

	if len(clusters) == 1 {
		endpoint,_ := creg.IngressIP(clusters[0])
		endpoint = endpoint + ":80"
		return endpoint
	}
	//minDistance := math.MaxFloat64
	//minCluster := ""
	//for _, cluster := range clusters {
	//	clat,_ := creg.Latitude(cluster)
	//	clon,_ := creg.Longitude(cluster)
	//	distance := calcDistance(tlat, tlon, clat, clon)
	//	if distance <= minDistance {
	//		minDistance = distance
	//		minCluster = cluster
	//	}
	//}
	//endpoint,_ := creg.IngressIP(minCluster)

	dscore := distanceScore(clusters, tcountry, tlat, tlon, creg)
	rscore := resourceScore(clusters, creg)
	cluster := selectCluster(dscore, rscore)
	endpoint,_ := creg.IngressIP(cluster)
	endpoint = endpoint + ":80"
	return endpoint
}

func selectCluster(dscore map[string]float64, rscore map[string]float64) string {
	fmt.Println("----Select Cluster----")
	distancePolicyWeight := 1.0
	resourcePolicyWeight := 1.0
	maxScore := 0.0
	maxCluster := ""
	for cluster,_ := range dscore {
		sumScore := (dscore[cluster] * distancePolicyWeight) + (rscore[cluster] * resourcePolicyWeight)
		if maxScore <= sumScore {
			maxScore = sumScore
			maxCluster = cluster
		}
	}
	return maxCluster
}


func loadBalance(host, tip, network, servicePath string, reg ingressregistry.Registry, creg clusterregistry.Registry) (net.Conn, error) {
	fmt.Println("----LoadBalance----")

	endpoints, err := reg.Lookup(host, servicePath)
	if err != nil {
		return nil, err
	}
	for {
		//Lunux
		//tcountry , tlat, tlon := extractGeo(tip)
		//Window
		//tcountry, tlat, tlon := "US", 37.751, -97.822
		tcountry, tlat, tlon := "US", 37.5215, 126.97416

		endpoint := scoring(endpoints, tcountry, tlat, tlon, creg)
		//if len(endpoints) == 0 {
		//	break
		//}
		//i := rand.Int() % len(endpoints)
		//endpoint := endpoints[i]
		//k := rand.Int() % 10
		//
		//if k >= 0 && k <= 5 {
		//	endpoint = endpoints[0]
		//	i = 0
		//}else if k >= 6 && k <= 8 {
		//	endpoint = endpoints[1]
		//	i = 1
		//}else {
		//	endpoint = endpoints[2]
		//	i = 2
		//}

		conn, err := net.Dial(network, endpoint)

		if err != nil {
			reg.Failure(host, servicePath, endpoint, err)
			//endpoints = append(endpoints[:i], endpoints[i+1:]...)
			continue
		}
		return conn, nil
	}
	return nil, fmt.Errorf("No endpoint available for %s", servicePath)
}


func NewMultipleHostReverseProxy(reg ingressregistry.Registry, creg clusterregistry.Registry) http.HandlerFunc {
	fmt.Println("----NewMultipleHostReversProxy----")

	return func(w http.ResponseWriter, req *http.Request) {
		host := req.Host
		ip, _ := ExtractIP(req.RemoteAddr)
		path, err := ExtractPath(req.URL)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		(&httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = path
			},
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				Dial: func(network, addr string) (net.Conn, error) {
					addr = strings.Split(addr, ":")[0]
					return LoadBalance(host, ip,  network, addr, reg, creg)
				},
				TLSHandshakeTimeout: 10 * time.Second,
			},
		}).ServeHTTP(w, req)
	}
}
