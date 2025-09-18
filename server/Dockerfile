FROM golang:1.25 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./server

FROM gcr.io/distroless/base-debian11 AS build-release-stage
COPY --from=build-stage /app/server /server
COPY --from=build-stage /app/config.yaml /config.yaml
COPY --from=build-stage /app/database/sql /database/sql
COPY --from=build-stage /app/templates /templates
# COPY --from=build-stage /app/static /static
EXPOSE 8080
USER nonroot:nonroot
CMD ["/server"]