FROM ubuntu:latest
LABEL authors="grigorijmatukov"

ENTRYPOINT ["top", "-b"]