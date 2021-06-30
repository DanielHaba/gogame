FROM golang:1.16

ARG USER_ID=1000

ENV GOROOT="/usr/local/go"
ENV GOPATH="/home/go"
ENV GOOS="windows"
ENV GOARCH="386"
ENV CGO_ENABLED=1
ENV CXX_FOR_TARGET=i686-w64-mingw32-g++ 
ENV CC_FOR_TARGET=i686-w64-mingw32-gcc
ENV CC=i686-w64-mingw32-gcc-win32

RUN apt update -y && apt install -y \
	gcc-multilib\
	gcc-mingw-w64

RUN mkdir -p /usr/include/KHR \
	&& curl https://www.khronos.org/registry/EGL/api/KHR/khrplatform.h > /usr/include/KHR/khrplatform.h \
	&& mkdir -p /usr/i686-w64-mingw32/include/KHR \
	&& cp /usr/include/KHR/khrplatform.h /usr/i686-w64-mingw32/include/KHR/khrplatform.h

COPY profile/* /etc/profile.d/

RUN useradd -u ${USER_ID} go
USER go:go
WORKDIR /home/go

#RUN go get -x \
#	gopkgs \
#	go-outline \
#	gotests \
#	gomodifytags \
#	impl \
#	goplay \
#	dlv \
#	dlv-dap \
#	staticcheck \
#	gopls
