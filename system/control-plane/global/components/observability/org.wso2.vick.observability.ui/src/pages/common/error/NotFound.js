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
        color: "#6e6e6e"
    },
    notFoundTitle: {
        margin: theme.spacing.unit,
        fontSize: "1.5em",
        fontWeight: 400,
        color: "#6e6e6e"
    },
    notFoundDescription: {
        fontSize: "1em",
        fontWeight: 300,
        color: "#808080"
    }
});

const NotFound = (props) => (
    <div className={props.classes.notFoundContainer}>
        <div className={props.classes.notFound}>
            <ErrorOutline className={props.classes.notFoundContentIndicator}/>
            <div className={props.classes.notFoundTitle}>
                {props.title ? props.title : "Unable to Find What You were Looking For"}
            </div>
            {
                props.description
                    ? (
                        <div className={props.classes.notFoundDescription}>
                            {props.description}
                        </div>
                    )
                    : null
            }
        </div>
    </div>
);

NotFound.propTypes = {
    classes: PropTypes.object.isRequired,
    title: PropTypes.string,
    description: PropTypes.string
};

export default withStyles(styles, {withTheme: true})(NotFound);
