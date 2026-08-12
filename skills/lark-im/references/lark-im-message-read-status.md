# IM message read status

> **Prerequisite:** Read [`../lark-shared/SKILL.md`](../../lark-shared/SKILL.md) first for authentication and global parameters.

Use two focused shortcuts for message read-status queries:

- `im +messages-read-status` queries whether the current user has read 1–50 messages.
- `im +message-read-users` lists users who have read one message and supports automatic pagination.

Both underlying OpenAPIs support user identity through a user access token (UAT). `+message-read-users` additionally supports bot identity through a tenant access token (TAT).

## Identity and scopes

| Shortcut | Identity | Scope |
|---|---|---|
| `+messages-read-status` | user only | `im:message.read_status:readonly` |
| `+message-read-users` | user | `im:message:get_as_user` |
| `+message-read-users` | bot | `im:message:readonly` |

For `+message-read-users`, a user can query only messages they sent within the last seven days. A bot must be in the chat and can query only messages sent by that bot within the last seven days.

## Batch query the current user's read status

```bash
# Preview one request
lark-cli im +messages-read-status \
  --message-ids om_xxx,om_yyy \
  --as user \
  --dry-run

# Execute with a user access token
lark-cli im +messages-read-status \
  --message-ids om_xxx,om_yyy \
  --as user \
  --json
```

The command accepts 1–50 comma-separated `om_` message IDs. The response keeps the server values `read`, `unread`, and `unexpected` unchanged.

Never treat `unexpected` as `unread`. Inspect `unexpected_reason`:

- `invalid`: the message ID is invalid.
- `no_permission`: the current user cannot inspect the message.
- `not_support`: the message does not support read-status lookup.

## List users who read one message

```bash
# Fetch one page as the current user
lark-cli im +message-read-users \
  --message-id om_xxx \
  --as user \
  --json

# Fetch every page as a bot, bounded to ten pages by default
lark-cli im +message-read-users \
  --message-id om_xxx \
  --user-id-type open_id \
  --page-all \
  --as bot \
  --json
```

Pagination flags:

- `--page-size`: 1–100, default 100.
- `--page-token`: start from a known cursor.
- `--page-all`: continue until the endpoint is exhausted.
- `--page-limit`: maximum pages with `--page-all`; default 10, and 0 means unlimited.

The command preserves each server item, including `user_id_type`, `user_id`, `timestamp`, and `tenant_key`. If pagination reaches its limit, stderr states that results may be incomplete and the JSON response retains `has_more` and `page_token` for resumption.

## Raw API commands

When Registry MR !128 is published, the corresponding raw commands remain available:

```bash
lark-cli im messages batch_get_read_status --data '{"message_ids":["om_xxx"]}' --as user
lark-cli im messages read_users --params '{"message_id":"om_xxx","user_id_type":"open_id"}' --as user
```

Prefer the shortcuts for flag validation, identity-specific scope hints, and read-users auto-pagination.

## Troubleshooting

| Symptom | Meaning | Action |
|---|---|---|
| `--as bot is not supported` for read status | The batch endpoint requires user identity | Switch to `--as user` |
| Missing `im:message.read_status:readonly` | The allowlisted high-sensitivity scope has not taken effect for this app | Verify the app ID and publication state in the Scope platform; do not request it through OAuth |
| Missing `im:message:get_as_user` | The allowlisted high-sensitivity scope has not taken effect for this app | Verify the app ID and publication state in the Scope platform; do not request it through OAuth |
| Bot permission denied | The application lacks a bot scope | Open the `console_url` from the typed error and enable the requested scope |
| Empty read-user list | No user has read the message, or sender/time constraints are not met | Verify the message sender and seven-day window |

## References

- [lark-im](../SKILL.md)
- [lark-shared](../../lark-shared/SKILL.md)
