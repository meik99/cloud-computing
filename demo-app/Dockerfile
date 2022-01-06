FROM node:latest as builder

RUN mkdir /app
WORKDIR /app

ADD . .

RUN npm install && \
    npm install -g @angular/cli && \
    ng build --configuration=production


FROM nginx:latest
COPY --from=builder /app/dist/demo-app/ /usr/share/nginx/html

EXPOSE 80
