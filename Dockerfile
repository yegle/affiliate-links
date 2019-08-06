FROM golang as build_env

ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o affiliate_link .

FROM gcr.io/distroless/static
COPY --from=build_env /app/affiliate_link /affiliate_link
ENV PORT 8999
ENTRYPOINT ["/affiliate_link"]
