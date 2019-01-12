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

import AccountCircle from "@material-ui/icons/AccountCircle";
import AppBar from "@material-ui/core/AppBar/AppBar";
import AuthUtils from "../../utils/api/authUtils";
import Avatar from "@material-ui/core/Avatar/Avatar";
import ColorGenerator from "../common/color/colorGenerator";
import Constants from "../../utils/constants";
import Divider from "@material-ui/core/Divider/Divider";
import FileCopy from "@material-ui/icons/FileCopyOutlined";
import FormatColorFillOutlined from "@material-ui/icons/FormatColorFillOutlined";
import HttpUtils from "../../utils/api/httpUtils";
import IconButton from "@material-ui/core/IconButton/IconButton";
import InputBase from "@material-ui/core/InputBase/InputBase";
import Menu from "@material-ui/core/Menu/Menu";
import MenuIcon from "@material-ui/icons/Menu";
import MenuItem from "@material-ui/core/MenuItem/MenuItem";
import Paper from "@material-ui/core/Paper/Paper";
import Popover from "@material-ui/core/Popover/Popover";
import QueryUtils from "../../utils/common/queryUtils";
import React from "react";
import ShareIcon from "../../icons/ShareIcon";
import StateHolder from "../common/state/stateHolder";
import Toolbar from "@material-ui/core/Toolbar/Toolbar";
import Tooltip from "@material-ui/core/Tooltip/Tooltip";
import Typography from "@material-ui/core/Typography/Typography";
import classNames from "classnames";
import {deepPurple} from "@material-ui/core/colors";
import withColor from "../common/color";
import withGlobalState from "../common/state";
import {withRouter} from "react-router-dom";
import {withStyles} from "@material-ui/core";
import * as PropTypes from "prop-types";

const styles = (theme) => ({
    appBar: {
        zIndex: theme.zIndex.drawer + 1,
        transition: theme.transitions.create(["width", "margin"], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen
        })
    },
    appBarShift: {
        marginLeft: Constants.Dashboard.SIDE_NAV_BAR_WIDTH,
        width: `calc(100% - ${Constants.Dashboard.SIDE_NAV_BAR_WIDTH}px)`,
        transition: theme.transitions.create(["width", "margin"], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen
        })
    },
    menuButton: {
        marginLeft: 12,
        marginRight: 20
    },
    hide: {
        display: "none"
    },
    grow: {
        flexGrow: 1
    },
    popoverContent: {
        margin: theme.spacing.unit * 3
    },
    copyContainer: {
        marginTop: theme.spacing.unit * 2,
        boxShadow: "none",
        border: "1px solid #eee"
    },
    copyInput: {
        margin: theme.spacing.unit,
        flex: 1
    },
    iconButton: {
        marginLeft: theme.spacing.unit,
        marginRight: theme.spacing.unit
    },
    avatarContainer: {
        marginBottom: theme.spacing.unit * 2,
        pointerEvents: "none"
    },
    userAvatar: {
        marginRight: theme.spacing.unit * 1.5,
        color: "#fff",
        backgroundColor: deepPurple[500]
    }
});

