// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package im

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/larksuite/cli/errs"
	"github.com/larksuite/cli/shortcuts/common"
	"github.com/spf13/cobra"
)

func newMessageReadUsersTestRuntime(t *testing.T, transport http.RoundTripper, stringsMap map[string]string, bools map[string]bool, ints map[string]int) *common.RuntimeContext {
	t.Helper()

	runtime := newUserShortcutRuntime(t, transport)
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("message-id", "", "")
	cmd.Flags().String("user-id-type", "open_id", "")
	cmd.Flags().Int("page-size", 100, "")
	cmd.Flags().String("page-token", "", "")
	cmd.Flags().Bool("page-all", false, "")
	cmd.Flags().Int("page-limit", 10, "")
	for name, value := range stringsMap {
		if err := cmd.Flags().Set(name, value); err != nil {
			t.Fatalf("Flags().Set(%q) error = %v", name, err)
		}
	}
	for name, value := range bools {
		if err := cmd.Flags().Set(name, strconv.FormatBool(value)); err != nil {
			t.Fatalf("Flags().Set(%q) error = %v", name, err)
		}
	}
	for name, value := range ints {
		if err := cmd.Flags().Set(name, strconv.Itoa(value)); err != nil {
			t.Fatalf("Flags().Set(%q) error = %v", name, err)
		}
	}
	runtime.Cmd = cmd
	return runtime
}

func TestMessageReadUsersScopesByIdentity(t *testing.T) {
	if !reflect.DeepEqual(ImMessageReadUsers.ScopesForIdentity("user"), []string{"im:message:get_as_user"}) {
		t.Fatalf("user scopes = %v", ImMessageReadUsers.ScopesForIdentity("user"))
	}
	if !reflect.DeepEqual(ImMessageReadUsers.ScopesForIdentity("bot"), []string{"im:message:readonly"}) {
		t.Fatalf("bot scopes = %v", ImMessageReadUsers.ScopesForIdentity("bot"))
	}
	if !reflect.DeepEqual(ImMessageReadUsers.AuthTypes, []string{"user", "bot"}) {
		t.Fatalf("AuthTypes = %v", ImMessageReadUsers.AuthTypes)
	}
}

func TestBuildMessageReadUsersParams(t *testing.T) {
	runtime := newMessageReadUsersTestRuntime(t, shortcutRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("unexpected request")
	}), map[string]string{"message-id": "om_test"}, nil, nil)

	got, err := buildMessageReadUsersParams(runtime, "cursor")
	if err != nil {
		t.Fatalf("buildMessageReadUsersParams() error = %v", err)
	}
	want := map[string]interface{}{
		"user_id_type": "open_id",
		"page_size":    100,
		"page_token":   "cursor",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildMessageReadUsersParams() = %#v, want %#v", got, want)
	}
}

func TestFetchMessageReadUsersAggregatesPages(t *testing.T) {
	calls := 0
	transport := shortcutRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/open-apis/im/v1/messages/om_test/read_users" {
			return nil, fmt.Errorf("unexpected path: %s", req.URL.Path)
		}
		calls++
		if req.URL.Query().Get("page_token") == "" {
			return shortcutJSONResponse(200, map[string]interface{}{
				"code": 0,
				"data": map[string]interface{}{
					"items":      []interface{}{map[string]interface{}{"user_id": "ou_one", "timestamp": "1"}},
					"has_more":   true,
					"page_token": "next",
				},
			}), nil
		}
		return shortcutJSONResponse(200, map[string]interface{}{
			"code": 0,
			"data": map[string]interface{}{
				"items":    []interface{}{map[string]interface{}{"user_id": "ou_two", "tenant_key": "tenant"}},
				"has_more": false,
			},
		}), nil
	})
	runtime := newMessageReadUsersTestRuntime(t, transport, map[string]string{"message-id": "om_test"}, map[string]bool{"page-all": true}, map[string]int{"page-limit": 0})

	got, err := fetchMessageReadUsers(context.Background(), runtime)
	if err != nil {
		t.Fatalf("fetchMessageReadUsers() error = %v", err)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
	items, _ := got["items"].([]interface{})
	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}
	if got["has_more"] != false || got["page_token"] != "" || got["total"] != 2 {
		t.Fatalf("result metadata = %#v", got)
	}
}

func TestMessageReadUsersValidationRejectsInvalidInputs(t *testing.T) {
	tests := []struct {
		name       string
		stringsMap map[string]string
		ints       map[string]int
	}{
		{name: "invalid message", stringsMap: map[string]string{"message-id": "oc_invalid"}},
		{name: "invalid id type", stringsMap: map[string]string{"message-id": "om_test", "user-id-type": "email"}},
		{name: "page size too large", stringsMap: map[string]string{"message-id": "om_test"}, ints: map[string]int{"page-size": 101}},
		{name: "negative page limit", stringsMap: map[string]string{"message-id": "om_test"}, ints: map[string]int{"page-limit": -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runtime := newMessageReadUsersTestRuntime(t, shortcutRoundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("unexpected request")
			}), tt.stringsMap, nil, tt.ints)
			err := validateMessageReadUsers(runtime)
			problem, ok := errs.ProblemOf(err)
			if !ok || problem.Subtype != errs.SubtypeInvalidArgument {
				t.Fatalf("problem = %#v, err = %v", problem, err)
			}
		})
	}
}
