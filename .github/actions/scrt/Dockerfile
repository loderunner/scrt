FROM loderunner/scrt:0.3.3 as scrt

FROM alpine:3.16.0

COPY --from=scrt /scrt /scrt
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]