FROM alpine:3.13

# set enviroment variables
ENV TYPING_PORT=7070

# install all depedencies in the container
RUN apk add --no-cache libc6-compat

# copy the local binary into the container
COPY ./bin/linux-amd64/typing /usr/local/typing

# set entrypoint for the image
ENTRYPOINT ["/usr/local/typing"]

# expose the container port
EXPOSE 7070
