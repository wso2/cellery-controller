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

import MUIDataTable from "mui-datatables";
import PropTypes from "prop-types";
import React from "react";
import {withStyles} from "@material-ui/core/styles";

const styles = (theme) => ({
    table: {
        minWidth: 1020
    },
    tableWrapper: {
        overflowX: "auto",
        marginLeft: theme.spacing.unit * 3,
        marginRight: theme.spacing.unit * 3
    }
});

const DataTable = (props) => {
    const {classes, options, columns, data} = props;
    const defaultOptions = {
        download: false,
        selectableRows: false,
        print: false,
        filterType: "dropdown",
        responsive: "scroll"
    };

    return (
        <div className={classes.tableWrapper}>
            <MUIDataTable data={data} columns={columns} options={{...defaultOptions, ...options}}/>
        </div>
    );
};

DataTable.propTypes = {
    data: PropTypes.arrayOf(PropTypes.arrayOf(PropTypes.any)).isRequired,
    options: PropTypes.object.isRequired,
    columns: PropTypes.arrayOf(PropTypes.any).isRequired,
    classes: PropTypes.object.isRequired
};

export default withStyles(styles)(DataTable);


