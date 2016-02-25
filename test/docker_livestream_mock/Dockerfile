FROM restreamer
ENV FFMPEG_VERSION 2.8.1
RUN mkdir /video
ADD loop.txt /video/
ADD small.mp4 /video/
ADD startLiveStream /usr/local/bin/
ADD nginx.conf /video/
RUN chmod +x /usr/local/bin/startLiveStream
WORKDIR /video

CMD ["/usr/local/bin/startLiveStream"]
