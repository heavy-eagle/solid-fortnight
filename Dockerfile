FROM amazoncorretto:21-alpine-jdk

# node version
ENV NODE_VERSION=22.22.3

# Install CA certificates and common CI tools
RUN apk update && apk upgrade --no-cache && apk add --no-cache curl ca-certificates xz git make tar docker-cli bash jq grep zstd pandoc-cli mkdocs tectonic font-urw-base35 kubectl openssl

# install hugo
RUN curl -L https://github.com/gohugoio/hugo/releases/download/v0.158.0/hugo_extended_0.158.0_linux-amd64.tar.gz --output hugo.tgz && \
  tar xzf hugo.tgz && mv hugo /usr/local/bin/hugo && chmod a+x /usr/local/bin/hugo && rm -f hugo.tgz

# install tea
RUN curl https://dl.gitea.com/tea/0.11.1/tea-0.11.1-linux-amd64 --output /usr/local/bin/tea && chmod a+x /usr/local/bin/tea

# install go
RUN LATEST=$(curl -s https://go.dev/VERSION?m=text | head -n 1) && curl -L "https://go.dev/dl/${LATEST}.linux-amd64.tar.gz" --output go.tgz && \
  tar -C /usr/local -xzf go.tgz && rm go.tgz

# Add Go to PATH
ENV PATH="/usr/local/go/bin:${PATH}"

# Build EST Client
COPY estcli estcli
RUN cd estcli && go build -o /usr/local/bin/estcli && cd - && rm -rf estcli

# Install ACME client
RUN curl https://get.acme.sh | sh -s email=my@example.com --install --home /usr/local/acme.sh && ln -s /usr/local/acme.sh/acme.sh /usr/local/bin/acme.sh

# install node
RUN curl -fsSL \
  https://unofficial-builds.nodejs.org/download/release/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64-musl.tar.xz \
  -o /tmp/node.tar.xz \
 && tar -xJf /tmp/node.tar.xz -C /usr/local --strip-components=1 \
 && rm /tmp/node.tar.xz

# add node apps
RUN npm install -g renovate @quasar/cli wrangler @usebruno/cli

# Add CA
COPY root-ca.crt /usr/local/share/ca-certificates/root-ca.crt

# Update trust store
RUN update-ca-certificates

# Default shell (bash)
CMD ["/bin/bash"]
