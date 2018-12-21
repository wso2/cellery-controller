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
/* eslint max-len: ["off"] */

import "react-vis/dist/style.css";
import "./index.css";
import Avatar from "@material-ui/core/Avatar";
import CheckCircleOutline from "@material-ui/icons/CheckCircleOutline";
import ErrorIcon from "@material-ui/icons/ErrorOutline";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import ExpansionPanel from "@material-ui/core/ExpansionPanel";
import ExpansionPanelDetails from "@material-ui/core/ExpansionPanelDetails";
import ExpansionPanelSummary from "@material-ui/core/ExpansionPanelSummary";
import Grey from "@material-ui/core/colors/grey";
import IconButton from "@material-ui/core/IconButton";
import {Link} from "react-router-dom";
import MUIDataTable from "mui-datatables";
import PropTypes from "prop-types";
import React from "react";
import StateHolder from "../common/state/stateHolder";
import SvgIcon from "@material-ui/core/SvgIcon";
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
        paddingLeft: 10
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

const CellIcon = (props) => (
    <SvgIcon viewBox="0 0 14 14" {...props}>
        <path d="M7,0.1L1,3.5l0,6.9l6,3.5l6-3.4l0-6.9L7,0.1z M12,9.9l-5,2.9L2,9.9l0-5.7l5-2.9l5,2.9L12,9.9z"/>
    </SvgIcon>
);

const MicroserviceIcon = (props) => (
    <SvgIcon viewBox="0 0 14 14" {...props}>
        <path d="M7,5.4C5.6,5.4,4.4,6.5,4.4,8s1.2,2.6,2.6,2.6c1.4,0,2.6-1.2,2.6-2.6S8.5,5.4,7,5.4z M7,9.6C6.2,9.6,5.4,8.8,5.4,8
c0-0.9,0.7-1.6,1.6-1.6c0.9,0,1.6,0.7,1.6,1.6C8.6,8.8,7.9,9.6,7,9.6z M7,3.5c0.4,0,0.8,0.3,0.8,0.8S7.5,5,7,5S6.3,4.7,6.3,4.3
S6.6,3.5,7,3.5z M11,10c-0.2,0.4-0.7,0.5-1,0.3s-0.5-0.7-0.3-1c0.2-0.4,0.7-0.5,1-0.3C11,9.1,11.2,9.6,11,10z M4.5,9.6
c0,0.4-0.3,0.8-0.8,0.8S2.9,10,2.9,9.6s0.3-0.8,0.8-0.8S4.5,9.2,4.5,9.6z M12.4,9.6c-0.8,0-1.5,0.7-1.5,1.5s0.7,1.5,1.5,1.5
s1.5-0.7,1.5-1.5S13.2,9.6,12.4,9.6z M12.4,11.5c-0.3,0-0.5-0.2-0.5-0.5c0-0.3,0.2-0.5,0.5-0.5c0.3,0,0.5,0.2,0.5,0.5
C12.9,11.3,12.7,11.5,12.4,11.5z M1.6,9.6c-0.8,0-1.5,0.7-1.5,1.5s0.7,1.5,1.5,1.5S3,11.9,3,11.1S2.4,9.6,1.6,9.6z M1.6,11.6
c-0.3,0-0.5-0.2-0.5-0.5s0.2-0.5,0.5-0.5S2,10.8,2,11.1S1.8,11.6,1.6,11.6z M7,3.2c0.8,0,1.5-0.7,1.5-1.5S7.8,0.2,7,0.2
S5.6,0.9,5.6,1.7S6.2,3.2,7,3.2z M7,1.2c0.3,0,0.5,0.2,0.5,0.5S7.3,2.2,7,2.2S6.6,2,6.6,1.7S6.8,1.2,7,1.2z M10.2,5.5
c0-0.6,0.5-1.1,1.1-1.1c0.6,0,1.1,0.5,1.1,1.1c0,0.6-0.5,1.1-1.1,1.1C10.7,6.6,10.2,6.1,10.2,5.5z M8.3,12.7c0,0.6-0.5,1.1-1.1,1.1
s-1.1-0.5-1.1-1.1c0-0.6,0.5-1.1,1.1-1.1S8.3,12,8.3,12.7z M1.8,5.5c0-0.6,0.5-1.1,1.1-1.1C3.5,4.4,4,4.9,4,5.5
c0,0.6-0.5,1.1-1.1,1.1C2.3,6.6,1.8,6.1,1.8,5.5z"/>
    </SvgIcon>
);

const HttpTrafficIcon = (props) => (
    <SvgIcon viewBox="0 0 14 14" {...props}>
        <path d="M1.1,6.2V2.9c0-0.2,0.2-0.4,0.4-0.4h8.8V0.2l2.4,2.6c0.1,0.1,0.1,0.3,0,0.4l-2.5,2.6l0-2.3H2.1l0,2.9c0,0.2-0.2,0.4-0.4,0.4
H1.5C1.3,6.6,1.1,6.4,1.1,6.2z M12.5,7.4h-0.3c-0.2,0-0.4,0.2-0.4,0.4l0,2.9H3.7l0-2.3l-2.5,2.6c-0.1,0.1-0.1,0.3,0,0.4l2.4,2.6
v-2.3h8.8c0.2,0,0.4-0.2,0.4-0.4V7.8C12.9,7.6,12.7,7.4,12.5,7.4z"/>
    </SvgIcon>
);

const CellsIcon = (props) => (
    <SvgIcon viewBox="0 0 14 14" {...props}>
        <path d="M10.7,13.7l-3.1-1.8V8.4l3.1-1.8l3.1,1.8v3.6L10.7,13.7z M7,7.5L3.8,5.7V2L7,0.2L10.2,2v3.6L7,7.5z M5.1,4.9 l1.9,1l1.9
-1V 2.8L7,1.8l-1.9,1V4.9z M3.4,13.8L0.2,12V8.4l3.2-1.8l3.2,1.8V12L3.4,13.8z M8.8,11.2l1.9,1l1.9-1V9.1l-1.9 -1l-1.9,
1V11.2z M1.5,11.2l1.9,1l1.9-1V9.1l-1.9-1l-1.9,1V11.2z"/>
    </SvgIcon>
);

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
                            <CellIcon className={classes.titleIcon} fontSize="small"/>
                            <Typography color="inherit"
                                className={classes.sideBarContentTitle}> Cell:</Typography>
                            {/* TODO : Change the to URL cell value and cell name to selected cell*/}
                            <Typography component={Link} to={"/cells/cell1"} className={classes.cellName}>
                                {summary.topic}</Typography>
                        </div>
                    }
                    <HttpTrafficIcon className={classes.titleIcon} fontSize="small"/>
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
                                    {`${trafficTooltip.title} :
                                    ${trafficTooltip.percentage}% (${trafficTooltip.count})`}
                                </div>
                            </Hint>}
                        </XYPlot>
                    </div>
                </div>
                <div className={classes.sidebarContainer}>
                    {isOverview === true
                        ? <CellsIcon className={classes.titleIcon} fontSize="small"/>
                        : <MicroserviceIcon className={classes.titleIcon} fontSize="small"/>}
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