class MainAppBar extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            isDashboardLinkCopiedTooltipOpen: false,
            dashBoardSharePopoverElement: null,
            accountPopoverElement: null
        };

        this.dashboardShareableLinkRef = React.createRef();
    }

    handleDashboardSharePopoverOpen = (event) => {
        this.setState({
            dashBoardSharePopoverElement: event.currentTarget
        });
    };

    handleDashboardSharePopoverClose = () => {
        this.setState({
            dashBoardSharePopoverElement: null
        });
    };

    handleAccountPopoverOpen = (event) => {
        this.setState({
            accountPopoverElement: event.currentTarget
        });
    };

    handleAccountPopoverClose = () => {
        this.setState({
            accountPopoverElement: null
        });
    };

    handleSideNavBarOpen = () => {
        const {onSideNavBarOpen} = this.props;
        if (onSideNavBarOpen) {
            onSideNavBarOpen();
        }
    };

    copyShareableLinkToClipboard = () => {
        if (this.dashboardShareableLinkRef.current) {
            // Copying the link onto the clipboard
            this.dashboardShareableLinkRef.current.select();
            document.execCommand("copy");

            this.setState({
                isDashboardLinkCopiedTooltipOpen: true
            });
        }
    };

    onDashboardLinkCopiedTooltipClose = () => {
        this.setState({
            isDashboardLinkCopiedTooltipOpen: false
        });
    };

    resetColorScheme = () => {
        const {colorGenerator} = this.props;
        colorGenerator.resetColors();
    };

    getShareableLink = () => {
        const {globalState, location} = this.props;

        const globalFilter = globalState.get(StateHolder.GLOBAL_FILTER);
        const queryParams = HttpUtils.parseQueryParams(location.search);
        queryParams.globalFilterStartTime = QueryUtils.parseTime(globalFilter.startTime).valueOf();
        queryParams.globalFilterEndTime = QueryUtils.parseTime(globalFilter.endTime).valueOf();

        return window.location.origin + location.pathname + HttpUtils.generateQueryParamString(queryParams);
    };

    render() {
        const {classes, isSideNavBarOpen, globalState} = this.props;
        const {dashBoardSharePopoverElement, isDashboardLinkCopiedTooltipOpen, accountPopoverElement} = this.state;

        const isAccountPopoverOpen = Boolean(accountPopoverElement);
        const isDashboardSharePopoverOpen = Boolean(dashBoardSharePopoverElement);
        const loggedInUser = globalState.get(StateHolder.USER);

        return (
            <AppBar position="fixed"
                className={classNames(classes.appBar, {
                    [classes.appBarShift]: isSideNavBarOpen
                })}>
                <Toolbar disableGutters={!isSideNavBarOpen}>
                    <IconButton color="inherit" aria-label="Open drawer" onClick={this.handleSideNavBarOpen}
                        className={classNames(classes.menuButton, {
                            [classes.hide]: isSideNavBarOpen
                        })}>
                        <MenuIcon/>
                    </IconButton>
                    <Typography variant="h6" color="inherit" className={classes.grow}>
                        Cellery Observability
                    </Typography>
                    {
                        loggedInUser
                            ? (
                                <div>
                                    <Tooltip title="Change color scheme" placement="bottom">
                                        <IconButton onClick={this.resetColorScheme} color="inherit">
                                            <FormatColorFillOutlined/>
                                        </IconButton>
                                    </Tooltip>
                                    <Tooltip title="Get shareable dashboard link" placement="bottom">
                                        <IconButton color="inherit" aria-haspopup="true" variant="contained"
                                            aria-owns={isDashboardSharePopoverOpen ? "share-dashboard" : undefined}
                                            onClick={this.handleDashboardSharePopoverOpen}>
                                            <ShareIcon/>
                                        </IconButton>
                                    </Tooltip>
                                    <Popover
                                        id="share-dashboard" open={isDashboardSharePopoverOpen}
                                        anchorEl={dashBoardSharePopoverElement}
                                        onClose={this.handleDashboardSharePopoverClose}
                                        anchorOrigin={{
                                            vertical: "bottom",
                                            horizontal: "right"
                                        }}
                                        transformOrigin={{
                                            vertical: "top",
                                            horizontal: "right"
                                        }}>
                                        <div className={classes.popoverContent}>
                                            <Typography color="textSecondary" variant="subtitle2" gutterBottom>
                                                Share Dashboard
                                            </Typography>
                                            <Divider/>
                                            <Paper className={classes.copyContainer} elevation={1}>
                                                <InputBase className={classes.copyInput}
                                                    placeholder="Sharable Link" value={this.getShareableLink()}
                                                    inputRef={this.dashboardShareableLinkRef}/>
                                                <Tooltip title="Copied!" disableFocusListener={false}
                                                    disableHoverListener={false} placement="top"
                                                    disableTouchListener={false} open={isDashboardLinkCopiedTooltipOpen}
                                                    onClose={this.onDashboardLinkCopiedTooltipClose}>
                                                    <IconButton color="primary" className={classes.iconButton}
                                                        aria-label="Copy" onClick={this.copyShareableLinkToClipboard}>
                                                        <FileCopy/>
                                                    </IconButton>
                                                </Tooltip>
                                            </Paper>
                                        </div>
                                    </Popover>
                                    <IconButton
                                        aria-owns={isAccountPopoverOpen ? "user-info-appbar" : undefined}
                                        color="inherit" aria-haspopup="true" onClick={this.handleAccountPopoverOpen}>
                                        <AccountCircle/>
                                    </IconButton>
                                    <Menu id="user-info-appbar" anchorEl={accountPopoverElement}
                                        anchorOrigin={{
                                            vertical: "top",
                                            horizontal: "right"
                                        }}
                                        transformOrigin={{
                                            vertical: "top",
                                            horizontal: "right"
                                        }}
                                        open={isAccountPopoverOpen}
                                        onClose={this.handleAccountPopoverClose}>
                                        <MenuItem onClick={this.handleAccountPopoverClose}
                                            className={classes.avatarContainer}>
                                            <Avatar className={classes.userAvatar}>
                                                {loggedInUser.username.substr(0, 1).toUpperCase()}
                                            </Avatar>
                                            {loggedInUser.username}
                                        </MenuItem>
                                        <MenuItem onClick={() => AuthUtils.signOut(globalState)}>
                                            Sign Out
                                        </MenuItem>
                                    </Menu>
                                </div>
                            )
                            : null
                    }
                </Toolbar>
            </AppBar>
        );
    }

}

MainAppBar.propTypes = {
    classes: PropTypes.object.isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator).isRequired,
    location: PropTypes.shape({
        pathname: PropTypes.string.isRequired
    }),
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    onSideNavBarOpen: PropTypes.func.isRequired,
    isSideNavBarOpen: PropTypes.bool.isRequired
};

export default withStyles(styles, {withTheme: true})(withRouter(withGlobalState(withColor(MainAppBar))));
