package food

import (
	"cube/lib/context"
	"cube/lib/database"
)

type SetTagArgs struct {
	VendorID   uint64 `arg:"positional"`
	FoodTagID  uint64 `arg:"positional"`
	VendorName string `arg:"positional"`
	VendorCode string `arg:"positional"`
	VendorURL  string `arg:"positional"`
}

func (h *Food) attachTag(req *context.ChatContext, args *SetTagArgs) context.IResponse {
	record := database.FoodTagRelation{
		FoodTagID:  args.FoodTagID,
		VendorID:   args.VendorID,
		VendorName: args.VendorName,
		VendorCode: args.VendorCode,
		VendorURL:  args.VendorURL,
	}

	tx := h.DB.Where(map[string]interface{}{"vendor_id": args.VendorID}).Delete(&database.FoodTagRelation{})
	if tx.Error != nil {
		return context.NewTextResponse(tx.Error.Error())
	}

	if args.FoodTagID <= 0 {
		return context.Success
	}

	tx = h.DB.Create(&record)
	if tx.Error != nil {
		return context.NewTextResponse(tx.Error.Error())
	}

	return context.Success
}
