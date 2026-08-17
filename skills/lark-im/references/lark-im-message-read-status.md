# IM message read status

> **Prerequisite:** Read [`../lark-shared/SKILL.md`](../../lark-shared/SKILL.md) first for authentication and global parameters.

Use two focused shortcuts for message read-status queries:

- `im +messages-read-status` queries whether the current user has read 1–50 messages.
- `im +message-read-users` lists users who have read one message and supports automatic pagination.

Both underlying OpenAPIs support user identity through a user access token (UAT). `+message-read-users` additionally supports bot identity through a tenant access token (TAT).

## Identity and scopes

| Shortcut | Identity | Scope |
|---|---|---|
| `+messages-read-status` | user only | `im:message:readonly` (recommended), `im:message`, or `im:message:get_as_user` |
| `+message-read-users` | user | `im:message:get_as_user` |
| `+message-read-users` | bot | `im:message:readonly` |

For `+message-read-users`, the caller must still be in the chat. A user can query only messages they sent within the last seven days, while a bot can query only messages sent by that bot within the last seven days.

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

The command accepts 1–50 comma-separated `om_` message IDs. The three scopes above are alternatives; any one is sufficient, and the CLI recommends the least-privileged OAuth scope `im:message:readonly`. The response keeps the OpenAPI response unchanged:

- `items[].message_id` and `items[].is_read` contain statuses the server could determine.
- `invalid_message_ids` contains messages that do not exist, are not visible to the current user, or do not support this query. The API deliberately does not expose a more specific reason.

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
lark-cli im messages read_status --data '{"message_ids":["om_xxx"]}' --as user
lark-cli im messages read_users --params '{"message_id":"om_xxx","user_id_type":"open_id"}' --as user
```

Prefer the shortcuts for flag validation, identity-specific scope hints, and read-users auto-pagination.

## Troubleshooting

| Symptom | Meaning | Action |
|---|---|---|
| `--as bot is not supported` for read status | The batch endpoint requires user identity | Switch to `--as user` |
| Missing `im:message:readonly` or `im:message` | A regular OAuth scope has not been granted | Follow the CLI authorization hint to grant one supported scope |
| Missing `im:message:get_as_user` | The allowlisted high-sensitivity scope has not taken effect for this app | Verify the app ID and publication state in the Scope platform; do not request it through OAuth |
| Bot permission denied | The application lacks a bot scope | Open the `console_url` from the typed error and enable the requested scope |
| Empty read-user list | No user has read the message, or sender/time constraints are not met | Verify chat membership, the message sender, and the seven-day window |

## References

- [lark-im](../SKILL.md)
- [lark-shared](../../lark-shared/SKILL.md)
