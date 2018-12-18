FROM node:11.2.0-slim
MAINTAINER WSO2 VICK Maintainers "dev@wso2.org"

ARG PORTAL_PORT=3000
ARG APP_DIR=/portal

WORKDIR ${APP_DIR}

COPY ./target/files ${APP_DIR}/

ENV PORTAL_PORT=${PORTAL_PORT}

EXPOSE ${PORTAL_PORT}

RUN npm install

ENTRYPOINT ["node", "serve.js"]
