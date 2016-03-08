FROM node:4.2.2-slim

RUN apt-get update && \
    apt-get install -y gem ruby-dev make autoconf python-dev

COPY . /restreamer
WORKDIR /restreamer

RUN apt-get install -y ruby rubygems-integration build-essential

RUN gem install RedCloth bundler
RUN bundle install

CMD ["bundle", "exec", "jekyll", "serve", "-H", "0.0.0.0", "-w"]