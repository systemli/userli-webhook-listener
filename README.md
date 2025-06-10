# userli-webhook-listener

[![Integration](https://github.com/systemli/userli-webhook-listener/actions/workflows/integration.yml/badge.svg)](https://github.com/systemli/userli-webhook-listener/actions/workflows/integration.yml) [![Quality](https://github.com/systemli/userli-webhook-listener/actions/workflows/quality.yml/badge.svg)](https://github.com/systemli/userli-webhook-listener/actions/workflows/quality.yml) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=systemli_userli-webhook-listener&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=systemli_userli-webhook-listener) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=systemli_userli-webhook-listener&metric=coverage)](https://sonarcloud.io/summary/new_code?id=systemli_userli-webhook-listener) [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=systemli_userli-webhook-listener&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=systemli_userli-webhook-listener)

## Debugging

Cheat sheet for sending requests with X-Signature header:

```shell
WEBHOOK_LISTENER_URL="https://example.org"
SECRET="secret"
PAYLOAD='{"type":"user.deleted","timestamp":"2025-01-01T00:00:00.000000Z","data":{"email":"user@example.org"}}'
SIGNATURE=$(printf '%s' "$PAYLOAD" | openssl dgst -sha256 -hmac "$SECRET" | sed 's/^.* //')

curl -i "$WEBHOOK_URL 
	-H "Content-Type: application/json" 
	-H "X-Signature: $SIGNATURE" 
	-d "$PAYLOAD"
```
