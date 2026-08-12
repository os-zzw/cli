// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package im

import (
	"context"
	"net/http"

	"github.com/larksuite/cli/errs"
	"github.com/larksuite/cli/shortcuts/common"
)

const maxReadStatusMessageIDs = 50

var ImMessagesReadStatus = common.Shortcut{
	Service:     "im",
	Command:     "+messages-read-status",
	Description: "Batch query whether the current user has read up to 50 messages; preserves read, unread, and unexpected statuses",
	Risk:        "read",
	// 高敏权限由 Scope 平台白名单校验，不会出现在 UAT 的 scope 字段中，
	// 因此不能参与本地 OAuth scope 预检或生成重新授权提示。
	Scopes:    []string{},
	AuthTypes: []string{"user"},
	Flags: []common.Flag{
		{Name: "message-ids", Aliases: []string{"message-id"}, Required: true, Desc: "message IDs, comma-separated (1-50 om_xxx IDs)"},
	},
	Validate: func(ctx context.Context, runtime *common.RuntimeContext) error {
		_, err := buildMessagesReadStatusBody(runtime)
		return err
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		body, _ := buildMessagesReadStatusBody(runtime)
		return common.NewDryRunAPI().
			POST("/open-apis/im/v1/messages/batch_query_read_status").
			Body(body)
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		body, err := buildMessagesReadStatusBody(runtime)
		if err != nil {
			return err
		}
		data, err := runtime.CallAPITyped(http.MethodPost, "/open-apis/im/v1/messages/batch_query_read_status", nil, body)
		if err != nil {
			return normalizeAllowlistedUserScopeError(err, runtime.As(), "im:message.read_status:readonly")
		}
		runtime.OutFormat(data, nil, nil)
		return nil
	},
}

func buildMessagesReadStatusBody(runtime *common.RuntimeContext) (map[string]interface{}, error) {
	const param = "--message-ids"
	ids := common.SplitCSV(runtime.Str("message-ids"))
	if len(ids) == 0 {
		return nil, errs.NewValidationError(errs.SubtypeInvalidArgument, "--message-ids requires at least one om_ message ID").WithParam(param)
	}
	if len(ids) > maxReadStatusMessageIDs {
		return nil, errs.NewValidationError(errs.SubtypeInvalidArgument, "--message-ids supports at most %d IDs per request (got %d)", maxReadStatusMessageIDs, len(ids)).WithParam(param)
	}
	for _, id := range ids {
		if _, err := validateMessageIDForParam(id, param); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{"message_ids": ids}, nil
}
