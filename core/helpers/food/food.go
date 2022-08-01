package food

import (
	"net/http"

	"cube/config"
	"cube/lib/context"
	"cube/lib/database"
	"cube/lib/logger"
	"cube/lib/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VendorData struct {
	ID       uint64   `json:"id"`
	Lat      float64  `json:"lat"`
	Lng      float64  `json:"lng"`
	Icon     string   `json:"icon"`
	Rating   float64  `json:"rating"`
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Code     string   `json:"code"`
	Cuisines []string `json:"cuisines"`
}

var log = logger.Log

type Args struct {
	List      *ListArgs   `arg:"subcommand:list"`
	AttachTag *SetTagArgs `arg:"subcommand:attach_tag"`
	Poll      *PollArgs   `arg:"subcommand:poll"`
}

type Food struct {
	DB *gorm.DB
}

func (h *Food) Handle(req *context.ChatContext, args *Args) context.IResponse {
	switch {
	case args.List != nil:
		return h.list(req, args.List)
	case args.AttachTag != nil:
		return h.attachTag(req, args.AttachTag)
	case args.Poll != nil:
		return h.poll(req, args.Poll)
	}
	return utils.PrintHelp("food", args)
}

func (h *Food) RenderMap(c *gin.Context) {
	cuisineTypes := []CuisineType{}
	vendorMarkers := []VendorData{}
	for _, v := range listVendors(config.Conf.Map.Lat, config.Conf.Map.Lng, cuisineTypes) {
		if !v.IsActive || !v.IsDeliveryEnabled || !v.Metadata.IsDeliveryAvailable {
			continue
		}

		cuisines := []string{}
		for _, c := range v.Cuisines {
			cuisines = append(cuisines, c.Name)
		}

		vendorMarkers = append(vendorMarkers, VendorData{
			v.ID,
			v.Latitude,
			v.Longtitude,
			v.HeroImage,
			v.Rating,
			v.Name,
			v.RedirectionURL,
			v.Code,
			cuisines,
		})
	}

	relations := []database.FoodTagRelation{}
	tx := h.DB.Find(&relations)
	if tx.Error != nil {
		c.AbortWithError(500, tx.Error)
		return
	}

	c.HTML(http.StatusOK, "foodmap.tmpl", gin.H{
		"apiKey":    config.Conf.Map.Key,
		"lat":       config.Conf.Map.Lat,
		"lng":       config.Conf.Map.Lng,
		"vendors":   vendorMarkers,
		"relations": relations,
	})
}
