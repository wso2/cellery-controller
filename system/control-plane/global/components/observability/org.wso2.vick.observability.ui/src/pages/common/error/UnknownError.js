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

import ErrorOutline from "@material-ui/icons/ErrorOutline";
import React from "react";
import {withStyles} from "@material-ui/core";
import * as PropTypes from "prop-types";

const styles = (theme) => ({
    unknownErrorContainer: {
        zIndex: -1,
        position: "absolute",
        top: 0,
        left: 0,
        height: "100%",
        width: "100%",
        display: "grid"
    },
    unknownError: {
        margin: "auto",
        textAlign: "center"
    },
    unknownErrorContentIndicator: {
        margin: theme.spacing.unit * 3,
        fontSize: "4em",
        color: "#808080"
    },
    unknownErrorContent: {
        fontSize: "1.5em",
        fontWeight: 400,
        color: "#808080"
    }
});

const UnknownError = (props) => (
    <div className={props.classes.unknownErrorContainer}>
        <div className={props.classes.unknownError}>
            <ErrorOutline className={props.classes.unknownErrorContentIndicator}/>
            <div className={props.classes.unknownErrorContent}>
                {props.message ? props.message : "Something Went Wrong"}
            </div>
        </div>
    </div>
);

UnknownError.propTypes = {
    classes: PropTypes.object.isRequired,
    message: PropTypes.string
};

export default withStyles(styles, {withTheme: true})(UnknownError);
