package food

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"cube/config"
	"cube/lib/context"
	"cube/lib/database"
)

type VendorSortWrapper struct {
	vs []Vendor
	by func(p, q *Vendor) bool
}

func (w VendorSortWrapper) Len() int { // 重寫 Len() 方法
	return len(w.vs)
}

func (w VendorSortWrapper) Swap(i, j int) { // 重寫 Swap() 方法
	w.vs[i], w.vs[j] = w.vs[j], w.vs[i]
}

func (w VendorSortWrapper) Less(i, j int) bool { // 重寫 Less() 方法
	return w.by(&w.vs[i], &w.vs[j])
}

type ListArgs struct {
	Limit        int64    `arg:"positional"`
	CuisineTypes []string `arg:"positional"`
}

func (h *Food) list(req *context.ChatContext, args *ListArgs) context.Response {
	cuisineTypes := []CuisineType{}
	limit := 6

	for _, t := range args.CuisineTypes {
		cuisineTypes = append(cuisineTypes, CuisineType(t))
	}

	vendors := listVendors(config.Conf.Map.Lat, config.Conf.Map.Lng, cuisineTypes)
	relations := []database.FoodTagRelation{}
	tx := h.DB.Find(&relations)
	if tx.Error != nil {
		return context.Response(tx.Error.Error())
	}

	relationMap := map[uint64]database.FoodTagRelation{}
	for _, r := range relations {
		relationMap[r.VendorID] = r
	}

	allowedVendors := []Vendor{}
	for _, v := range vendors {
		r, ok := relationMap[v.ID]
		if ok && r.FoodTagID != 1 {
			continue
		}
		allowedVendors = append(allowedVendors, v)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allowedVendors), func(i, j int) { allowedVendors[i], allowedVendors[j] = allowedVendors[j], allowedVendors[i] })

	if args.Limit > 0 {
		allowedVendors = allowedVendors[:args.Limit]
	} else if len(allowedVendors) > limit {
		allowedVendors = allowedVendors[:limit]
	}

	sort.Sort(VendorSortWrapper{allowedVendors, func(p, q *Vendor) bool {
		return q.Distance > p.Distance
	}})

	var resp string
	for i, v := range allowedVendors {
		cuisines := []string{}
		for _, c := range v.Cuisines {
			cuisines = append(cuisines, c.Name)
		}

		resp += fmt.Sprintf("%02d %v★(%v) %.2f km \"%v\" \"%v\" %v\n", i+1,
			v.Rating,
			v.ReviewNumber,
			v.Distance,
			v.Name,
			strings.Join(cuisines, " "),
			v.RedirectionURL,
		)
	}

	return context.Response(fmt.Sprintf("```%v```", resp))
}

func listVendors(lat, lng float64, cuisineTypes []CuisineType) (vs []Vendor) {
	url := "https://disco.deliveryhero.io/search/api/v1/feed"
	contentType := "application/json"

	payload := map[string]interface{}{
		"q": "",
		"location": map[string]interface{}{
			"point": map[string]interface{}{
				"longitude": lng,
				"latitude":  lat,
			},
		},
		"budgets":                 []int{},
		"config":                  "Variant17",
		"vertical_types":          []string{"restaurants"},
		"include_component_types": []string{"vendors"},
		"include_fields":          []string{"feed"},
		"language_id":             "6",
		"opening_type":            "delivery",
		"platform":                "web",
		"language_code":           "zh",
		"customer_type":           "regular",
		"limit":                   48,
		"offset":                  0,
		"dynamic_pricing":         0,
		"brand":                   "foodpanda",
		"country_code":            "tw",
		"use_free_delivery_label": false,
		"sort":                    "distance_asc",
	}

	if len(cuisineTypes) > 0 {
		targetCuisineType := []int{}
		for _, t := range cuisineTypes {
			targetCuisineType = append(targetCuisineType, cuisineType[t])
		}
		payload["cuisine"] = targetCuisineType
	}

	reqBody, _ := json.Marshal(payload)
	resp, err := http.Post(url, contentType, bytes.NewReader(reqBody))
	if err != nil {
		log.Errorf("Post: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Post: %v", resp.StatusCode)
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("ReadAll: %v", err)
		return
	}

	var nearBy NearByResp
	err = json.Unmarshal(respBody, &nearBy)
	if err != nil {
		log.Errorf("Unmarshal: %v", err)
		return
	}

	return nearBy.Feed.Items[0].Items
}
