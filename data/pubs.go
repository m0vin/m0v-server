package data

import (
	"encoding/json"
        "fmt"
        "github.com/golang/glog"
        "math/rand"
        "net/http"
        "strconv"
        "time"
	rtree "github.com/dhconnelly/rtreego"
	geojson "github.com/paulmach/go.geojson"
)

// Pubs is an RTree housing the pubs. Args are number of spatial dimensions, minimum and maximum branching factors
var Pubdeets = rtree.NewTree(2, 25, 50)
var states = []string{"up", "down"}

func LoadPubdeetsLocal() {
	stationsGeojson := GeoJSON["pubs.geojson"]
	fc, err := geojson.UnmarshalFeatureCollection(stationsGeojson)
	if err != nil {
		// Note: this will take down the GAE instance by exiting this process.
		glog.Errorf("%v \n", err)
	}
	for _, f := range fc.Features {
		Pubdeets.Insert(&Station{f})
                glog.Infof("Feature inserted into tree %v %v\n", f.Geometry, f.Properties)
	}
}

func LoadPubdeets() error {
        pbs, err := GetPubs(10)
        if err != nil {
                glog.Error("GetPubs %v \n", err)
                return err
        }
        rand.Seed(time.Now().UnixNano())
        for _,pb := range pbs {
                f := geojson.NewPointFeature([]float64{float64(pb.Longitude), float64(pb.Latitude)})
                pbc, err := GetPubConfigByHash(pb.Hash)
                if err != nil {
                        glog.Error("GetPubs %v \n", err)
                        return err
                        continue
                }
                f.Properties["name"] = pbc.Nickname
                f.Properties["notes"] = fmt.Sprintf("%s %s %f Kw", pbc.Kwpmake, pbc.Kwrmake, pbc.Kwr)
                /*f.Properties["Kwpmake"] = pbc.Kwpmake
                f.Properties["Kwr"] = pbc.Kwr
                f.Properties["Kwrmake"] = pbc.Kwrmake*/
                f.Properties["state"] = states[rand.Intn(len(states))]
                Pubdeets.DeleteWithComparator(&Station{f}, ComparePubDeets)
                Pubdeets.Insert(&Station{f})
                glog.Infof("Pubdeets inserted into tree %v %v\n", f.Geometry, f.Properties)
        }
        return nil
}

func LoadDummyPubdeets() error {
        pbs, err := GetPubDummies(50)
        if err != nil {
                glog.Error("GetPubs %v \n", err)
                return err
        }
        rand.Seed(time.Now().UnixNano())
        for _,pb := range pbs {
                f := geojson.NewPointFeature([]float64{float64(pb.Longitude), float64(pb.Latitude)})
                if err != nil {
                        glog.Error("GetPubs %v \n", err)
                        return err
                        continue
                }
                f.Properties["name"] = pb.Nickname
                f.Properties["notes"] = fmt.Sprintf("%s %s %f Kw", pb.Kwpmake, pb.Kwrmake, pb.Kwr)
                /*f.Properties["Kwpmake"] = pbc.Kwpmake
                f.Properties["Kwr"] = pbc.Kwr
                f.Properties["Kwrmake"] = pbc.Kwrmake*/
                f.Properties["state"] = states[rand.Intn(len(states))]
                Pubdeets.DeleteWithComparator(&Station{f}, ComparePubDeets)
                Pubdeets.Insert(&Station{f})
                glog.Infof("Pubdeets inserted into tree %v %v\n", f.Geometry, f.Properties)
        }
        return nil
}

func ComparePubDeets(a, b rtree.Spatial) bool {
        return a.(*Station).feature.Properties["name"] == b.(*Station).feature.Properties["name"] //&& a.(*Station).feature.Properties["notes"] == b.(*Station).feature.Properties["notes"]
}

// pubdeetsHandler reads r for a "viewport" query parameter
// and writes a GeoJSON response of the features contained in
// that viewport into w.
func PubDeetsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	vp := r.FormValue("viewport")
	rect, err := newRect(vp)
	if err != nil {
		str := fmt.Sprintf("Couldn't parse viewport: %s", err)
                glog.Error("PubDeetshandler couldn't parse viewport %v \n", err)
		http.Error(w, str, 400)
		return
	}
	zm, err := strconv.ParseInt(r.FormValue("zoom"), 10, 0)
	if err != nil {
		str := fmt.Sprintf("Couldn't parse zoom: %s", err)
                glog.Error("PubDeetshandler couldn't parse zoom %v \n", err)
		http.Error(w, str, 400)
		return
	}

        glog.Infof("Tree size %v \n", Pubdeets.Size())
        glog.Infof("Rect %v \n", rect)
	s := Pubdeets.SearchIntersect(rect)
        if len(s) == 0 {
                glog.Infof("Nothing in s %v \n", s)
                glog.Infof("Rect %v \n", rect)
        }
	fc, err := clusterStations(s, int(zm))
	if err != nil {
		str := fmt.Sprintf("Couldn't cluster results: %s", err)
                glog.Error("PubDeetshandler couldn't cluster results %v \n", err)
		http.Error(w, str, 500)
		return
	}
	err = json.NewEncoder(w).Encode(fc)
	if err != nil {
		str := fmt.Sprintf("Couldn't encode results: %s", err)
                glog.Error("PubDeetshandler couldn't encode results %v \n", err)
		http.Error(w, str, 500)
		return
	}
}
