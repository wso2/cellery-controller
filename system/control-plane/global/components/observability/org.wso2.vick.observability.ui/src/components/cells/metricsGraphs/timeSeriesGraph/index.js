/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

import Button from "@material-ui/core/Button/Button";
import Card from "@material-ui/core/Card/Card";
import CardContent from "@material-ui/core/CardContent/CardContent";
import CardHeader from "@material-ui/core/CardHeader/CardHeader";
import Constants from "../../../../utils/constants";
import InfoIcon from "@material-ui/icons/InfoOutlined";
import React from "react";
import Timeline from "@material-ui/icons/Timeline";
import Tooltip from "@material-ui/core/Tooltip/Tooltip";
import TracesDialog from "./TracesDialog";
import moment from "moment";
import {withStyles} from "@material-ui/core";
import {
    Crosshair, DiscreteColorLegend, Highlight, HorizontalGridLines, LineMarkSeries, VerticalGridLines, XAxis, XYPlot,
    YAxis, makeWidthFlexible
} from "react-vis/es";
import * as PropTypes from "prop-types";

const styles = (theme) => ({
    cardHeader: {
        borderBottom: "1px solid #eee",
        paddingTop: theme.spacing.unit,
        paddingBottom: theme.spacing.unit
    },
    title: {
        fontSize: 16,
        fontWeight: 500,
        color: "#4d4d4d"
    },
    viewTracesButton: {
        fontSize: 12,
        verticalAlign: "middle",
        marginRight: theme.spacing.unit
    },
    viewTracesContent: {
        paddingLeft: theme.spacing.unit
    },
    infoIcon: {
        color: "#999",
        marginRight: 27,
        fontSize: 18,
        verticalAlign: "middle"
    },
    card: {
        boxShadow: "none",
        border: "1px solid #eee"
    },
    cardActions: {
        marginTop: theme.spacing.unit / 4
    },
    toolTipHead: {
        fontWeight: 600
    }
});

const FlexibleWidthXYPlot = makeWidthFlexible(XYPlot);

class TimeSeriesGraph extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            tooltipData: [],
            lastDrawArea: null
        };

        this.traceDialogRef = React.createRef();
    }

    render = () => {
        const {classes, traceSearchFilter, title, data, yAxis} = this.props;
        const {lastDrawArea, tooltipData} = this.state;
        return (
            <Card className={classes.card}>
                <CardHeader title={title} className={classes.cardHeader}
                    classes={{
                        title: classes.title,
                        action: classes.cardActions
                    }}
                    action={
                        (
                            <React.Fragment>
                                {
                                    lastDrawArea
                                        ? (
                                            <React.Fragment>
                                                <Tooltip title="View traces for the selected time range">
                                                    <Button className={classes.viewTracesButton} variant="outlined"
                                                        onClick={this.handleClickOpen}>
                                                        <Timeline color="action"/>
                                                        <span className={classes.viewTracesContent}>View Traces</span>
                                                    </Button>
                                                </Tooltip>
                                                <TracesDialog filter={traceSearchFilter} innerRef={this.traceDialogRef}
                                                    selectedArea={lastDrawArea}/>
                                            </React.Fragment>
                                        )
                                        : null
                                }
                                <Tooltip title={"Click and drag in the plot area to zoom in, click anywhere in the "
                                    + "graph to zoom out."}>
                                    <InfoIcon className={classes.infoIcon}/>
                                </Tooltip>
                            </React.Fragment>
                        )
                    }
                />
                <CardContent>
                    <div>
                        <FlexibleWidthXYPlot xType="time" height={400} animation
                            xDomain={
                                lastDrawArea
                                    ? [
                                        lastDrawArea.left,
                                        lastDrawArea.right
                                    ]
                                    : null
                            }
                            onMouseLeave={() => this.setState({tooltipData: []})}>
                            <HorizontalGridLines/>
                            <VerticalGridLines/>
                            <XAxis title="Time"/>
                            <YAxis title={`${yAxis.title} (${yAxis.unit})`}/>
                            {
                                data.map((datum) => (
                                    <LineMarkSeries key={datum.title} color={datum.color} size={3}
                                        data={datum.points.map((point) => ({
                                            x: point.timestamp,
                                            y: point.value
                                        }))}
                                        onNearestX={(d, {index}) => this.setState({
                                            tooltipData: data.map((datum) => ({
                                                ...datum.points[index],
                                                name: datum.title
                                            }))
                                        })}/>
                                ))
                            }
                            {
                                tooltipData.length > 0
                                    ? (
                                        <Crosshair values={tooltipData.map((tooltipDatum) => ({
                                            x: tooltipDatum.timestamp,
                                            y: tooltipDatum.value
                                        }))}>
                                            <div className="rv-hint__content">
                                                <div className={classes.toolTipHead}>
                                                    {
                                                        `${moment(tooltipData[0].timestamp)
                                                            .format(Constants.Pattern.DATE_TIME)}`
                                                    }
                                                </div>
                                                {
                                                    tooltipData.map((tooltipDatum) => (
                                                        <div key={tooltipDatum.name}>
                                                            {
                                                                `${tooltipDatum.name}:
                                                                    ${Math.round(tooltipDatum.value)} ${yAxis.unit}`
                                                            }
                                                        </div>
                                                    ))
                                                }
                                            </div>
                                        </Crosshair>
                                    )
                                    : null
                            }
                            <Highlight enableY={false}
                                onBrushEnd={(area) => this.setState({lastDrawArea: area})}
                                onDrag={(area) => {
                                    this.setState({
                                        lastDrawArea: {
                                            bottom: lastDrawArea.bottom + (area.top - area.bottom),
                                            left: lastDrawArea.left - (area.right - area.left),
                                            right: lastDrawArea.right - (area.right - area.left),
                                            top: lastDrawArea.top + (area.top - area.bottom)
                                        }
                                    });
                                }}/>
                        </FlexibleWidthXYPlot>
                        <DiscreteColorLegend orientation="horizontal"
                            items={data.map((datum) => ({
                                title: datum.title,
                                color: datum.color
                            }))}/>
                    </div>
                </CardContent>
            </Card>
        );
    };

    handleClickOpen = () => {
        if (this.traceDialogRef.current && this.traceDialogRef.current.handleClickOpen) {
            this.traceDialogRef.current.handleClickOpen();
        }
    };

}

TimeSeriesGraph.propTypes = {
    classes: PropTypes.object.isRequired,
    title: PropTypes.string.isRequired,
    data: PropTypes.arrayOf(PropTypes.shape({
        title: PropTypes.string.isRequired,
        points: PropTypes.arrayOf(PropTypes.shape({
            timestamp: PropTypes.number.isRequired,
            value: PropTypes.number.isRequired
        })).isRequired,
        color: PropTypes.string.isRequired
    })).isRequired,
    yAxis: PropTypes.shape({
        title: PropTypes.string.isRequired,
        unit: PropTypes.string.isRequired
    }).isRequired,
    traceSearchFilter: PropTypes.shape({
        cell: PropTypes.string,
        microservice: PropTypes.string,
        operation: PropTypes.string,
        tags: PropTypes.object,
        minDuration: PropTypes.number,
        minDurationMultiplier: PropTypes.number,
        maxDuration: PropTypes.number,
        maxDurationMultiplier: PropTypes.number
    }).isRequired
};

export default withStyles(styles)(TimeSeriesGraph);
