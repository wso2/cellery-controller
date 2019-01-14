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

import React from "react";
import UnknownError from "./UnknownError";
import * as PropTypes from "prop-types";

/**
 * Error Boundary to catch error in React Components.
 * This Component can be used to wrap areas of the React App and catch any errors that occur inside them.
 *
 * Example:- A graph can be wrapped and the title can be set to "Invalid Data" to make sure that the users sees
 * this message instead of a blank screen if an error occurs.
 *
 * This will not affect the dev servers and the errors will still be shown.
 *
 * @returns {React.Component} Error Boundary React Component
 */
class ErrorBoundary extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            hasError: false
        };
    }

    static getDerivedStateFromError = () => ({
        hasError: true
    });

    render = () => {
        const {children, title, description, showNavigationButtons} = this.props;
        const {hasError} = this.state;

        let content;
        if (hasError) {
            content = (
                <UnknownError title={title} description={description} showNavigationButtons={showNavigationButtons}/>
            );
        } else {
            content = children;
        }
        return content;
    };

}

ErrorBoundary.propTypes = {
    children: PropTypes.any.isRequired,
    title: PropTypes.string,
    description: PropTypes.string,
    showNavigationButtons: PropTypes.bool
};

export default ErrorBoundary;
