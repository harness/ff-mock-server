############################
# STEP 1 build executable binary
############################
FROM golang:1.16 as builder

ARG FF_COMMIT
ARG FF_VERSION

WORKDIR /app

# Fetch dependencies.
COPY go.mod .
COPY go.sum .
COPY Makefile .
RUN make dep 

# Fetch all
COPY . .

# Generate Code and Build
RUN make build

RUN chmod +x ./wait-for-it.sh

############################
# STEP 2 build a small image
############################
FROM alpine:latest
RUN apk update && apk add --no-cache bash curl
# Import from builder.
COPY --from=builder /app/wait-for-it.sh /app/wait-for-it.sh
# Copy our static executable
COPY --from=builder /app/cmd/server/server /app/server
# Use an unprivileged user.
USER nobody:nogroup
ENTRYPOINT ["/app/server"]