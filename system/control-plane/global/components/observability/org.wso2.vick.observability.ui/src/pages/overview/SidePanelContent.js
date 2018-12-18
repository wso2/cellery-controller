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

/* eslint react/display-name: "off" */
import "react-vis/dist/style.css";
import "./index.css";
import Avatar from "@material-ui/core/Avatar";
import BlurOnIcon from "@material-ui/icons/BlurOn";
import CellIcon from "@material-ui/icons/CheckBoxOutlineBlank";
import CheckCircleOutline from "@material-ui/icons/CheckCircleOutline";
import ErrorIcon from "@material-ui/icons/ErrorOutline";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import ExpansionPanel from "@material-ui/core/ExpansionPanel";
import ExpansionPanelDetails from "@material-ui/core/ExpansionPanelDetails";
import ExpansionPanelSummary from "@material-ui/core/ExpansionPanelSummary";
import Grain from "@material-ui/icons/Grain";
import Grey from "@material-ui/core/colors/grey";
import HttpTrafficIcon from "@material-ui/icons/CallSplit";
import IconButton from "@material-ui/core/IconButton";
import {Link} from "react-router-dom";
import MUIDataTable from "mui-datatables";
import PropTypes from "prop-types";
import React from "react";
import StateHolder from "../common/state/stateHolder";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Timeline from "@material-ui/icons/Timeline";
import Typography from "@material-ui/core/Typography";
import withGlobalState from "../common/state";
import {withStyles} from "@material-ui/core/styles";
import {
    Hint,
    HorizontalBarSeries,
    HorizontalGridLines,
    VerticalGridLines,
    XAxis,
    XYPlot,
    YAxis
} from "react-vis";
import withColor, {ColorGenerator} from "../common/color";

const styles = (theme) => ({
    drawerContent: {
        padding: 20
    },
    sideBarContentTitle: {
        fontSize: 14,
        fontWeight: 500,
        display: "inline-flex",
        paddingLeft: 5
    },
    titleIcon: {
        verticalAlign: "middle"
    },
    sidebarTableCell: {
        padding: 10
    },
    avatar: {
        width: 25,
        height: 25,
        fontSize: 10,
        fontWeight: 600,
        color: "#fff"
    },
    sidebarContainer: {
        marginBottom: 30
    },
    expansionSum: {
        padding: 0,
        borderBottomWidth: 1,
        borderBottomStyle: "solid",
        borderBottomColor: Grey[300]
    },
    cellIcon: {
        verticalAlign: "middle"
    },
    panel: {
        marginTop: 15,
        boxShadow: "none",
        borderTopWidth: 1,
        borderTopStyle: "solid",
        borderTopColor: Grey[200]
    },
    secondaryHeading: {
        paddingRight: 10
    },
    panelDetails: {
        padding: 0,
        marginBottom: 100
    },
    sidebarListTableText: {
        fontSize: 12

    },
    listIcon: {
        width: 20
    },
    cellNameContainer: {
        marginTop: 10,
        marginBottom: 25
    },
    cellName: {
        display: "inline-flex",
        paddingLeft: 10
    },
    barChart: {
        marginTop: 20
    }
});

