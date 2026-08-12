// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package im

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/larksuite/cli/errs"
	"github.com/larksuite/cli/internal/output"
	"github.com/larksuite/cli/internal/validate"
	"github.com/larksuite/cli/shortcuts/common"
)

const (
	messageReadUsersDefaultPageSize  = 100
	messageReadUsersMaxPageSize      = 100
	messageReadUsersDefaultPageLimit = 10
)

var ImMessageReadUsers = common.Shortcut{
	Service:     "im",
	Command:     "+message-read-users",
	Description: "List users who have read one message; supports user and bot identities with optional automatic pagination",
	Risk:        "read",
	UserScopes:  []string{"im:message:get_as_user"},
	BotScopes:   []string{"im:message:readonly"},
	AuthTypes:   []string{"user", "bot"},
	Flags: []common.Flag{
		{Name: "message-id", Required: true, Desc: "message ID (om_xxx)"},
		{Name: "user-id-type", Default: "open_id", Desc: "user ID type returned in each item", Enum: []string{"open_id", "union_id", "user_id"}},
		{Name: "page-size", Aliases: []string{"limit"}, Type: "int", Default: fmt.Sprintf("%d", messageReadUsersDefaultPageSize), Desc: fmt.Sprintf("page size (1-%d)", messageReadUsersMaxPageSize)},
		{Name: "page-token", Desc: "starting pagination cursor"},
		{Name: "page-all", Type: "bool", Desc: "automatically paginate through all pages (capped by --page-limit)"},
		{Name: "page-limit", Type: "int", Default: fmt.Sprintf("%d", messageReadUsersDefaultPageLimit), Desc: fmt.Sprintf("max pages to fetch with --page-all (default %d, 0 = unlimited)", messageReadUsersDefaultPageLimit)},
	},
	Validate: func(ctx context.Context, runtime *common.RuntimeContext) error {
		return validateMessageReadUsers(runtime)
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		messageID := strings.TrimSpace(runtime.Str("message-id"))
		params, _ := buildMessageReadUsersParams(runtime, strings.TrimSpace(runtime.Str("page-token")))
		dry := common.NewDryRunAPI().
			GET("/open-apis/im/v1/messages/:message_id/read_users").
			Set("message_id", messageID).
			Params(params)
		if runtime.Bool("page-all") {
			dry.Desc("Auto-paginates through read-user pages, capped by --page-limit when greater than zero")
		}
		return dry
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		data, err := fetchMessageReadUsers(ctx, runtime)
		if err != nil {
			return err
		}
		count, _ := data["total"].(int)
		runtime.OutFormat(data, &output.Meta{Count: count}, nil)
		return nil
	},
}

func validateMessageReadUsers(runtime *common.RuntimeContext) error {
	if _, err := validateMessageIDForParam(runtime.Str("message-id"), "--message-id"); err != nil {
		return err
	}
	switch idType := strings.TrimSpace(runtime.Str("user-id-type")); idType {
	case "open_id", "union_id", "user_id":
	default:
		return errs.NewValidationError(errs.SubtypeInvalidArgument, "--user-id-type must be one of open_id, union_id, or user_id").WithParam("--user-id-type")
	}
	if _, err := common.ValidatePageSizeTyped(runtime, "page-size", messageReadUsersDefaultPageSize, 1, messageReadUsersMaxPageSize); err != nil {
		return err
	}
	if runtime.Int("page-limit") < 0 {
		return errs.NewValidationError(errs.SubtypeInvalidArgument, "--page-limit must be a non-negative integer").WithParam("--page-limit")
	}
	return nil
}

func buildMessageReadUsersParams(runtime *common.RuntimeContext, pageToken string) (map[string]interface{}, error) {
	pageSize, err := common.ValidatePageSizeTyped(runtime, "page-size", messageReadUsersDefaultPageSize, 1, messageReadUsersMaxPageSize)
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"user_id_type": strings.TrimSpace(runtime.Str("user-id-type")),
		"page_size":    pageSize,
	}
	if pageToken != "" {
		params["page_token"] = pageToken
	}
	return params, nil
}

func fetchMessageReadUsers(ctx context.Context, runtime *common.RuntimeContext) (map[string]interface{}, error) {
	messageID, err := validateMessageIDForParam(runtime.Str("message-id"), "--message-id")
	if err != nil {
		return nil, err
	}
	pageToken := strings.TrimSpace(runtime.Str("page-token"))
	pageAll := runtime.Bool("page-all")
	pageLimit := runtime.Int("page-limit")
	items := []interface{}{}
	hasMore := false
	nextToken := ""
	apiPath := fmt.Sprintf("/open-apis/im/v1/messages/%s/read_users", validate.EncodePathSegment(messageID))

	for page := 1; ; page++ {
		if pageAll {
			fmt.Fprintf(runtime.IO().ErrOut, "[page %d] fetching read users...\n", page)
		}
		params, err := buildMessageReadUsersParams(runtime, pageToken)
		if err != nil {
			return nil, err
		}
		data, err := runtime.CallAPITyped(http.MethodGet, apiPath, params, nil)
		if err != nil {
			return nil, err
		}
		if pageItems, ok := data["items"].([]interface{}); ok {
			items = append(items, pageItems...)
		}
		hasMore, nextToken = common.PaginationMeta(data)
		if !pageAll || !hasMore || nextToken == "" {
			break
		}
		if nextToken == pageToken {
			fmt.Fprintln(runtime.IO().ErrOut, "Stopping pagination: server returned a non-advancing page_token.")
			break
		}
		if pageLimit > 0 && page >= pageLimit {
			fmt.Fprintf(runtime.IO().ErrOut, "[pagination] reached page limit (%d); results may be incomplete. Use --page-limit 0 to fetch all pages.\n", pageLimit)
			break
		}
		pageToken = nextToken
	}

	return map[string]interface{}{
		"items":      items,
		"has_more":   hasMore,
		"page_token": nextToken,
		"total":      len(items),
	}, nil
}
