/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
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

import Checkbox from "@material-ui/core/Checkbox";
import Constants from "../../utils/constants";
import FormControl from "@material-ui/core/FormControl/FormControl";
import Grid from "@material-ui/core/Grid";
import Input from "@material-ui/core/Input/Input";
import InputLabel from "@material-ui/core/InputLabel/InputLabel";
import ListItemText from "@material-ui/core/ListItemText";
import MenuItem from "@material-ui/core/MenuItem/MenuItem";
import PropTypes from "prop-types";
import React from "react";
import Select from "@material-ui/core/Select/Select";
import Span from "../../utils/span";
import TimelineView from "./TimelineView";
import TracingUtils from "../../utils/tracingUtils";
import classNames from "classnames";
import withStyles from "@material-ui/core/styles/withStyles";

const styles = (theme) => ({
    formControl: {
        margin: theme.spacing.unit
    }
});

class Timeline extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            selectedServiceTypes: [
                Constants.Span.ComponentType.VICK,
                Constants.Span.ComponentType.ISTIO
            ]
        };

        this.handleServiceTypeChange = this.handleServiceTypeChange.bind(this);
    }


    handleServiceTypeChange(event) {
        this.setState({
            selectedServiceTypes: event.target.value
        });
    }

    /**
     * Get the filtered spans from the available spans.
     *
     * @returns {Array.<Span>} The filtered list of spans
     */
    getFilteredSpans() {
        const spans = [];
        for (let i = 0; i < this.props.spans.length; i++) {
            spans.push(this.props.spans[i].shallowClone());
        }
        TracingUtils.buildTree(spans);

        const filteredSpans = [];
        for (let i = 0; i < spans.length; i++) {
            const span = spans[i];
            let isSelected = this.state.selectedServiceTypes.length === 0;

            // Apply service type filter
            isSelected = this.state.selectedServiceTypes.includes(span.componentType)
                || span.componentType === Constants.Span.ComponentType.MICROSERVICE;

            if (isSelected) {
                filteredSpans.push(span);
            } else {
                // Remove the span from the tree without harming the tree structure
                TracingUtils.removeSpanFromTree(span);
            }
        }
        return filteredSpans;
    }

    render() {
        const {classes} = this.props;

        // Finding the service types to be shown in the filter
        const serviceTypes = [];
        for (const filterName in Constants.Span.ComponentType) {
            if (Constants.Span.ComponentType.hasOwnProperty(filterName)) {
                const serviceType = Constants.Span.ComponentType[filterName];
                if (serviceType !== Constants.Span.ComponentType.MICROSERVICE) {
                    serviceTypes.push(serviceType);
                }
            }
        }

        return (
            <div>
                <Grid container justify={"flex-start"} spacing={24}>
                    <Grid item xs={3}>
                        <FormControl className={classNames(classes.formControl)} fullWidth={true}>
                            <InputLabel htmlFor="select-multiple-checkbox">Type</InputLabel>
                            <Select multiple value={this.state.selectedServiceTypes}
                                onChange={this.handleServiceTypeChange}
                                input={<Input id="select-multiple-checkbox"/>}
                                renderValue={(selected) => selected.join(", ")}
                                MenuProps={{
                                    PaperProps: {
                                        style: {
                                            maxHeight: 48 * 4.5 + 8,
                                            width: 250
                                        }
                                    }
                                }}
                            >
                                {
                                    serviceTypes.map((serviceType) => {
                                        const checked = this.state.selectedServiceTypes.indexOf(serviceType) > -1;
                                        return (
                                            <MenuItem key={serviceType} value={serviceType}>
                                                <Checkbox checked={checked}/>
                                                <ListItemText primary={serviceType}/>
                                            </MenuItem>
                                        );
                                    })
                                }
                                <MenuItem key={Constants.Span.ComponentType.MICROSERVICE}
                                    value={Constants.Span.ComponentType.MICROSERVICE} style={{pointerEvents: "none"}}>
                                    <Checkbox checked={true} disabled={true}/>
                                    <ListItemText primary={Constants.Span.ComponentType.MICROSERVICE}/>
                                </MenuItem>
                            </Select>
                        </FormControl>
                    </Grid>
                </Grid>
                <TimelineView spans={this.getFilteredSpans()}/>
            </div>
        );
    }

}

Timeline.propTypes = {
    classes: PropTypes.object.isRequired,
    spans: PropTypes.arrayOf(
        PropTypes.instanceOf(Span).isRequired
    ).isRequired
};

export default withStyles(styles)(Timeline);
