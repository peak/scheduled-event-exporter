FROM alpine
RUN apk add -U --no-cache ca-certificates
ADD app /app
CMD ["./app"]
