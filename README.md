# userli-webhook-listener

This listens to webhooks sent by userli and reacts on them. The idea is to
auto-provisision and deprovision accounts in external services using requests to
these.

Supported webhook endpoints:
- User got created
- User got deleted

upported service to provision/deprovision accounts based on the event:
- Nextcloud
