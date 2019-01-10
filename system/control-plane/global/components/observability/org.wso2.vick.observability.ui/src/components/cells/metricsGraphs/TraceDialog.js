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

import Button from "@material-ui/core/Button/Button";
import Constants from "../../../utils/constants";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";
import React from "react";
import moment from "moment/moment";
import {withStyles} from "@material-ui/core/styles";
import * as PropTypes from "prop-types";

const styles = (theme) => ({
    subTitle: {
        marginLeft: theme.spacing.unit,
        marginTop: theme.spacing.unit * 1.5,
        color: "#666",
        fontSize: 14
    },
    light: {
        color: "#999",
        fontStyle: "Italic"
    }
});

class TraceDialog extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            open: false
        };
    }

    handleClickOpen = () => {
        this.setState({open: true});
    };

    handleClose = () => {
        this.setState({open: false});
    };

    render = () => {
        const {classes, selectedArea} = this.props;
        const {open} = this.state;

        return (
            <Dialog
                fullWidth={true}
                maxWidth="lg"
                open={open}
                onClose={this.handleClose}
                aria-labelledby="max-width-dialog-title"
            >
                <DialogTitle id="max-width-dialog-title">Traces <span className={classes.subTitle}>
                    <span className={classes.light}> From</span> {
                        selectedArea
                            ? moment(selectedArea.left.valueOf()).format(Constants.Pattern.GRAPH_DATE_TIME)
                            : null}
                    <span className={classes.light}> to</span> {
                        selectedArea
                            ? moment(selectedArea.right.valueOf()).format(Constants.Pattern.GRAPH_DATE_TIME)
                            : null}</span> </DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Trace List Table
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={this.handleClose} color="primary">
                        Close
                    </Button>
                </DialogActions>
            </Dialog>
        );
    }

}

TraceDialog.propTypes = {
    classes: PropTypes.object.isRequired,
    selectedArea: PropTypes.objectOf(moment)
};

export default withStyles(styles)(TraceDialog);
