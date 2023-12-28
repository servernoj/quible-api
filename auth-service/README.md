# Setup
The `auth-service` is dependent on several generated assets that can be produced by running `go generate ./...` from the service root directory. That command will generate
- error strigifier used to produce the output of `docs/errors` response
- email templates defined in `services/emailService` directory