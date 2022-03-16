package food

import (
	"cube/lib/context"
	"cube/lib/database"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type PollArgs struct {
	Limit int64 `arg:"positional"`
}

func (h *Food) poll(req *context.ChatContext, args *PollArgs) context.Response {
	limit := 6

	allowedVendors := []database.FoodTagRelation{}
	tx := h.DB.Find(&allowedVendors, map[string]interface{}{"food_tag_id": 1})
	if tx.Error != nil {
		return context.Response(tx.Error.Error())
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allowedVendors), func(i, j int) { allowedVendors[i], allowedVendors[j] = allowedVendors[j], allowedVendors[i] })

	if args.Limit > 0 {
		limit = int(args.Limit)
	}

	var resp string
	respCount := 0
	for _, v := range allowedVendors {
		if len(v.VendorCode) < 1 {
			log.Infof("%v: no code", v.VendorCode)
			continue
		}

		time.Sleep(200 * time.Millisecond)

		info, err := getVendorInfo(v.VendorCode)
		if err != nil {
			log.Error(err)
			continue
		}

		if !info.IsActive || !info.IsDeliveryEnabled || !info.Metadata.IsDeliveryAvailable {
			log.Infof("%v: IsActive=%v IsDeliveryEnabled=%v IsDeliveryAvailable=%v", v.VendorCode, info.IsActive, info.IsDeliveryEnabled, info.Metadata.IsDeliveryAvailable)
			continue
		}

		cuisines := []string{}
		for _, c := range info.Cuisines {
			cuisines = append(cuisines, c.Name)
		}

		resp += fmt.Sprintf("%02d %v★(%v) \"%v\" \"%v\" %v\n",
			respCount+1,
			info.Rating,
			info.ReviewNumber,
			info.Name,
			strings.Join(cuisines, " "),
			v.VendorURL,
		)

		respCount++
		if respCount >= limit {
			break
		}
	}

	return context.Response(fmt.Sprintf("```%v```", resp))
}

func getVendorInfo(code string) (v Vendor, err error) {
	url := fmt.Sprintf("https://tw.fd-api.com/api/v5/vendors/%v", code)
	resp, err := http.Get(url)
	if err != nil {
		return v, fmt.Errorf("get: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return v, fmt.Errorf("get: %v", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return v, fmt.Errorf("geadAll: %w", err)
	}

	body := map[string]interface{}{}
	err = json.Unmarshal(respBody, &body)
	if err != nil {
		return v, fmt.Errorf("unmarshal body: %w", err)
	}

	data, err := json.Marshal(body["data"])
	if err != nil {
		return v, fmt.Errorf("marshal: %w", err)
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return v, fmt.Errorf("unmarshal data: %w", err)
	}

	return
}
