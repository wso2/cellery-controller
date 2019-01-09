/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import CheckCircleOutline from "@material-ui/icons/CheckCircleOutline";
import ErrorOutline from "@material-ui/icons/ErrorOutline";
import {HelpOutline} from "@material-ui/icons";
import React from "react";
import StateHolder from "./state/stateHolder";
import withGlobalState from "./state";
import withColor, {ColorGenerator} from "./color";
import * as PropTypes from "prop-types";

const HealthIndicator = ({colorGenerator, globalState, value}) => {
    const color = colorGenerator.getColorForPercentage(value, globalState);
    let Icon;
    if (value < 0 || value > 1) {
        Icon = HelpOutline;
    } else if (value < globalState.get(StateHolder.CONFIG).percentageRangeMinValue.errorThreshold) {
        Icon = ErrorOutline;
    } else if (value < globalState.get(StateHolder.CONFIG).percentageRangeMinValue.warningThreshold) {
        Icon = ErrorOutline;
    } else {
        Icon = CheckCircleOutline;
    }
    return <Icon style={{color: color}}/>;
};

HealthIndicator.propTypes = {
    colorGenerator: PropTypes.instanceOf(ColorGenerator).isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    value: PropTypes.number.isRequired
};

export default withGlobalState(withColor(HealthIndicator));
