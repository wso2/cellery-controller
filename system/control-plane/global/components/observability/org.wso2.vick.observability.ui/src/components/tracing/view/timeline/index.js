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

import Checkbox from "@material-ui/core/Checkbox";
import Constants from "../../../../utils/constants";
import FormControl from "@material-ui/core/FormControl/FormControl";
import Grid from "@material-ui/core/Grid";
import Input from "@material-ui/core/Input/Input";
import InputLabel from "@material-ui/core/InputLabel/InputLabel";
import ListItemText from "@material-ui/core/ListItemText";
import MenuItem from "@material-ui/core/MenuItem/MenuItem";
import PropTypes from "prop-types";
import React from "react";
import Select from "@material-ui/core/Select/Select";
import Span from "../../../../utils/tracing/span";
import TimelineView from "./TimelineView";
import TracingUtils from "../../../../utils/tracing/tracingUtils";
import withStyles from "@material-ui/core/styles/withStyles";

const styles = (theme) => ({
    formControl: {
        marginTop: theme.spacing.unit * 4,
        marginBottom: theme.spacing.unit * 0.5
    },
    componentTypeMenuItem: {
        pointerEvents: "none"
    }
});

class Timeline extends React.PureComponent {

    constructor(props) {
        super(props);

        this.state = {
            selectedServiceTypes: [
                Constants.CelleryType.COMPONENT,
                Constants.CelleryType.SYSTEM
            ],
            filteredSpans: []
        };

        this.timelineViewRef = React.createRef();
    }

    componentDidMount = () => {
        this.timelineViewRef.current.drawTimeline();
    };

    handleServiceTypeChange = (event) => {
        const serviceType = event.target.value;
        this.setState({
            selectedServiceTypes: serviceType
        });
    };

    static getDerivedStateFromProps = (props, state) => {
        const spans = [];
        for (let i = 0; i < props.spans.length; i++) {
            spans.push(props.spans[i].shallowClone());
        }
        if (spans.length > 0) {
            TracingUtils.buildTree(spans);
        }

        const filteredSpans = [];
        if (spans.length > 0) {
            for (let i = 0; i < spans.length; i++) {
                const span = spans[i];

                // Apply service type filter
                const isSelected = state.selectedServiceTypes.includes(span.componentType);

                if (isSelected) {
                    filteredSpans.push(span);
                } else {
                    // Remove the span from the tree without harming the tree structure
                    TracingUtils.removeSpanFromTree(span);
                }
            }
        }
        if (filteredSpans.length > 0) {
            const filteredTree = TracingUtils.getTreeRoot(filteredSpans);
            TracingUtils.labelSpanTree(filteredTree);
        }
        return {
            ...state,
            filteredSpans: filteredSpans
        };
    };

    render = () => {
        const {classes, selectedComponent} = this.props;
        const {filteredSpans} = this.state;

        // Finding the service types to be shown in the filter
        const serviceTypes = [];
        for (const filterName in Constants.CelleryType) {
            if (Constants.CelleryType.hasOwnProperty(filterName)) {
                const serviceType = Constants.CelleryType[filterName];
                if (serviceType !== Constants.CelleryType.COMPONENT) {
                    serviceTypes.push(serviceType);
                }
            }
        }

        return (
            <React.Fragment>
                <Grid container justify={"flex-start"}>
                    <Grid item xs={3}>
                        <FormControl className={classes.formControl}>
                            <InputLabel htmlFor="select-multiple-checkbox">Type</InputLabel>
                            <Select multiple value={this.state.selectedServiceTypes}
                                onChange={this.handleServiceTypeChange}
                                input={<Input id="select-multiple-checkbox"/>}
                                renderValue={(selected) => selected.join(", ")}>
                                {
                                    serviceTypes.map((serviceType) => {
                                        const checked = this.state.selectedServiceTypes
                                            .filter((type) => type !== Constants.CelleryType.COMPONENT)
                                            .indexOf(serviceType) > -1;
                                        return (
                                            <MenuItem key={serviceType} value={serviceType}>
                                                <Checkbox checked={checked}/>
                                                <ListItemText primary={serviceType}/>
                                            </MenuItem>
                                        );
                                    })
                                }
                                <MenuItem key={Constants.CelleryType.COMPONENT}
                                    value={Constants.CelleryType.COMPONENT}
                                    className={classes.componentTypeMenuItem}>
                                    <Checkbox checked={true}/>
                                    <ListItemText primary={Constants.CelleryType.COMPONENT}/>
                                </MenuItem>
                            </Select>
                        </FormControl>
                    </Grid>
                </Grid>
                <TimelineView spans={filteredSpans} selectedComponent={selectedComponent}
                    innerRef={this.timelineViewRef}/>
            </React.Fragment>
        );
    };

}

Timeline.propTypes = {
    classes: PropTypes.object.isRequired,
    spans: PropTypes.arrayOf(
        PropTypes.instanceOf(Span).isRequired
    ).isRequired,
    clicked: PropTypes.bool,
    selectedComponent: PropTypes.shape({
        cellName: PropTypes.string.isRequired,
        serviceName: PropTypes.string.isRequired
    })
};

export default withStyles(styles)(Timeline);
