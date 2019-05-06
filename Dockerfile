FROM alpine:latest

RUN apk update && apk upgrade && \
    apk add ruby-dev ruby-rdoc ruby-bigdecimal ruby-webrick ruby-etc build-base zlib-dev

WORKDIR /gh-pages

COPY run.sh /usr/local/bin
COPY Gemfile /gh-pages

RUN gem install bundler -v '1.17.3' && bundle install

EXPOSE 4000
VOLUME /gh-pages

CMD ["/usr/local/bin/run.sh"]
