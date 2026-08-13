// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package im

import (
	"errors"
	"fmt"
	"slices"

	"github.com/larksuite/cli/errs"
	"github.com/larksuite/cli/internal/core"
)

// normalizeAllowlistedUserScopeError 保留服务端错误证据，同时移除对白名单权限无效的 OAuth 恢复建议。
func normalizeAllowlistedUserScopeError(err error, identity core.Identity, scope string) error {
	if err == nil || identity != core.AsUser {
		return err
	}
	var permissionErr *errs.PermissionError
	if !errors.As(err, &permissionErr) {
		return err
	}
	// 只有服务端明确返回高敏白名单权限时才替换恢复建议，普通权限仍走 OAuth 授权流程。
	if !slices.Contains(permissionErr.MissingScopes, scope) {
		return err
	}
	permissionErr.WithMissingScopes()
	permissionErr.Hint = fmt.Sprintf("verify that allowlisted scope %s is published for the current app in the Scope platform; this permission cannot be granted through OAuth", scope)
	return err
}
