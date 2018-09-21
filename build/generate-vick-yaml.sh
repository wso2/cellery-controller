#!/usr/bin/env bash

SCRIPT_ROOT=$(cd $(dirname ${BASH_SOURCE}) && pwd)
PROJECT_ROOT=$(cd ${SCRIPT_ROOT}/.. && pwd)
CONTROLLER_ROOT=${PROJECT_ROOT}/system/controller
BUILD_ROOT=${PROJECT_ROOT}/build/target

VICK_YAML_PATH=${BUILD_ROOT}/vick.yaml

rm -f ${VICK_YAML_PATH}

for yaml in ${CONTROLLER_ROOT}/artifacts/*.yaml; do
 cat ${yaml} >> ${BUILD_ROOT}/vick.yaml
 echo "---" >> ${BUILD_ROOT}/vick.yaml
done
