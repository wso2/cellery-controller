FROM scratch
COPY vick-controller /
ENTRYPOINT ["/vick-controller","-logtostderr=true"]
