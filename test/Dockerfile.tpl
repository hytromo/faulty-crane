FROM alpine:latest

RUN echo RANDOM_REPLACE >/tmp/for_randomization

CMD ["sleep", "inf"]
