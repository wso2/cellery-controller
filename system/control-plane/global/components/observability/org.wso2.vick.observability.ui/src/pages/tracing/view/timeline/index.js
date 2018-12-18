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
import Constants from "../../../common/constants";
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
import withStyles from "@material-ui/core/styles/withStyles";

const styles = (theme) => ({
    formControl: {
        margin: theme.spacing.unit
    },
    microserviceTypeMenuItem: {
        pointerEvents: "none"
    }
});

class Timeline extends React.PureComponent {

    constructor(props) {
        super(props);

        this.state = {
            selectedServiceTypes: [
                Constants.ComponentType.MICROSERVICE,
                Constants.ComponentType.VICK,
                Constants.ComponentType.ISTIO
            ],
            filteredSpans: []
        };

        this.timelineViewRef = React.createRef();
    }

    componentDidMount = () => {
        this.timelineViewRef.current.drawTimeline();
    };

    componentDidUpdate = () => {
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
        const {classes, selectedMicroservice} = this.props;
        const {filteredSpans} = this.state;

        // Finding the service types to be shown in the filter
        const serviceTypes = [];
        for (const filterName in Constants.ComponentType) {
            if (Constants.ComponentType.hasOwnProperty(filterName)) {
                const serviceType = Constants.ComponentType[filterName];
                if (serviceType !== Constants.ComponentType.MICROSERVICE) {
                    serviceTypes.push(serviceType);
                }
            }
        }

        return (
            <React.Fragment>
                <Grid container justify={"flex-start"} spacing={24}>
                    <Grid item xs={3}>
                        <FormControl className={classes.formControl} fullWidth={true}>
                            <InputLabel htmlFor="select-multiple-checkbox">Type</InputLabel>
                            <Select multiple value={this.state.selectedServiceTypes}
                                onChange={this.handleServiceTypeChange}
                                input={<Input id="select-multiple-checkbox"/>}
                                renderValue={(selected) => selected.join(", ")}>
                                {
                                    serviceTypes.map((serviceType) => {
                                        const checked = this.state.selectedServiceTypes
                                            .filter((type) => type !== Constants.ComponentType.MICROSERVICE)
                                            .indexOf(serviceType) > -1;
                                        return (
                                            <MenuItem key={serviceType} value={serviceType}>
                                                <Checkbox checked={checked}/>
                                                <ListItemText primary={serviceType}/>
                                            </MenuItem>
                                        );
                                    })
                                }
                                <MenuItem key={Constants.ComponentType.MICROSERVICE}
                                    value={Constants.ComponentType.MICROSERVICE}
                                    className={classes.microserviceTypeMenuItem}>
                                    <Checkbox checked={true}/>
                                    <ListItemText primary={Constants.ComponentType.MICROSERVICE}/>
                                </MenuItem>
                            </Select>
                        </FormControl>
                    </Grid>
                </Grid>
                <TimelineView spans={filteredSpans} selectedMicroservice={selectedMicroservice}
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
    selectedMicroservice: PropTypes.shape({
        cellName: PropTypes.string.isRequired,
        serviceName: PropTypes.string.isRequired
    }).isRequired
};

export default withStyles(styles)(Timeline);
