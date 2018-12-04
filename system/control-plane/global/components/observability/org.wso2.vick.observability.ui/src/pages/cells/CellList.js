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

import CheckCircleOutline from "@material-ui/icons/CheckCircleOutline";
import ColorGenerator from "../common/color/colorGenerator";
import DataTable from "../common/DataTable";
import {Link} from "react-router-dom";
import NotificationUtils from "../common/utils/notificationUtils";
import Paper from "@material-ui/core/Paper";
import PropTypes from "prop-types";
import React from "react";
import StateHolder from "../common/state/stateHolder";
import TopToolbar from "../common/toptoolbar";
import withColor from "../common/color";
import withGlobalState from "../common/state";
import {withStyles} from "@material-ui/core/styles";

const styles = (theme) => ({
    root: {
        margin: theme.spacing.unit
    }
});

class CellList extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            data: []
        };
    }

    loadCellInfo = (isUserAction) => {
        const {globalState} = this.props;
        const self = this;

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Cell Info", globalState);
        }
        // TODO : Change to a backend call to fetch data.
        setTimeout(() => {
            const cellData = [
                [1, "Cell A", 0, 0.1, 800, 10],
                [0.8, "Cell B", 0.2, 0.01, 800, 10],
                [0.6, "Cell C", 0.3, 0.2, 800, 10]
            ];
            self.setState({
                data: cellData
            });
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }, 1000);
    };

    render = () => {
        const {classes, colorGenerator, globalState, match} = this.props;
        const {data} = this.state;
        const columns = [
            {
                name: "Health",
                options: {
                    customBodyRender: (value) => {
                        const color = colorGenerator.getColorForPercentage(value, globalState);
                        return <CheckCircleOutline style={{color: color}}/>;
                    }
                }
            },
            {
                name: "Cell",
                options: {
                    customBodyRender: (value) => <Link to={`${match.url}/${value}`}>{value}</Link>
                }
            },
            {
                name: "Inbound Error Rate",
                option: {
                    customBodyRender: (value) => `${value * 100} %`
                }
            },
            {
                name: "Outbound Error Rate",
                option: {
                    customBodyRender: (value) => `${value * 100} %`
                }
            },
            {
                name: "Average Response Time (ms)"
            },
            {
                name: "Average Request Count (requests/s)"
            }
        ];
        return (
            <React.Fragment>
                <TopToolbar title={"Cells"} onUpdate={this.loadCellInfo}/>
                <Paper className={classes.root}>
                    <DataTable columns={columns} data={data}/>
                </Paper>
            </React.Fragment>
        );
    };

}

CellList.propTypes = {
    classes: PropTypes.object.isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator).isRequired,
    match: PropTypes.shape({
        url: PropTypes.string.isRequired
    }).isRequired
};

export default withStyles(styles)(withGlobalState(withColor(CellList)));
