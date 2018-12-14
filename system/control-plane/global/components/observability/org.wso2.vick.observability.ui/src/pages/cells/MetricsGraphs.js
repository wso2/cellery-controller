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

import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";
import CardHeader from "@material-ui/core/CardHeader";
import Constants from "../common/constants";
import Grid from "@material-ui/core/Grid";
import PropTypes from "prop-types";
import React from "react";
import Typography from "@material-ui/core/Typography";
import moment from "moment";
import {withStyles} from "@material-ui/core/styles";
import {
    Crosshair,
    DiscreteColorLegend,
    Hint,
    HorizontalBarSeries,
    HorizontalGridLines,
    LineSeries,
    RadialChart,
    VerticalGridLines,
    XAxis,
    XYPlot,
    YAxis,
    makeWidthFlexible
} from "react-vis";
import withColor, {ColorGenerator} from "../common/color";

const styles = {
    root: {
        flexGrow: 1
    },
    card: {
        boxShadow: "none",
        border: "1px solid #eee"
    },
    pos: {
        marginBottom: 12
    },
    cardHeader: {
        borderBottom: "1px solid #eee"
    }, title: {
        fontSize: 16,
        fontWeight: 500,
        color: "#4d4d4d"
    },
    cardDetails: {
        fontSize: 28,
        display: "inline-block"
    },
    cardDetailsSecondary: {
        fontSize: 16,
        display: "inline-block",
        paddingLeft: 5
    },
    contentGrid: {
        height: 186
    },
    toolTipHead: {
        fontWeight: 600
    },
    barChart: {
        marginTop: 40,
        marginBottom: 40
    }
};

const FlexibleWidthXYPlot = makeWidthFlexible(XYPlot);

