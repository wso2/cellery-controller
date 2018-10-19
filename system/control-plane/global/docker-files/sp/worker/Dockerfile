from docker.wso2.com/wso2sp-worker:4.3.0

ARG VICK_PROFILE_FILES=./target/files

# copy the osgified jars
COPY --chown=wso2carbon:wso2 ${VICK_PROFILE_FILES}/lib/*.jar ${WSO2_SERVER_HOME}/lib/

# expose ports
EXPOSE 9090 9443 9712 9612 7711 7611 9091 9411 9092