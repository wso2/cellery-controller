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
import Constants from "../../common/constants";
import Grid from "@material-ui/core/Grid";
import PropTypes from "prop-types";
import React from "react";
import moment from "moment";
import {withStyles} from "@material-ui/core/styles";
import {
    Crosshair,
    DiscreteColorLegend,
    HorizontalGridLines,
    LineSeries,
    VerticalGridLines,
    XAxis,
    XYPlot,
    YAxis,
    makeWidthFlexible
} from "react-vis";
import withColor, {ColorGenerator} from "../../common/color";

const styles = {
    card: {
        boxShadow: "none",
        border: "1px solid #eee"
    },
    cardHeader: {
        borderBottom: "1px solid #eee"
    }, title: {
        fontSize: 16,
        fontWeight: 500,
        color: "#4d4d4d"
    },
    contentGrid: {
        height: 186
    },
    toolTipHead: {
        fontWeight: 600
    }
};

const FlexibleWidthXYPlot = makeWidthFlexible(XYPlot);

class Metrics extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            tooltip: []
        };
    }

    render = () => {
        const {classes, colorGenerator, graphName, graphData} = this.props;
        const {tooltip} = this.state;
        const dateTimeFormat = Constants.Pattern.DATE_TIME;

        return (

            <Grid item md={6} sm={12} xs={12}>
                <Card className={classes.card}>
                    <CardHeader
                        classes={{
                            title: classes.title
                        }}
                        title={graphName}
                        className={classes.cardHeader}
                    />
                    <CardContent className={classes.content}>
                        <div>
                            <FlexibleWidthXYPlot xType="time" height={400}
                                onMouseLeave={() => this.setState({tooltip: []})}>
                                <HorizontalGridLines/>
                                <VerticalGridLines/>
                                <XAxis title="Time"/>
                                <YAxis title="Usage (%)"/>

                                {
                                    graphData.map((dataItem) => (
                                        <LineSeries
                                            key={dataItem.name}
                                            data={dataItem.data}
                                            onNearestX={(d, {index}) => this.setState({
                                                tooltip: graphData.map((d) => ({
                                                    ...d.data[index],
                                                    name: d.name
                                                }))
                                            })}
                                            color={colorGenerator.getColor(dataItem.name)}
                                        />))
                                }

                                <Crosshair values={tooltip}>
                                    {
                                        tooltip.length > 0
                                            ? (
                                                <div className="rv-hint__content">
                                                    <div className={classes.toolTipHead}>
                                                        {`${moment(tooltip[0].x).format(dateTimeFormat)}`}
                                                    </div>
                                                    {tooltip.map((tooltipItem, {index}) => <div key={tooltipItem.name}>
                                                        {`${tooltipItem.name}: ${tooltipItem.y}%`}</div>)
                                                    }
                                                </div>
                                            )
                                            : null
                                    }

                                </Crosshair>
                            </FlexibleWidthXYPlot>
                            <DiscreteColorLegend
                                orientation="horizontal"
                                items={graphData.map((d) => ({title: d.name, color: colorGenerator.getColor(d.name)}))}
                            />
                        </div>
                    </CardContent>
                </Card>
            </Grid>
        );
    }

}

Metrics.propTypes = {
    classes: PropTypes.object.isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator),
    graphName: PropTypes.string.isRequired,
    graphData: PropTypes.arrayOf(PropTypes.shape({
        name: PropTypes.string.isRequired,
        data: PropTypes.arrayOf(PropTypes.shape({
            x: PropTypes.number.isRequired,
            y: PropTypes.number.isRequired
        })).isRequired
    })).isRequired
};

export default withStyles(styles)(withColor(Metrics));
