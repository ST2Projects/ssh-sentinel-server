FROM golang:1.22
LABEL authors="st2projects"

COPY ssh-sentinel-server ./

EXPOSE 8080
RUN ["mkdir", "/config"]
ENTRYPOINT ["./ssh-sentinel-server", "serve", "--config", "/resources/config.json"]
