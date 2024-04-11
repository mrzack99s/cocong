FROM --platform=linux/amd64 oraclelinux:9
# FROM oraclelinux:9

RUN yum install -y libpcap gcc make git

RUN curl -L https://go.dev/dl/go1.21.8.linux-arm64.tar.gz -o go1.21.8.linux-arm64.tar.gz
RUN tar -C /usr/local -xzf go1.21.8.linux-arm64.tar.gz 

RUN rm -f go1.21.8.linux-arm64.tar.gz

RUN echo "export PATH=\$PATH:/usr/local/go/bin" >> /root/.bash_profile
RUN echo "export GOPATH=\$HOME/go" >> /root/.bash_profile
RUN echo "export PATH=\$PATH:\$GOPATH/bin" >> /root/.bash_profile