class MetricsGraphs extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            statusTooltip: false,
            trafficTooltip: false,
            sizeTooltip: [],
            volumeTooltip: false,
            durationTooltip: false
        };
    }

    render = () => {
        const {classes, colorGenerator} = this.props;
        const {statusTooltip, trafficTooltip, volumeTooltip, durationTooltip, sizeTooltip} = this.state;
        const successColor = colorGenerator.getColor(ColorGenerator.SUCCESS);
        const errColor = colorGenerator.getColor(ColorGenerator.ERROR);
        const warningColor = colorGenerator.getColor(ColorGenerator.WARNING);
        const redirectionColor = colorGenerator.getColor(ColorGenerator.REDIRECTION);
        const statusCodesColors = [successColor, redirectionColor, warningColor, errColor];
        const reqResColors = ["#5bbd5a", "#76c7e3"];
        const BarSeries = HorizontalBarSeries;
        const dateTimeFormat = Constants.Pattern.DATE_TIME;

        // TODO: Remove when replaced with actual data
        const MSEC_DAILY = 86400000;
        const timestamp = new Date("December 9 2018").getTime();

        const statusData = [
            {title: "Error", count: 20, percentage: 2, color: errColor},
            {title: "Success", count: 80, percentage: 8, color: successColor}
        ];

        const trafficData = [
            {y: "Out", x: 70, title: "OK", count: 700, percentage: 70},
            {y: "Out", x: 20, title: "3xx", count: 200, percentage: 20},
            {y: "Out", x: 5, title: "4xx", count: 50, percentage: 5},
            {y: "Out", x: 5, title: "5xx", count: 50, percentage: 5}
        ];

        const reqVolumeData = [
            {x: timestamp + MSEC_DAILY, y: 3},
            {x: timestamp + MSEC_DAILY * 2, y: 5},
            {x: timestamp + MSEC_DAILY * 3, y: 15},
            {x: timestamp + MSEC_DAILY * 4, y: 10},
            {x: timestamp + MSEC_DAILY * 5, y: 6},
            {x: timestamp + MSEC_DAILY * 6, y: 3},
            {x: timestamp + MSEC_DAILY * 7, y: 9},
            {x: timestamp + MSEC_DAILY * 8, y: 11}
        ];
        const reqDurationData = [
            {x: timestamp + MSEC_DAILY, y: 3},
            {x: timestamp + MSEC_DAILY * 2, y: 5},
            {x: timestamp + MSEC_DAILY * 3, y: 15},
            {x: timestamp + MSEC_DAILY * 4, y: 10},
            {x: timestamp + MSEC_DAILY * 5, y: 6},
            {x: timestamp + MSEC_DAILY * 6, y: 3},
            {x: timestamp + MSEC_DAILY * 7, y: 9},
            {x: timestamp + MSEC_DAILY * 8, y: 11}
        ];

        const reqResSizeData = [
            {
                name: "Request",
                data: [
                    {x: timestamp + MSEC_DAILY, y: 3},
                    {x: timestamp + MSEC_DAILY * 2, y: 5},
                    {x: timestamp + MSEC_DAILY * 3, y: 15},
                    {x: timestamp + MSEC_DAILY * 4, y: 10},
                    {x: timestamp + MSEC_DAILY * 5, y: 6},
                    {x: timestamp + MSEC_DAILY * 6, y: 3},
                    {x: timestamp + MSEC_DAILY * 7, y: 9},
                    {x: timestamp + MSEC_DAILY * 8, y: 11}
                ]
            },
            {
                name: "Response",
                data: [
                    {x: timestamp + MSEC_DAILY, y: 10},
                    {x: timestamp + MSEC_DAILY * 2, y: 4},
                    {x: timestamp + MSEC_DAILY * 3, y: 2},
                    {x: timestamp + MSEC_DAILY * 4, y: 15},
                    {x: timestamp + MSEC_DAILY * 5, y: 13},
                    {x: timestamp + MSEC_DAILY * 6, y: 6},
                    {x: timestamp + MSEC_DAILY * 7, y: 7},
                    {x: timestamp + MSEC_DAILY * 8, y: 2}
                ]
            }
        ];

        return (
            <div className={classes.root}>
                <Grid container spacing={24}>
                    <Grid item md={3} sm={6} xs={12}>
                        <Card className={classes.card}>
                            <CardHeader
                                classes={{
                                    title: classes.title
                                }}
                                title="Success/ Failure Rate"
                                className={classes.cardHeader}
                            />
                            <CardContent className={classes.content} align="center">
                                <RadialChart
                                    data={statusData}
                                    innerRadius={60}
                                    radius={85}
                                    getAngle={(d) => d.percentage}
                                    onValueMouseOver={(v) => this.setState({statusTooltip: v})}
                                    onSeriesMouseOut={(v) => this.setState({statusTooltip: false})}
                                    width={180}
                                    height={180}
                                    colorType="literal"
                                >
                                    {statusTooltip && <Hint value={statusTooltip}>
                                        <div className="rv-hint__content">
                                            {`${statusTooltip.title} :
                                            ${statusTooltip.percentage}% (${statusTooltip.count} Requests)`}
                                        </div>
                                    </Hint>}
                                </RadialChart>
                                <div>
                                    <DiscreteColorLegend items={statusData.map((d) => d.title)}
                                        colors={[errColor, successColor]}
                                        orientation="horizontal"/>
                                </div>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item md={3} sm={6} xs={12} alignItems="center">
                        <Grid item sm={12} className={classes.contentGrid}>
                            <Card className={classes.card}>
                                <CardHeader
                                    classes={{
                                        title: classes.title
                                    }}
                                    title="Average Response Time"
                                    className={classes.cardHeader}
                                />
                                <CardContent align="center">
                                    <Typography className={classes.cardDetails}>
                                        200
                                    </Typography>
                                    <Typography color="textSecondary" className={classes.cardDetailsSecondary}>
                                        ms
                                    </Typography>
                                </CardContent>
                            </Card>
                        </Grid>
                        <Grid item sm={12} className={classes.contentGrid}>
                            <Card className={classes.card}>
                                <CardHeader
                                    classes={{
                                        title: classes.title
                                    }}
                                    title="Average Request Count"
                                    className={classes.cardHeader}
                                />
                                <CardContent align="center">
                                    <Typography className={classes.cardDetails}>
                                        25
                                    </Typography>
                                    <Typography color="textSecondary" className={classes.cardDetailsSecondary}>
                                        Requests/s
                                    </Typography>
                                </CardContent>
                            </Card>
                        </Grid>
                    </Grid>
                    <Grid item md={6} sm={12} xs={12}>
                        <Card className={classes.card}>
                            <CardHeader
                                classes={{
                                    title: classes.title
                                }}
                                title="HTTP Traffic"
                                className={classes.cardHeader}
                            />
                            <CardContent className={classes.content}>
                                <div>
                                    <FlexibleWidthXYPlot
                                        yType="ordinal"
                                        stackBy="x"
                                        height={100}
                                        className={classes.barChart}>
                                        <VerticalGridLines/>
                                        <HorizontalGridLines/>
                                        <XAxis/>
                                        <YAxis/>

                                        {
                                            trafficData.map((dataItem, index) => (
                                                <BarSeries
                                                    key={dataItem.y}
                                                    data={[dataItem]}
                                                    color={statusCodesColors[index]}
                                                    onValueMouseOver={(v) => this.setState({trafficTooltip: v})}
                                                    onSeriesMouseOut={(v) => this.setState({trafficTooltip: false})}
                                                />
                                            ))

                                        }
                                        {trafficTooltip && <Hint value={trafficTooltip}>
                                            <div className="rv-hint__content">
                                                {`${trafficTooltip.title} :
                                                ${trafficTooltip.percentage}% (${trafficTooltip.count} Requests)`}
                                            </div>
                                        </Hint>}
                                    </FlexibleWidthXYPlot>
                                    <DiscreteColorLegend
                                        orientation="horizontal"
                                        items={[
                                            {
                                                title: "OK",
                                                color: successColor
                                            },
                                            {
                                                title: "3xx",
                                                color: warningColor
                                            },
                                            {
                                                title: "4xx",
                                                color: redirectionColor
                                            },
                                            {
                                                title: "5xx",
                                                color: errColor
                                            }
                                        ]}
                                    />
                                </div>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item md={6} sm={12} xs={12}>
                        <Card className={classes.card}>
                            <CardHeader
                                classes={{
                                    title: classes.title
                                }}
                                title="Request Volume"
                                className={classes.cardHeader}
                            />
                            <CardContent className={classes.content}>
                                <div className={classes.lineChart}>
                                    <FlexibleWidthXYPlot xType="time" height={400}
                                        onMouseLeave={() => this.setState({volumeTooltip: false})}>
                                        <HorizontalGridLines/>
                                        <VerticalGridLines/>
                                        <XAxis title="Time"/>
                                        <YAxis title="Volume(ops)"/>
                                        <LineSeries
                                            data={reqVolumeData}
                                            onNearestX={(d) => this.setState({volumeTooltip: d})}
                                        />
                                        {volumeTooltip && <Crosshair values={[volumeTooltip]}>
                                            <div className="rv-hint__content">
                                                {`${moment(volumeTooltip.x).format(Constants.Pattern.DATE_TIME)} :
                                                ${volumeTooltip.y} Requests`}</div>
                                        </Crosshair>}
                                    </FlexibleWidthXYPlot>
                                    <DiscreteColorLegend
                                        orientation="horizontal"
                                        items={[
                                            {
                                                title: "Request",
                                                color: "#12939a"
                                            }
                                        ]}
                                    />
                                </div>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item md={6} sm={12} xs={12}>
                        <Card className={classes.card}>
                            <CardHeader
                                classes={{
                                    title: classes.title
                                }}
                                title="Request Duration"
                                className={classes.cardHeader}
                            />
                            <CardContent className={classes.content}>
                                <div className={classes.lineChart}>
                                    <FlexibleWidthXYPlot xType="time" height={400}
                                        onMouseLeave={() => this.setState({durationTooltip: false})}>
                                        <HorizontalGridLines/>
                                        <VerticalGridLines/>
                                        <XAxis title="Time"/>
                                        <YAxis title="Duration(s)"/>
                                        <LineSeries color="#3f51b5"
                                            data={reqDurationData}
                                            onNearestX={(d) => this.setState({durationTooltip: d})}

                                        />
                                        {durationTooltip && <Crosshair values={[durationTooltip]}>
                                            <div className="rv-hint__content">
                                                {`${moment(durationTooltip.x).format(Constants.Pattern.DATE_TIME)} :
                                                ${durationTooltip.y}`}</div>
                                        </Crosshair>}
                                    </FlexibleWidthXYPlot>
                                    <DiscreteColorLegend
                                        orientation="horizontal"
                                        items={[
                                            {
                                                title: "Request",
                                                color: "#3f51b5"
                                            }
                                        ]}
                                    />
                                </div>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item md={6} sm={12} xs={12}>
                        <Card className={classes.card}>
                            <CardHeader
                                classes={{
                                    title: classes.title
                                }}
                                title="Request/Response Size"
                                className={classes.cardHeader}
                            />
                            <CardContent className={classes.content}>
                                <div>
                                    <FlexibleWidthXYPlot xType="time" height={400}
                                        onMouseLeave={() => this.setState({sizeTooltip: []})}>
                                        <HorizontalGridLines/>
                                        <VerticalGridLines/>
                                        <XAxis title="Time"/>
                                        <YAxis title="Size"/>
                                        {
                                            reqResSizeData.map((dataItem, index) => (
                                                <LineSeries
                                                    key={dataItem.name}
                                                    data={dataItem.data}
                                                    onNearestX={(d, {index}) => this.setState({
                                                        sizeTooltip: reqResSizeData.map((d) => ({
                                                            ...d.data[index],
                                                            name: d.name
                                                        }))
                                                    })}
                                                    color={reqResColors[index]}
                                                />))
                                        }

                                        <Crosshair values={sizeTooltip}>
                                            {
                                                sizeTooltip.length > 0
                                                    ? (
                                                        <div className="rv-hint__content">
                                                            <div className={classes.toolTipHead}>
                                                                {`${moment(sizeTooltip[0].x).format(dateTimeFormat)}`}
                                                            </div>
                                                            {sizeTooltip.map(
                                                                (tooltipItem, {index}) => <div key={tooltipItem.name}>
                                                                    {`${tooltipItem.name}: ${tooltipItem.y}`}</div>)
                                                            }
                                                        </div>
                                                    )
                                                    : null
                                            }
                                        </Crosshair>
                                    </FlexibleWidthXYPlot>
                                    <DiscreteColorLegend
                                        orientation="horizontal"
                                        items={
                                            reqResSizeData.map((d, index) => ({
                                                title: d.name,
                                                color: reqResColors[index]
                                            }))}
                                    />
                                </div>
                            </CardContent>
                        </Card>
                    </Grid>
                </Grid>
            </div>
        );
    }

}

MetricsGraphs.propTypes = {
    classes: PropTypes.object.isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator).isRequired
};

export default withStyles(styles)(withColor(MetricsGraphs));
