<<<<<<< HEAD
FROM node:4.2.2-slim

RUN apt-get update && \
    apt-get install -y gem ruby-dev make autoconf python-dev
=======
FROM node:5.5.0-slim

MAINTAINER datarhei <info@datarhei.org>

ENV FFMPEG_VERSION 2.8.5
ENV YASM_VERSION 1.3.0
ENV LAME_VERSION 3_99_5
ENV NGINX_VERSION 1.9.9
ENV NGINX_RTMP_VERSION 1.1.7

ENV SRC /usr/local
ENV LD_LIBRARY_PATH ${SRC}/lib
ENV PKG_CONFIG_PATH ${SRC}/lib/pkgconfig

ENV BUILDDEPS "autoconf automake gcc g++ libtool make nasm zlib1g-dev libssl-dev xz-utils cmake perl build-essential libpcre3-dev"

RUN rm -rf /var/lib/apt/lists/* && \
    apt-get update && \
    apt-get install -y --force-yes curl wget git libpcre3 tar ${BUILDDEPS}

# yasm
RUN DIR=$(mktemp -d) && cd ${DIR} && \
    curl -LOks https://www.tortall.net/projects/yasm/releases/yasm-${YASM_VERSION}.tar.gz && \
      tar xzvf yasm-${YASM_VERSION}.tar.gz && \
    cd yasm-${YASM_VERSION} && \
    ./configure --prefix="$SRC" --bindir="${SRC}/bin" && \
    make && \
    make install && \
    make distclean && \
    rm -rf ${DIR}

# x264
RUN DIR=$(mktemp -d) && cd ${DIR} && \
    git clone --depth 1 git://git.videolan.org/x264 && \
    cd x264 && \
    ./configure --prefix="$SRC" --bindir="${SRC}/bin" --enable-static && \
    make && \
    make install && \
    make distclean && \
    rm -rf ${DIR}

# libmp3lame
RUN DIR=$(mktemp -d) && cd ${DIR} && \
    curl -LOks https://github.com/rbrito/lame/archive/RELEASE__${LAME_VERSION}.tar.gz && \
      tar xzvf RELEASE__${LAME_VERSION}.tar.gz && \
    cd lame-RELEASE__${LAME_VERSION} && \
    ./configure --prefix="${SRC}" --bindir="${SRC}/bin" --disable-shared --enable-nasm && \
    make && \
    make install && \
    make distclean&& \
    rm -rf ${DIR}

# ffmpeg
RUN DIR=$(mktemp -d) && cd ${DIR} && \
    curl -LOks https://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz && \
      tar xzvf ffmpeg-${FFMPEG_VERSION}.tar.gz && \
    cd ffmpeg-${FFMPEG_VERSION} && \
    ./configure --prefix="${SRC}" --extra-cflags="-I${SRC}/include" --extra-ldflags="-L${SRC}/lib" --bindir="${SRC}/bin" \
      --extra-libs=-ldl --enable-version3 --enable-libmp3lame --enable-libx264 --enable-gpl \
      --enable-postproc --enable-nonfree --enable-avresample --disable-debug --enable-small --enable-openssl \
      --disable-doc --disable-ffserver && \
    make && \
    make install && \
    make distclean && \
    hash -r && \
    cd tools && \
    make qt-faststart && \
    cp qt-faststart ${SRC}/bin && \
    rm -rf ${DIR}
RUN echo "/usr/local/lib" > /etc/ld.so.conf.d/libc.conf
RUN ffmpeg -buildconf

# nginx-rtmp
RUN DIR=$(mktemp -d) && cd ${DIR} && \
    curl -LOks https://github.com/nginx/nginx/archive/release-${NGINX_VERSION}.tar.gz && \
      tar xzvf release-${NGINX_VERSION}.tar.gz && \
    curl -LOks https://github.com/arut/nginx-rtmp-module/archive/v${NGINX_RTMP_VERSION}.tar.gz && \
      tar xzvf v${NGINX_RTMP_VERSION}.tar.gz && \
    cd nginx-release-${NGINX_VERSION} && \
    auto/configure --with-http_ssl_module --add-module=../nginx-rtmp-module-${NGINX_RTMP_VERSION} && \
    make && \
    make install && \
    rm -rf ${DIR}

RUN apt-get purge -y --auto-remove ${BUILDDEPS} && \
    apt-get install -y --force-yes git && \
    rm -rf /tmp/*
>>>>>>> master

COPY . /restreamer
WORKDIR /restreamer

<<<<<<< HEAD
RUN apt-get install -y ruby rubygems-integration build-essential
=======
RUN npm install -g bower grunt-cli nodemon public-ip eslint@v2.0.0-beta.3 && \
    npm install && \
    grunt build && \
    npm prune --productio && \
    npm cache clean && \
    bower cache clean --allow-root

ENV RESTREAMER_USERNAME admin
ENV RESTREAMER_PASSWORD datarhei
>>>>>>> master

RUN gem install RedCloth bundler
RUN bundle install

CMD ["bundle", "exec", "jekyll", "serve", "-H", "0.0.0.0", "-w"]
