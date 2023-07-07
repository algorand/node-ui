# Build this Dockerfile with goreleaser.
# The binary must be present at /conduit
FROM debian:bullseye-slim

RUN groupadd --gid=999 --system algorand && \
    useradd --uid=999 --no-log-init --create-home --system --gid algorand algorand && \
    apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    update-ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

USER algorand

COPY nodeui /usr/local/bin/nodeui

ENTRYPOINT ["nodeui"]
