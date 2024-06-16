FROM golang:1.22

LABEL authors="st2projects"
LABEL org.opencontainers.image.source=https://github.com/st2projects/ssh-sentinel-server

COPY dist/bin/ssh-sentinel-server ./
COPY dist/bin/ssh-sentinel-client ./

EXPOSE 8080
RUN ["mkdir", "/resources"]
ENTRYPOINT ["./ssh-sentinel-server", "serve", "--config", "/resources/config.json"]
