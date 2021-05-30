FROM linuxserver/openssh-server

# Install git
RUN apk add --no-cache git

# Init repository
RUN mkdir /config/repos \
  && cd /config/repos \
  && git config --global init.defaultBranch main \
  && git init --bare scrt-test.git

COPY scrt_id_rsa.pub /config/scrt_id_rsa.pub