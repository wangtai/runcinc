FROM golang
ENV GOPROXY=https://goproxy.cn
WORKDIR /home/runcic
COPY . /home/runcic
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
FROM centos
COPY --from=0 /home/runcic/runcic* /usr/bin/
RUN yum install -y podman e4fsprogs strace &&yum clean all &&rm -rf /tmp/ && sed -i "s/\"overlay\"/\"overlay2\"/g" /etc/containers/storage.conf
RUN mkdir /image /cic /cic/up /cic/work &&chmod +x /usr/bin/runcic.sh