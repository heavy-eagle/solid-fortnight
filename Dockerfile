FROM amazoncorretto:17-alpine-jdk

# Install CA certificates and common CI tools
RUN apk update && apk upgrade --no-cache && apk add --no-cache curl nodejs npm ca-certificates git make tar docker-cli go bash jq grep zstd pandoc-cli mkdocs tectonic font-urw-base35 docker-cli kubectl

# install hugo
RUN curl -LO https://github.com/gohugoio/hugo/releases/download/v0.158.0/hugo_extended_0.158.0_linux-amd64.tar.gz --output hugo.tgz && \
  tar xzf hugo.tgz && mv hugo /usr/local/bin/hugo && chmod a+x /usr/local/bin/hugo && rm -f hugo.tgz

# install tea
RUN curl https://dl.gitea.com/tea/0.11.1/tea-0.11.1-linux-amd64 --output /usr/local/bin/tea && chmod a+x /usr/local/bin/tea

# add node apps
RUN npm install -g renovate @quasar/cli

# Add CA
COPY root-ca.crt /usr/local/share/ca-certificates/root-ca.crt

# Update trust store
RUN update-ca-certificates

# Default shell (bash)
CMD ["/bin/bash"]
