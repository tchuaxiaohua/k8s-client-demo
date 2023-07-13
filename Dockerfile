FROM golang:1.20 as builder

WORKDIR /apps

COPY ./ /apps
RUN export GOPROXY=https://goproxy.cn \
    && go build  -ldflags "-s -w" -o get-pod \
    && chmod +x get-pod

FROM alpine
LABEL maintainer="tchua"
COPY --from=builder /apps/get-pod  /apps/

RUN echo -e  "http://mirrors.aliyun.com/alpine/v3.15/main\nhttp://mirrors.aliyun.com/alpine/v3.15/community" >  /etc/apk/repositories \
&& apk  update && apk --no-cache add tzdata gcompat libc6-compat \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Shanghai/Asia" > /etc/timezone \
&& apk del tzdata \
&& ln -s /lib/libc.so.6 /usr/lib/libresolv.so.2

WORKDIR /apps


CMD ["./get-pod"]