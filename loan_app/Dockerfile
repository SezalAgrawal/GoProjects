FROM golang:1.20 AS build-env
ENV APP_HOME /loan_app
WORKDIR $APP_HOME
COPY go.mod ./
COPY go.sum ./
COPY . $APP_HOME
ENV GOBIN $APP_HOME/bin
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -mod=vendor ./...

FROM ruby:3.1.2-slim
RUN apt-get update && apt-get install -y build-essential libpq-dev && rm -rf /var/lib/apt/lists/*
WORKDIR /repo
COPY migrations migrations
RUN bundle install --gemfile migrations/Gemfile
COPY scripts/docker-entrypoint.sh /usr/bin/
RUN chmod +x /usr/bin/docker-entrypoint.sh
COPY --from=build-env /loan_app/bin/* ./

ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["./loan_app"]