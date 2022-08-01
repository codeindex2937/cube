package action

import (
	"cube/lib/context"
	"cube/lib/database"
	"cube/lib/utils"
	"errors"
	"fmt"
	"strings"

	"github.com/google/shlex"
	"gorm.io/gorm"
)

type Args struct {
	Create *CreateArgs `arg:"subcommand:create"`
	List   *ListArgs   `arg:"subcommand:list"`
	Delete *DeleteArgs `arg:"subcommand:delete"`
	Show   *ShowArgs   `arg:"subcommand:show"`
}

type Action struct {
	DB *gorm.DB
}

func (h *Action) Handle(req *context.ChatContext, args *Args) context.IResponse {
	switch {
	case args.Create != nil:
		return h.create(req, args.Create)
	case args.List != nil:
		return h.list(req, args.List)
	case args.Delete != nil:
		return h.delete(req, args.Delete)
	case args.Show != nil:
		return h.show(req, args.Show)
	}
	return utils.PrintHelp("action", args)
}

func displayAction(action database.Action) string {
	return fmt.Sprintf(`ID=%v "%v" Items=%v`,
		action.ActionID, action.Message, action.Items)
}

func NewActionResponse(db *gorm.DB, acitonID uint64, userID int) (resp context.ActionResponse, err error) {
	action := database.Action{}
	tx := db.First(&action, map[string]interface{}{
		"user_id":   userID,
		"action_id": acitonID,
	})
	if tx.Error != nil {
		return resp, tx.Error
	}

	resp = context.ActionResponse{
		CallbackID: fmt.Sprintf("%v", action.ActionID),
		Message:    action.Message,
	}

	items, err := shlex.Split(action.Items)
	if err != nil {
		return resp, errors.New("invalid action")
	}

	for _, a := range items {
		pair := strings.SplitN(a, ":", 2)

		if len(pair) < 2 {
			return resp, errors.New("invalid action")
		}
		text, value := pair[0], pair[1]
		resp.Actions = append(resp.Actions, context.Action{
			ActionType: "button",
			Name:       text,
			Value:      value,
			Text:       text,
			Style:      "green",
		})
	}
	resp.Actions = append(resp.Actions, context.Action{
		ActionType: "button",
		Name:       "close",
		Value:      "",
		Text:       "close",
		Style:      "green",
	})
	return
}
