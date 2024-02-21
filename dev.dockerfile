FROM golang:1.22.0-bookworm

WORKDIR /home/app

CMD ["tail", "-f", "/dev/null"]
