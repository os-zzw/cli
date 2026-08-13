// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package im

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/larksuite/cli/errs"
	"github.com/larksuite/cli/internal/core"
	"github.com/larksuite/cli/shortcuts/common"
	"github.com/spf13/cobra"
)

func newMessagesReadStatusTestRuntime(t *testing.T, messageIDs string) *common.RuntimeContext {
	t.Helper()

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("message-ids", "", "")
	if err := cmd.Flags().Set("message-ids", messageIDs); err != nil {
		t.Fatalf("Flags().Set() error = %v", err)
	}
	return &common.RuntimeContext{Cmd: cmd}
}

func TestNormalizeAllowlistedUserScopeErrorRemovesOAuthRecovery(t *testing.T) {
	source := errs.NewPermissionError(errs.SubtypeMissingScope, "missing allowlisted scope").
		WithCode(99991679).
		WithLogID("log-id").
		WithMissingScopes("im:message:get_as_user").
		WithHint("run auth login")

	got := normalizeAllowlistedUserScopeError(source, core.AsUser, "im:message:get_as_user")
	var permissionErr *errs.PermissionError
	if !errors.As(got, &permissionErr) {
		t.Fatalf("errors.As() = false, err = %v", got)
	}
	if len(permissionErr.MissingScopes) != 0 {
		t.Fatalf("MissingScopes = %v, want none", permissionErr.MissingScopes)
	}
	if strings.Contains(permissionErr.Hint, "auth login") || !strings.Contains(permissionErr.Hint, "Scope platform") {
		t.Fatalf("Hint = %q", permissionErr.Hint)
	}
	if permissionErr.Code != 99991679 || permissionErr.LogID != "log-id" {
		t.Fatalf("server evidence was not preserved: %#v", permissionErr)
	}
}

func TestNormalizeAllowlistedUserScopeErrorLeavesBotErrorUnchanged(t *testing.T) {
	source := errs.NewPermissionError(errs.SubtypeMissingScope, "missing bot scope").WithHint("open console")
	got := normalizeAllowlistedUserScopeError(source, core.AsBot, "im:message:get_as_user")
	if got != source || source.Hint != "open console" {
		t.Fatalf("bot error changed: %#v", source)
	}
}

func TestNormalizeAllowlistedUserScopeErrorLeavesOAuthScopeErrorUnchanged(t *testing.T) {
	source := errs.NewPermissionError(errs.SubtypeMissingScope, "missing OAuth scope").
		WithMissingScopes("im:message:readonly").
		WithHint("run auth login")

	got := normalizeAllowlistedUserScopeError(source, core.AsUser, "im:message:get_as_user")
	if got != source {
		t.Fatalf("error instance changed: got %p, want %p", got, source)
	}
	if !reflect.DeepEqual(source.MissingScopes, []string{"im:message:readonly"}) {
		t.Fatalf("MissingScopes = %v, want [im:message:readonly]", source.MissingScopes)
	}
	if source.Hint != "run auth login" {
		t.Fatalf("Hint = %q, want OAuth recovery hint", source.Hint)
	}
}

func TestBuildMessagesReadStatusBody(t *testing.T) {
	runtime := newMessagesReadStatusTestRuntime(t, "om_one, om_two")

	got, err := buildMessagesReadStatusBody(runtime)
	if err != nil {
		t.Fatalf("buildMessagesReadStatusBody() error = %v", err)
	}
	want := map[string]interface{}{
		"message_ids": []string{"om_one", "om_two"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildMessagesReadStatusBody() = %#v, want %#v", got, want)
	}
}

func TestBuildMessagesReadStatusBodyRejectsInvalidInputs(t *testing.T) {
	tests := []struct {
		name       string
		messageIDs string
	}{
		{name: "empty", messageIDs: ""},
		{name: "invalid prefix", messageIDs: "oc_not_message"},
		{name: "more than fifty", messageIDs: strings.Join(makeReadStatusMessageIDs(51), ",")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runtime := newMessagesReadStatusTestRuntime(t, tt.messageIDs)
			_, err := buildMessagesReadStatusBody(runtime)
			problem, ok := errs.ProblemOf(err)
			if !ok {
				t.Fatalf("errs.ProblemOf() ok = false, err = %v", err)
			}
			if problem.Subtype != errs.SubtypeInvalidArgument {
				t.Fatalf("problem.Subtype = %q, want %q", problem.Subtype, errs.SubtypeInvalidArgument)
			}
		})
	}
}

func TestMessagesReadStatusShortcutContract(t *testing.T) {
	if !reflect.DeepEqual(ImMessagesReadStatus.AuthTypes, []string{"user"}) {
		t.Fatalf("AuthTypes = %v, want [user]", ImMessagesReadStatus.AuthTypes)
	}
	if got := ImMessagesReadStatus.ScopesForIdentity("user"); !reflect.DeepEqual(got, []string{"im:message:readonly"}) {
		t.Fatalf("user preflight scopes = %v, want [im:message:readonly]", got)
	}
	if got := ImMessagesReadStatus.DeclaredScopesForIdentity("user"); !reflect.DeepEqual(got, []string{"im:message:readonly"}) {
		t.Fatalf("declared user scopes = %v, want [im:message:readonly]", got)
	}
	if ImMessagesReadStatus.Risk != "read" {
		t.Fatalf("Risk = %q, want read", ImMessagesReadStatus.Risk)
	}
}

func makeReadStatusMessageIDs(count int) []string {
	ids := make([]string, count)
	for i := range ids {
		ids[i] = fmt.Sprintf("om_%d", i)
	}
	return ids
}
