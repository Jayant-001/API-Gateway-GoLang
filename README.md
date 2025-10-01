### Request lifecycle — top to bottom (conceptual)

A single incoming request flows through:

1. Edge LB (outside our app) — TLS termination and basic filtering.

2. Gateway Ingress (our Go server)

- Parse request, assign request_id and trace_id.

- Logging middleware (structured log + start timestamp).

- Authentication middleware (JWT / API Key / token introspection).

- Authorization (scope/role checks for the route).

- Rate limit middleware (check + consume token).

- Routing (match host/path to route config).

- Reverse proxy / upstream call (with retries/timeouts/circuit-breaker).

- Response middleware (add X-RateLimit-*, correlation headers).

- Emit metrics and finish logs; flush traces.

- Admin interactions (separate) update the config store; gateways pick changes up via watch/poll.

- Keep this flow in your head — every feature fits into one of these stages.