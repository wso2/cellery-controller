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

/* eslint max-lines: ["off"] */

import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";
import CardHeader from "@material-ui/core/CardHeader";
import Constants from "../common/constants";
import Grid from "@material-ui/core/Grid";
import InfoIcon from "@material-ui/icons/InfoOutlined";
import React from "react";
import Tooltip from "@material-ui/core/Tooltip";
import Typography from "@material-ui/core/Typography";
import moment from "moment";
import {withStyles} from "@material-ui/core/styles";
import {
    Crosshair, DiscreteColorLegend, Highlight, Hint, HorizontalBarSeries, HorizontalGridLines, LineMarkSeries,
    RadialChart, VerticalGridLines, XAxis, XYPlot, YAxis, makeWidthFlexible
} from "react-vis";
import withColor, {ColorGenerator} from "../common/color";
import * as PropTypes from "prop-types";

const styles = {
    root: {
        flexGrow: 1
    },
    card: {
        boxShadow: "none",
        border: "1px solid #eee"
    },
    cardHeader: {
        borderBottom: "1px solid #eee"
    },
    title: {
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
        marginTop: 48
    },
    toolTipHead: {
        fontWeight: 600
    },
    barChart: {
        marginTop: 40,
        marginBottom: 40
    },
    info: {
        color: "#999",
        marginTop: 8,
        marginRight: 27,
        fontSize: 18

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
            durationTooltip: false,
            lastDrawLocationVolumeChart: null,
            lastDrawLocationDurationChart: null,
            lastDrawLocationSizeChart: null
        };
    }

    calculateMetrics = () => {
        const {colorGenerator, data, direction} = this.props;
        const successColor = colorGenerator.getColor(ColorGenerator.SUCCESS);
        const errColor = colorGenerator.getColor(ColorGenerator.ERROR);

        let totalRequestsCount = 0;
        let totalErrorsCount = 0;
        let totalResponseTime = 0;
        let minTime = Number.MAX_SAFE_INTEGER;
        let maxTime = 0;
        const httpResponseGroupCounts = {
            "2xx": 0,
            "3xx": 0,
            "4xx": 0,
            "5xx": 0
        };
        for (const datum of data) {
            totalRequestsCount += datum.requestCount;
            totalResponseTime += datum.totalResponseTimeMilliSec;
            httpResponseGroupCounts[datum.httpResponseGroup] += datum.requestCount;

            if (datum.httpResponseGroup === "5xx") {
                totalErrorsCount += datum.requestCount;
            }

            if (datum.timestamp < minTime) {
                minTime = datum.timestamp;
            }
            if (datum.timestamp > maxTime) {
                maxTime = datum.timestamp;
            }
        }
        const timeRange = maxTime > minTime ? maxTime - minTime : 0;

        // Preparing data for the Success / Failure rate Pie Chart
        const totalErrorPercentage = totalRequestsCount === 0 ? 0 : totalErrorsCount * 100 / totalRequestsCount;
        const totalSuccessPercentage = totalRequestsCount === 0
            ? 0
            : (totalRequestsCount - totalErrorsCount) * 100 / totalRequestsCount;
        const statusData = [];
        if (totalErrorPercentage > 0) {
            statusData.push({
                title: "Error",
                count: totalErrorsCount,
                percentage: totalErrorPercentage,
                color: errColor
            });
        }
        if (totalSuccessPercentage > 0) {
            statusData.push({
                title: "Success",
                count: totalRequestsCount - totalErrorsCount,
                percentage: totalSuccessPercentage,
                color: successColor
            });
        }

        // Preparing data for the HTTP traffic Bar Chart
        const trafficData = ["2xx", "3xx", "4xx", "5xx"]
            .map((datum) => ({
                x: totalRequestsCount === 0 ? 0 : httpResponseGroupCounts[datum] * 100 / totalRequestsCount,
                y: direction,
                title: (datum === "2xx" ? "OK" : datum),
                count: httpResponseGroupCounts[datum]
            }));

        // Aggregating the data by timestamp (time-series charts doesn't need to consider response code)
        const aggregatedData = [];
        let lastTimestamp = 0;
        let aggregatedDatum = null;
        for (let i = 0; i < data.length; i++) {
            const datum = data[i];
            if (datum.timestamp === lastTimestamp) {
                aggregatedDatum.totalRequestSizeBytes += datum.totalRequestSizeBytes;
                aggregatedDatum.totalResponseSizeBytes += datum.totalResponseSizeBytes;
                aggregatedDatum.totalResponseTimeMilliSec += datum.totalResponseTimeMilliSec;
                aggregatedDatum.requestCount += datum.requestCount;
            } else {
                if (aggregatedDatum) {
                    aggregatedData.push(aggregatedDatum);
                }

                lastTimestamp = datum.timestamp;
                aggregatedDatum = {
                    timestamp: datum.timestamp,
                    totalRequestSizeBytes: datum.totalRequestSizeBytes,
                    totalResponseSizeBytes: datum.totalResponseSizeBytes,
                    totalResponseTimeMilliSec: datum.totalResponseTimeMilliSec,
                    requestCount: datum.requestCount
                };
            }
        }

        // Preparing a proper time-series data-set with 0 requests points added
        const timeSeriesData = [];
        const addEmptyTimeSeriesPoint = (timestamp) => {
            timeSeriesData.push({
                timestamp: timestamp,
                totalRequestSizeBytes: 0,
                totalResponseSizeBytes: 0,
                totalResponseTimeMilliSec: 0,
                requestCount: 0
            });
        };
        for (let i = 0; i < aggregatedData.length; i++) {
            const datum = aggregatedData[i];

            if (i > 0 && timeSeriesData[timeSeriesData.length - 1].timestamp < datum.timestamp - 1000) {
                addEmptyTimeSeriesPoint(datum.timestamp - 1000);
            }
            timeSeriesData.push(datum);
            if (i < aggregatedData.length - 1 && aggregatedData[i + 1].timestamp > datum.timestamp + 1000) {
                addEmptyTimeSeriesPoint(datum.timestamp + 1000);
            }
        }

        // Preparing data for the Request Volume Line Chart
        const reqVolumeData = timeSeriesData.map((datum) => ({
            x: datum.timestamp,
            y: datum.requestCount
        }));

        // Preparing data for the Request Duration Line Chart
        const reqDurationData = timeSeriesData.map((datum) => ({
            x: datum.timestamp,
            y: datum.requestCount === 0 ? 0 : datum.totalResponseTimeMilliSec / datum.requestCount
        }));

        // Preparing data for the Response Size Line Chart
        const reqResSizeData = [
            {
                name: "Request",
                data: timeSeriesData.map((datum) => ({
                    x: datum.timestamp,
                    y: datum.requestCount === 0 ? 0 : datum.totalRequestSizeBytes / datum.requestCount
                }))
            },
            {
                name: "Response",
                data: timeSeriesData.map((datum) => ({
                    x: datum.timestamp,
                    y: datum.requestCount === 0 ? 0 : datum.totalResponseSizeBytes / datum.requestCount
                }))
            }
        ];

        return {
            totalResponseTime: totalResponseTime,
            totalRequestsCount: totalRequestsCount,
            timeRange: timeRange,
            statusData: statusData,
            trafficData: trafficData,
            reqVolumeData: reqVolumeData,
            reqDurationData: reqDurationData,
            reqResSizeData: reqResSizeData
        };
    };

    render = () => {
        const {classes, colorGenerator} = this.props;
        const {
            statusTooltip, trafficTooltip, volumeTooltip, durationTooltip, sizeTooltip, lastDrawLocationVolumeChart,
            lastDrawLocationDurationChart, lastDrawLocationSizeChart
        } = this.state;

        const successColor = colorGenerator.getColor(ColorGenerator.SUCCESS);
        const errColor = colorGenerator.getColor(ColorGenerator.ERROR);
        const warningColor = colorGenerator.getColor(ColorGenerator.WARNING);
        const redirectionColor = colorGenerator.getColor(ColorGenerator.CLIENT_ERROR);

        const statusCodesColors = [successColor, redirectionColor, warningColor, errColor];
        const reqResColors = ["#5bbd5a", "#76c7e3"];

        const {
            timeRange, statusData, trafficData, reqVolumeData, reqDurationData, reqResSizeData, totalRequestsCount,
            totalResponseTime
        } = this.calculateMetrics();

        const zoomHintToolTip = (
            <Tooltip title="Click and drag in the plot area to zoom in, click anywhere in the graph to zoom out.">
                <InfoIcon className={classes.info}/>
            </Tooltip>
        );

        return (
            <div className={classes.root}>
                <Grid container spacing={24}>
                    <Grid item md={3} sm={6} xs={12}>
                        <Card className={classes.card}>
                            <CardHeader title="Success / Failure Rate" className={classes.cardHeader}
                                classes={{
                                    title: classes.title
                                }}/>
                            <CardContent className={classes.content} align="center">
                                <RadialChart innerRadius={60} radius={85} width={180} height={180}
                                    colorType="literal"
                                    getAngle={(d) => d.percentage}
                                    onValueMouseOver={(v) => this.setState({statusTooltip: v})}
                                    onSeriesMouseOut={() => this.setState({statusTooltip: false})}
                                    data={statusData}>
                                    {
                                        statusTooltip
                                            ? (
                                                <Hint value={statusTooltip}>
                                                    <div className="rv-hint__content">
                                                        {
                                                            `${statusTooltip.title} :
                                                            ${Math.round(statusTooltip.percentage)}%
                                                            (${statusTooltip.count} Requests)`
                                                        }
                                                    </div>
                                                </Hint>
                                            ) : null
                                    }
                                </RadialChart>
                                <div>
                                    <DiscreteColorLegend items={statusData.map((d) => d.title)}
                                        colors={statusData.map((statusDatum) => statusDatum.color)}
                                        orientation="horizontal"/>
                                </div>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid container item md={3} sm={6} xs={12} alignItems="center">
                        <Grid item sm={12}>
                            <Card className={classes.card}>
                                <CardHeader title="Average Response Time" className={classes.cardHeader}
                                    classes={{
                                        title: classes.title
                                    }}/>
                                <CardContent align="center">
                                    <Typography className={classes.cardDetails}>
                                        {
                                            totalRequestsCount === 0
                                                ? 0
                                                : Math.round(totalResponseTime / totalRequestsCount)
                                        }
                                    </Typography>
                                    <Typography color="textSecondary" className={classes.cardDetailsSecondary}>
                                        ms
                                    </Typography>
                                </CardContent>
                            </Card>
                        </Grid>
                        <Grid item sm={12} className={classes.contentGrid}>
                            <Card className={classes.card}>
                                <CardHeader title="Average Request Count" className={classes.cardHeader}
                                    classes={{
                                        title: classes.title
                                    }}/>
                                <CardContent align="center">
                                    <Typography className={classes.cardDetails}>
                                        {
                                            timeRange === 0
                                                ? 0
                                                : Math.round(totalRequestsCount * 1000 * 100 / timeRange) / 100
                                        }
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
                            <CardHeader title="HTTP Traffic" className={classes.cardHeader}
                                classes={{
                                    title: classes.title
                                }}/>
                            <CardContent className={classes.content}>
                                <div>
                                    <FlexibleWidthXYPlot yType="ordinal" stackBy="x" height={100}
                                        className={classes.barChart}>
                                        <VerticalGridLines/>
                                        <HorizontalGridLines/>
                                        <XAxis/>
                                        <YAxis/>
                                        {
                                            trafficData.map((dataItem, index) => (
                                                <HorizontalBarSeries key={dataItem.title} data={[dataItem]}
                                                    color={statusCodesColors[index]}
                                                    onValueMouseOver={(v) => this.setState({trafficTooltip: v})}
                                                    onSeriesMouseOut={() => this.setState({trafficTooltip: false})}
                                                />
                                            ))
                                        }
                                        {
                                            trafficTooltip
                                                ? (
                                                    <Hint value={trafficTooltip}>
                                                        <div className="rv-hint__content">{
                                                            `${trafficTooltip.title} : ${Math.round(trafficTooltip.x)}%
                                                            (${trafficTooltip.count} Requests)`
                                                        }</div>
                                                    </Hint>
                                                )
                                                : null
                                        }
                                    </FlexibleWidthXYPlot>
                                    <DiscreteColorLegend orientation="horizontal"
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
                            <CardHeader title="Request Volume" className={classes.cardHeader}
                                classes={{
                                    title: classes.title
                                }}
                                action={zoomHintToolTip}
                            />
                            <CardContent className={classes.content}>
                                <div>
                                    <FlexibleWidthXYPlot xType="time" height={400} animation
                                        xDomain={
                                            lastDrawLocationVolumeChart && [
                                                lastDrawLocationVolumeChart.left,
                                                lastDrawLocationVolumeChart.right
                                            ]
                                        }
                                        onMouseLeave={() => this.setState({volumeTooltip: false})}>
                                        <HorizontalGridLines/>
                                        <VerticalGridLines/>
                                        <XAxis title="Time"/>
                                        <YAxis title="Volume (ops / s)"/>
                                        <LineMarkSeries data={reqVolumeData} color="#12939a" size={3}
                                            onNearestX={(d) => this.setState({volumeTooltip: d})}/>
                                        {
                                            volumeTooltip
                                                ? (
                                                    <Crosshair values={[volumeTooltip]}>
                                                        <div className="rv-hint__content">
                                                            {
                                                                `${moment(volumeTooltip.x)
                                                                    .format(Constants.Pattern.DATE_TIME)} :
                                                                ${Math.round(volumeTooltip.y)} Requests`
                                                            }
                                                        </div>
                                                    </Crosshair>
                                                )
                                                : null
                                        }
                                        <Highlight
                                            onBrushEnd={(area) => this.setState({lastDrawLocationVolumeChart: area})}
                                            onDrag={(area) => {
                                                this.setState({
                                                    lastDrawLocationVolumeChart: {
                                                        bottom: lastDrawLocationVolumeChart.bottom
                                                            + (area.top - area.bottom),
                                                        left: lastDrawLocationVolumeChart.left
                                                            - (area.right - area.left),
                                                        right: lastDrawLocationVolumeChart.right
                                                            - (area.right - area.left),
                                                        top: lastDrawLocationVolumeChart.top
                                                            + (area.top - area.bottom)
                                                    }
                                                });
                                            }}/>
                                    </FlexibleWidthXYPlot>
                                    <DiscreteColorLegend orientation="horizontal"
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
                            <CardHeader title="Request Duration" className={classes.cardHeader}
                                classes={{
                                    title: classes.title
                                }}
                                action={zoomHintToolTip}
                            />
                            <CardContent className={classes.content}>
                                <div>
                                    <FlexibleWidthXYPlot xType="time" height={400} animation
                                        xDomain={
                                            lastDrawLocationDurationChart && [
                                                lastDrawLocationDurationChart.left,
                                                lastDrawLocationDurationChart.right
                                            ]
                                        }
                                        onMouseLeave={() => this.setState({durationTooltip: false})}>
                                        <HorizontalGridLines/>
                                        <VerticalGridLines/>
                                        <XAxis title="Time"/>
                                        <YAxis title="Duration (ms)"/>
                                        <LineMarkSeries data={reqDurationData} color="#3f51b5" size={3}
                                            onNearestX={(d) => this.setState({durationTooltip: d})}/>
                                        {
                                            durationTooltip
                                                ? (
                                                    <Crosshair values={[durationTooltip]}>
                                                        <div className="rv-hint__content">
                                                            {
                                                                `${moment(durationTooltip.x)
                                                                    .format(Constants.Pattern.DATE_TIME)} :
                                                                ${Math.round(durationTooltip.y)} ms`
                                                            }
                                                        </div>
                                                    </Crosshair>
                                                )
                                                : null
                                        }
                                        <Highlight
                                            onBrushEnd={(area) => this.setState({lastDrawLocationDurationChart: area})}
                                            onDrag={(area) => {
                                                this.setState({
                                                    lastDrawLocationDurationChart: {
                                                        bottom: lastDrawLocationDurationChart.bottom
                                                        + (area.top - area.bottom),
                                                        left: lastDrawLocationDurationChart.left
                                                        - (area.right - area.left),
                                                        right: lastDrawLocationDurationChart.right
                                                        - (area.right - area.left),
                                                        top: lastDrawLocationDurationChart.top
                                                        + (area.top - area.bottom)
                                                    }
                                                });
                                            }}/>
                                    </FlexibleWidthXYPlot>
                                    <DiscreteColorLegend orientation="horizontal"
                                        items={[
                                            {
                                                title: "Request",
                                                color: "#3f51b5"
                                            }
                                        ]}/>
                                </div>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item md={6} sm={12} xs={12}>
                        <Card className={classes.card}>
                            <CardHeader title="Request/Response Size" className={classes.cardHeader}
                                classes={{
                                    title: classes.title
                                }}
                                action={zoomHintToolTip}
                            />
                            <CardContent className={classes.content}>
                                <div>
                                    <FlexibleWidthXYPlot xType="time" height={400} animation
                                        xDomain={
                                            lastDrawLocationSizeChart && [
                                                lastDrawLocationSizeChart.left,
                                                lastDrawLocationSizeChart.right
                                            ]
                                        }
                                        onMouseLeave={() => this.setState({sizeTooltip: []})}>
                                        <HorizontalGridLines/>
                                        <VerticalGridLines/>
                                        <XAxis title="Time"/>
                                        <YAxis title="Size (Bytes)"/>
                                        {
                                            reqResSizeData.map((dataItem, index) => (
                                                <LineMarkSeries key={dataItem.name} color={reqResColors[index]} size={3}
                                                    data={dataItem.data}
                                                    onNearestX={(d, {index}) => this.setState({
                                                        sizeTooltip: reqResSizeData.map((d) => ({
                                                            ...d.data[index],
                                                            name: d.name
                                                        }))
                                                    })}/>
                                            ))
                                        }
                                        <Crosshair values={sizeTooltip}>
                                            {
                                                sizeTooltip.length > 0
                                                    ? (
                                                        <div className="rv-hint__content">
                                                            <div className={classes.toolTipHead}>
                                                                {
                                                                    `${moment(sizeTooltip[0].x)
                                                                        .format(Constants.Pattern.DATE_TIME)}`
                                                                }
                                                            </div>
                                                            {
                                                                sizeTooltip.map((tooltipItem) => (
                                                                    <div key={tooltipItem.name}>
                                                                        {
                                                                            `${tooltipItem.name} Size:
                                                                            ${Math.round(tooltipItem.y)} Bytes`
                                                                        }
                                                                    </div>
                                                                ))
                                                            }
                                                        </div>
                                                    )
                                                    : null
                                            }
                                        </Crosshair>
                                        <Highlight
                                            onBrushEnd={(area) => this.setState({lastDrawLocationSizeChart: area})}
                                            onDrag={(area) => {
                                                this.setState({
                                                    lastDrawLocationSizeChart: {
                                                        bottom: lastDrawLocationSizeChart.bottom
                                                        + (area.top - area.bottom),
                                                        left: lastDrawLocationSizeChart.left - (area.right - area.left),
                                                        right: lastDrawLocationSizeChart.right
                                                        - (area.right - area.left),
                                                        top: lastDrawLocationSizeChart.top + (area.top - area.bottom)
                                                    }
                                                });
                                            }}/>
                                    </FlexibleWidthXYPlot>
                                    <DiscreteColorLegend orientation="horizontal"
                                        items={
                                            reqResSizeData.map((d, index) => ({
                                                title: d.name,
                                                color: reqResColors[index]
                                            }))
                                        }/>
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
    colorGenerator: PropTypes.instanceOf(ColorGenerator).isRequired,
    data: PropTypes.arrayOf(PropTypes.shape({
        timestamp: PropTypes.number.isRequired,
        httpResponseGroup: PropTypes.string.isRequired,
        totalRequestSizeBytes: PropTypes.number.isRequired,
        totalResponseSizeBytes: PropTypes.number.isRequired,
        totalResponseTimeMilliSec: PropTypes.number.isRequired,
        requestCount: PropTypes.number.isRequired
    })).isRequired,
    direction: PropTypes.string.isRequired
};

export default withStyles(styles)(withColor(MetricsGraphs));