class SidePanelContent extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            trafficTooltip: false
        };
    }

    render = () => {
        const {classes, summary, request, isOverview, colorGenerator, globalState, listData} = this.props;
        const {trafficTooltip} = this.state;
        const options = {
            download: false,
            selectableRows: false,
            print: false,
            filter: false,
            search: false,
            viewColumns: false,
            rowHover: false,
            rowsPerPageOptions: false
        };

        const columns = [
            {
                options: {
                    customBodyRender: (value) => {
                        const color = colorGenerator.getColorForPercentage(value, globalState);
                        if (value < globalState.get(StateHolder.CONFIG).percentageRangeMinValue.errorThreshold) {
                            return <ErrorIcon style={{color: color}} className={classes.listIcon}/>;
                        } else if (value
                            < globalState.get(StateHolder.CONFIG).percentageRangeMinValue.warningThreshold) {
                            return <ErrorIcon style={{color: color}} className={classes.listIcon}/>;
                        }
                        return <CheckCircleOutline style={{color: color}} className={classes.listIcon}/>;
                    }
                }
            },
            {
                options: {
                    customBodyRender: (value) => <Typography component={Link} className={classes.sidebarListTableText}
                                                             to={`/cells/${value}`}>{value}</Typography>
                }
            },
            {
                options: {
                    customBodyRender: (value) => (
                        // TODO : Change the to URL specific cell trace search
                        <IconButton size="small" color="inherit" component={Link} to="/tracing/search">
                            <Timeline/>
                        </IconButton>
                    )
                }
            }
        ];

        const BarSeries = HorizontalBarSeries;

        return (
            <div className={classes.drawerContent}>
                <div className={classes.sidebarContainer}>
                    {isOverview === true
                        ? ""
                        : <div className={classes.cellNameContainer}>
                            <CellIcon className={classes.titleIcon}/>
                            <Typography color="inherit"
                                        className={classes.sideBarContentTitle}> Cell:</Typography>
                            {/* TODO : Change the to URL cell value and cell name to selected cell*/}
                            <Typography component={Link} to={"/cells/cell1"} className={classes.cellName}>
                                {summary.topic}</Typography>
                        </div>
                    }
                    <HttpTrafficIcon className={classes.titleIcon}/>
                    <Typography color="inherit" className={classes.sideBarContentTitle}>
                        HTTP Traffic
                    </Typography>

                    <Table className={classes.table}>
                        <TableHead>
                            <TableRow>
                                <TableCell className={classes.sidebarTableCell}>Requests/s</TableCell>
                                {request.statusCodes[1].value === 0
                                    ? ""
                                    : <TableCell className={classes.sidebarTableCell}><Avatar
                                        className={classes.avatar}
                                        style={{
                                            backgroundColor: colorGenerator.getColor(ColorGenerator.SUCCESS)
                                        }}>OK</Avatar></TableCell>}
                                {request.statusCodes[2].value === 0
                                    ? ""
                                    : <TableCell className={classes.sidebarTableCell}><Avatar
                                        className={classes.avatar}
                                        style={{
                                            backgroundColor: colorGenerator.getColor(ColorGenerator.CLIENT_ERROR)
                                        }}>3xx</Avatar></TableCell>}
                                {request.statusCodes[3].value === 0
                                    ? ""
                                    : <TableCell className={classes.sidebarTableCell}><Avatar
                                        className={classes.avatar}
                                        style={{
                                            backgroundColor: colorGenerator.getColor(ColorGenerator.WARNING)
                                        }}>4xx</Avatar></TableCell>}
                                {request.statusCodes[4].value === 0
                                    ? ""
                                    : <TableCell className={classes.sidebarTableCell}><Avatar
                                        className={classes.avatar}
                                        style={{
                                            backgroundColor: colorGenerator.getColor(ColorGenerator.ERROR)
                                        }}>5xx</Avatar></TableCell>}
                                {request.statusCodes[5].value === 0
                                    ? ""
                                    : <TableCell className={classes.sidebarTableCell}><Avatar
                                        className={classes.avatar}
                                        style={{
                                            backgroundColor: colorGenerator.getColor(ColorGenerator.UNKNOWN)
                                        }}>xxx</Avatar></TableCell>}
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            <TableRow>
                                <TableCell
                                    className={classes.sidebarTableCell}>
                                    {request.statusCodes[0].value}</TableCell>
                                {request.statusCodes[1].value === 0
                                    ? ""
                                    : <TableCell className={classes.sidebarTableCell}>
                                        {request.statusCodes[1].value}%</TableCell>}
                                {request.statusCodes[2].value === 0
                                    ? ""
                                    : <TableCell
                                        className={classes.sidebarTableCell}>
                                        {request.statusCodes[2].value}%</TableCell>}
                                {request.statusCodes[3].value === 0
                                    ? ""
                                    : <TableCell
                                        className={classes.sidebarTableCell}>
                                        {request.statusCodes[3].value}%</TableCell>}
                                {request.statusCodes[4].value === 0
                                    ? ""
                                    : <TableCell
                                        className={classes.sidebarTableCell}>
                                        {request.statusCodes[4].value}%</TableCell>}
                                {request.statusCodes[5].value === 0
                                    ? ""
                                    : <TableCell
                                        className={classes.sidebarTableCell}>
                                        {request.statusCodes[5].value}%</TableCell>}
                            </TableRow>
                        </TableBody>
                    </Table>
                    <div className={classes.barChart}>
                        <XYPlot
                            yType="ordinal"
                            stackBy="x"
                            width={250}
                            height={isOverview === true
                                ? 80
                                : 100}>
                            <VerticalGridLines/>
                            <HorizontalGridLines/>
                            <XAxis/>
                            <YAxis/>
                            <BarSeries
                                color={colorGenerator.getColor(ColorGenerator.SUCCESS)}
                                data={[
                                    {
                                        y: "Total", x: request.statusCodes[1].value, title: request.statusCodes[1].key,
                                        percentage: request.statusCodes[1].value, count: request.statusCodes[1].count
                                    }
                                ]}
                                onValueMouseOver={(v) => this.setState({trafficTooltip: v})}
                                onSeriesMouseOut={(v) => this.setState({trafficTooltip: false})}
                            />
                            <BarSeries
                                color={colorGenerator.getColor(ColorGenerator.CLIENT_ERROR)}
                                data={[
                                    {
                                        y: "Total", x: request.statusCodes[2].value, title: request.statusCodes[2].key,
                                        percentage: request.statusCodes[2].value, count: request.statusCodes[2].count
                                    }
                                ]}
                                onValueMouseOver={(v) => this.setState({trafficTooltip: v})}
                                onSeriesMouseOut={(v) => this.setState({trafficTooltip: false})}
                            />
                            <BarSeries
                                color={colorGenerator.getColor(ColorGenerator.WARNING)}
                                data={[
                                    {
                                        y: "Total", x: request.statusCodes[3].value, title: request.statusCodes[3].key,
                                        percentage: request.statusCodes[3].value, count: request.statusCodes[3].count
                                    }
                                ]}
                                onValueMouseOver={(v) => this.setState({trafficTooltip: v})}
                                onSeriesMouseOut={(v) => this.setState({trafficTooltip: false})}
                            />
                            <BarSeries
                                color={colorGenerator.getColor(ColorGenerator.ERROR)}
                                data={[
                                    {
                                        y: "Total", x: request.statusCodes[4].value, title: request.statusCodes[4].key,
                                        percentage: request.statusCodes[4].value, count: request.statusCodes[4].count
                                    }
                                ]}
                                onValueMouseOver={(v) => this.setState({trafficTooltip: v})}
                                onSeriesMouseOut={(v) => this.setState({trafficTooltip: false})}
                            />
                            <BarSeries
                                color={colorGenerator.getColor(ColorGenerator.UNKNOWN)}
                                data={[
                                    {
                                        y: "Total", x: request.statusCodes[5].value, title: request.statusCodes[5].key,
                                        percentage: request.statusCodes[5].value, count: request.statusCodes[5].count
                                    }
                                ]}
                                onValueMouseOver={(v) => this.setState({trafficTooltip: v})}
                                onSeriesMouseOut={(v) => this.setState({trafficTooltip: false})}
                            />
                            {trafficTooltip && <Hint value={trafficTooltip}>
                                <div className="rv-hint__content">
                                    {`${trafficTooltip.title} : ${trafficTooltip.percentage}% (${trafficTooltip.count})`}
                                </div>
                            </Hint>}
                        </XYPlot>
                    </div>
                </div>
                <div className={classes.sidebarContainer}>
                    {isOverview === true
                        ? <Grain className={classes.titleIcon}/>
                        : <BlurOnIcon className={classes.titleIcon}/>}
                    <Typography color="inherit" className={classes.sideBarContentTitle}>
                        {isOverview === true
                            ? "Cells"
                            : "Microservices"} ({summary.content[0].value})
                    </Typography>
                    <ExpansionPanel className={classes.panel}>
                        <ExpansionPanelSummary expandIcon={<ExpandMoreIcon/>}
                                               className={classes.expansionSum}>
                            {summary.content[1].value === 0
                                ? ""
                                : <Typography className={classes.secondaryHeading}><CheckCircleOutline
                                    className={classes.cellIcon}
                                    style={{
                                        color: colorGenerator.getColor(ColorGenerator.SUCCESS)
                                    }}/> {summary.content[1].value}
                                </Typography>}
                            {summary.content[2].value === 0
                                ? ""
                                : <Typography className={classes.secondaryHeading}><ErrorIcon
                                    className={classes.cellIcon}
                                    style={{
                                        color: colorGenerator.getColor(ColorGenerator.ERROR)
                                    }}/> {summary.content[2].value}
                                </Typography>}
                        </ExpansionPanelSummary>
                        <ExpansionPanelDetails className={classes.panelDetails}>
                            <div className="overviewSidebarListTable">
                                <MUIDataTable columns={columns} data={listData} options={options}/>
                            </div>
                        </ExpansionPanelDetails>
                    </ExpansionPanel>
                </div>
            </div>
        );
    }

}

SidePanelContent.propTypes = {
    classes: PropTypes.object.isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator).isRequired,
    summary: PropTypes.object.isRequired,
    request: PropTypes.object.isRequired,
    isOverview: PropTypes.bool.isRequired,
    listData: PropTypes.arrayOf(PropTypes.any).isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired
};

export default withStyles(styles, {withTheme: true})(withGlobalState(withColor(SidePanelContent)));
