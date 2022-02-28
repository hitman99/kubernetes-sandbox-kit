FROM golang:1.17

COPY . /gosrc
WORKDIR /gosrc
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ksk

FROM node:16-alpine
COPY frontend /src
WORKDIR /src

RUN npm install && npm run build

FROM alpine:3.15
LABEL maintainer="tomas@adomavicius.com"

RUN apk --no-cache add ca-certificates && adduser gopher -D -H -u 1133
WORKDIR /workshop
COPY --from=0 /gosrc/ksk ksk
COPY --from=1 /src/dist frontend/dist
ENV PATH="/workshop/:${PATH}"

EXPOSE 8080
USER 1133

CMD ["ksk", "registration"]