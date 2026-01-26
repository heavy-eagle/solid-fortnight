FROM amazoncorretto:17-alpine-jdk

# Install CA certificates and common CI tools
RUN apk update && apk upgrade --no-cache && apk add --no-cache curl nodejs npm ca-certificates git make tar docker-cli go

# add node apps
RUN npm install -g renovate @quasar/cli

# Add CA
COPY root-ca.crt /usr/local/share/ca-certificates/root-ca.crt

# Update trust store
RUN update-ca-certificates

# Default shell (bash)
CMD ["/bin/bash"]
