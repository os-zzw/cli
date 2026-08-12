// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package im

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/larksuite/cli/errs"
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
	if !reflect.DeepEqual(ImMessagesReadStatus.Scopes, []string{"im:message.read_status:readonly"}) {
		t.Fatalf("Scopes = %v", ImMessagesReadStatus.Scopes)
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
