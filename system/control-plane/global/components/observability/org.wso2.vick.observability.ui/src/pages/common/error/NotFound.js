/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import ErrorOutline from "@material-ui/icons/ErrorOutline";
import PropTypes from "prop-types";
import React from "react";
import {withStyles} from "@material-ui/core";

const styles = (theme) => ({
    notFoundContainer: {
        zIndex: -1,
        position: "absolute",
        top: 0,
        left: 0,
        height: "100%",
        width: "100%",
        display: "grid"
    },
    notFound: {
        margin: "auto",
        textAlign: "center"
    },
    notFoundContentIndicator: {
        margin: theme.spacing.unit * 3,
        fontSize: "4em",
        color: "#808080"
    },
    notFoundContent: {
        fontSize: "1.5em",
        fontWeight: 400,
        color: "#808080"
    }
});

const NotFound = (props) => (
    <div className={props.classes.notFoundContainer}>
        <div className={props.classes.notFound}>
            <ErrorOutline className={props.classes.notFoundContentIndicator}/>
            <div className={props.classes.notFoundContent}>
                {props.message ? props.message : "Unable to Find What You were Looking For"}
            </div>
        </div>
    </div>
);

NotFound.propTypes = {
    classes: PropTypes.object.isRequired,
    message: PropTypes.string
};

export default withStyles(styles, {withTheme: true})(NotFound);
