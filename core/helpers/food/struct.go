package food

type CuisineType string

const (
	CHINESE   CuisineType = "chinese"
	PIZZA     CuisineType = "pizza"
	EASTERN   CuisineType = "eastern"
	EURO      CuisineType = "euro"
	HAMBURGER CuisineType = "hamburger"
	EXOTIC    CuisineType = "exotic"
	AMERICAN  CuisineType = "american"
	SOUTHERN  CuisineType = "southern"
	DESERT    CuisineType = "desert"
	DRINK     CuisineType = "vegeta"
	NOODLE    CuisineType = "drink"
	HEALTHY   CuisineType = "healthy"
	SNACK     CuisineType = "snack"
)

var cuisineType = map[CuisineType]int{
	"chinese":   166,
	"pizza":     165,
	"eastern":   164,
	"euro":      175,
	"hamburger": 177,
	"exotic":    183,
	"american":  179,
	"southern":  252,
	"noodle":    201,
	//
	"desert":  176,
	"vegeta":  186,
	"drink":   181,
	"healthy": 225,
	"snack":   214,
}

type NearByResp struct {
	Count     uint64                 `json:"count"`
	Feed      FeedData               `json:"feed"`
	RequestID string                 `json:"request_id"`
	Meta      map[string]interface{} `json:"meta"`
}

type FeedData struct {
	Count uint64       `json:"count"`
	Items []VendorList `json:"items"`
}

type FeedItem struct {
	Component uint64 `json:"count"`
	Headline  string `json:"headline"`
}

type VendorList struct {
	FeedItem
	Items []Vendor `json:"items"`
}

type Cuisine struct {
	ID     uint64 `json:"id"`
	Main   bool   `json:"main"`
	Name   string `json:"name"`
	UrlKey string `json:"url_key"`
}

type FoodCharactistic struct {
	ID           uint64 `json:"id"`
	IsHalal      bool   `json:"is_halal"`
	IsVegeTarian bool   `json:"is_vegetarian"`
	Name         string `json:"name"`
}

type VendorCharacteristics struct {
	Cuisines          []Cuisine          `json:"cuisines"`
	FoodCharactistics []FoodCharactistic `json:"food_characteristics"`
	PrimaryCuisine    Cuisine            `json:"primary_cuisine"`
}

type VendorMetadata struct {
	// "available_in": null,
	// "close_reasons": [],
	// "events": [],
	HasDiscount                bool   `json:"has_discount"`
	IsDeliveryAvailable        bool   `json:"is_delivery_available"`
	IsDineInAvailable          bool   `json:"is_dine_in_available"`
	IsExpressDeliveryAvailable bool   `json:"is_express_delivery_available"`
	IsFoodFeatureClosed        bool   `json:"is_flood_feature_closed"`
	IIsPickupAvailable         bool   `json:"is_pickup_available"`
	IsTemporaryClosed          bool   `json:"is_temporary_closed"`
	Timezone                   string `json:"timezone"`
}

type Tag struct {
	Code string `json:"code"`
	Text string `json:"text"`
}

type Vendor struct {
	// "accepts_instructions": true,
	// "chain": {
	//   "code": "",
	//   "main_vendor_code": "",
	//   "name": "",
	//   "url_key": ""
	// },
	// "city": {
	//   "name": "New Taipei City"
	// },
	// "custom_location_url": "",
	// "customer_type": "all",
	// "delivery_box": "",
	// "delivery_fee_type": "amount",
	// "delivery_provider_id": 0,
	// "description": "",
	// "discounts": [],
	// "discounts_info": [],
	// "payment_types": [],
	// "schedules": [],
	// "special_days": [],

	ID                uint64                `json:"id"`
	Code              string                `json:"code"`
	Name              string                `json:"name"`
	Address           string                `json:"address"`
	Distance          float64               `json:"distance"`
	Latitude          float64               `json:"latitude"`
	Longtitude        float64               `json:"longitude"`
	Metadata          VendorMetadata        `json:"metadata"`
	Tag               string                `json:"tag"`
	Rating            float64               `json:"rating"`
	ReviewNumber      uint64                `json:"review_number"`
	WebPath           string                `json:"web_path"`
	WebSite           string                `json:"website"`
	Charactistics     VendorCharacteristics `json:"characteristics"`
	FoodCharactistics []FoodCharactistic    `json:"food_characteristics"`
	Cuisines          []Cuisine             `json:"cuisines"`
	Budget            int                   `json:"budget"`
	IsActive          bool                  `json:"is_active"`
	Logo              string                `json:"logo"`
	// Score             float64               `json:"score"`
	//
	HasDeliveryProvider      bool   `json:"has_delivery_provider"`
	HasOnlinePayment         bool   `json:"has_online_payment"`
	HeroImage                string `json:"hero_image"`
	HeroListingImage         string `json:"hero_listing_image"`
	IsBestInCity             bool   `json:"is_best_in_city"`
	IsCheckoutCommentEnabled bool   `json:"is_checkout_comment_enabled"`
	IsDeliveryEnabled        bool   `json:"is_delivery_enabled"`
	IsNew                    bool   `json:"is_new"`
	IsNewUntil               string `json:"is_new_until"`
	IsPickupEnabled          bool   `json:"is_pickup_enabled"`
	IsPremium                bool   `json:"is_premium"`
	IsPreorderEnabled        bool   `json:"is_preorder_enabled"`
	IsPromoted               bool   `json:"is_promoted"`
	IsReplacementDishEnabled bool   `json:"is_replacement_dish_enabled"`
	IsServiceFeeEnabled      bool   `json:"is_service_fee_enabled"`
	IsServiceTaxEnabled      bool   `json:"is_service_tax_enabled"`
	IsServiceTaxVisible      bool   `json:"is_service_tax_visible"`
	IsTest                   bool   `json:"is_test"`
	IsVoucherEnabled         bool   `json:"is_voucher_enabled"`
	//
	IsVatDisabled               bool   `json:"is_vat_disabled"`
	IsVatIncludedInProductPrice bool   `json:"is_vat_included_in_product_price"`
	IsVatVisible                bool   `json:"is_vat_visible"`
	VatPercentageAmount         uint64 `json:"vat_percentage_amount"`
	//
	LoyaltyPercentageAmount    float64 `json:"loyalty_percentage_amount"`
	LoyaltyProgramEnabled      bool    `json:"loyalty_program_enabled"`
	MaximumExpressOrderAmount  uint64  `json:"maximum_express_order_amount"`
	MinimumDeliveryFee         float64 `json:"minimum_delivery_fee"`
	MinimumDeliveryTime        float64 `json:"minimum_delivery_time"`
	MinimumOrderAmount         float64 `json:"minimum_order_amount"`
	MinimumPickupTime          float64 `json:"minimum_pickup_time"`
	PremiumPosition            uint64  `json:"premium_position"`
	PrimaryCuisineID           uint64  `json:"primary_cuisine_id"`
	RedirectionURL             string  `json:"redirection_url"`
	ReviewWithCommentNumber    uint64  `json:"review_with_comment_number"`
	ServiceFeePercentageAmount uint64  `json:"service_fee_percentage_amount"`
	ServiceTaxPercentageAmount uint64  `json:"service_tax_percentage_amount"`
	URLKey                     string  `json:"url_key"`
	VendorPoints               uint64  `json:"vendor_points"`
	//
	Vertical        string   `json:"vertical"`
	VerticalParent  string   `json:"vertical_parent"`
	VerticalSegment string   `json:"vertical_segment"`
	VerticalTypeIDs []string `json:"vertical_type_ids"`
}
