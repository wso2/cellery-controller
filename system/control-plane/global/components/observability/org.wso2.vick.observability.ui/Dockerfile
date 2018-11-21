FROM node:11.2.0-slim
MAINTAINER WSO2 VICK Maintainers "dev@wso2.org"

ARG APP_DIR=/portal
ARG PORTAL_PORT=3000

WORKDIR ${APP_DIR}

COPY ./node-server/ ${APP_DIR}/
COPY ./build/ ${APP_DIR}/public/

ENV PORTAL_PORT=${PORTAL_PORT}

EXPOSE ${PORTAL_PORT}

RUN npm install
CMD node serve.js
