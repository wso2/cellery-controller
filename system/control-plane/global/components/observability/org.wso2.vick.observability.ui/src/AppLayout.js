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

/* eslint max-len: ["off"] */

import AccountCircle from "@material-ui/icons/AccountCircle";
import AppBar from "@material-ui/core/AppBar";
import AuthUtils from "./pages/common/utils/authUtils";
import CheckCircle from "@material-ui/icons/CheckCircle";
import ChevronLeftIcon from "@material-ui/icons/ChevronLeft";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import CircularProgress from "@material-ui/core/CircularProgress";
import CloseIcon from "@material-ui/icons/Close";
import Collapse from "@material-ui/core/Collapse";
import CssBaseline from "@material-ui/core/CssBaseline";
import Divider from "@material-ui/core/Divider";
import Drawer from "@material-ui/core/Drawer";
import Error from "@material-ui/icons/Error";
import ExpandLess from "@material-ui/icons/ExpandLess";
import ExpandMore from "@material-ui/icons/ExpandMore";
import IconButton from "@material-ui/core/IconButton";
import Info from "@material-ui/icons/Info";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import Menu from "@material-ui/core/Menu";
import MenuIcon from "@material-ui/icons/Menu";
import MenuItem from "@material-ui/core/MenuItem";
import NotificationUtils from "./pages/common/utils/notificationUtils";
import React from "react";
import Settings from "@material-ui/icons/SettingsOutlined";
import Snackbar from "@material-ui/core/Snackbar/Snackbar";
import SvgIcon from "@material-ui/core/SvgIcon";
import Timeline from "@material-ui/icons/Timeline";
import Toolbar from "@material-ui/core/Toolbar";
import Tooltip from "@material-ui/core/Tooltip";
import Typography from "@material-ui/core/Typography";
import Warning from "@material-ui/icons/Warning";
import classNames from "classnames";
import {withRouter} from "react-router-dom";
import {withStyles} from "@material-ui/core/styles";
import withGlobalState, {StateHolder} from "./pages/common/state";
import * as PropTypes from "prop-types";

const drawerWidth = 240;

const styles = (theme) => ({
    root: {
        display: "flex",
        flexGrow: 1,
        minHeight: "100%"
    },
    grow: {
        flexGrow: 1
    },
    appBar: {
        zIndex: theme.zIndex.drawer + 1,
        transition: theme.transitions.create(["width", "margin"], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen
        })
    },
    appBarShift: {
        marginLeft: drawerWidth,
        width: `calc(100% - ${drawerWidth}px)`,
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
    drawer: {
        width: drawerWidth,
        flexShrink: 0,
        whiteSpace: "nowrap"
    },
    drawerOpen: {
        width: drawerWidth,
        transition: theme.transitions.create("width", {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen
        })
    },
    drawerClose: {
        transition: theme.transitions.create("width", {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen
        }),
        overflowX: "hidden",
        width: theme.spacing.unit * 7 + 1,
        [theme.breakpoints.up("sm")]: {
            width: theme.spacing.unit * 9 + 1
        }
    },
    toolbar: {
        display: "flex",
        alignItems: "center",
        justifyContent: "flex-end",
        padding: "0 8px",
        ...theme.mixins.toolbar
    },
    content: {
        position: "relative",
        flexGrow: 1,
        padding: theme.spacing.unit * 3,
        minHeight: "100%"
    },
    progressOverlay: {
        position: "absolute",
        zIndex: 9999,
        top: 0,
        left: 0,
        width: "100%",
        height: "100%",
        backgroundColor: "rgb(0, 0, 0, 0.5)"
    },
    progress: {
        textAlign: "center",
        margin: "auto"
    },
    progressIndicator: {
        margin: theme.spacing.unit * 2
    },
    progressContent: {
        fontSize: "large",
        fontWeight: 500,
        width: "100%",
        color: "#ffffff"
    },
    nested: {
        paddingLeft: theme.spacing.unit * 3
    },
    active: {
        color: theme.palette.primary.main,
        fontWeight: 500
    },
    list: {
        paddingTop: 0
    },
    snackbarIcon: {
        fontSize: "1.5em"
    },
    snackbarMessageContainer: {
        display: "flex",
        alignItems: "center"
    },
    snackbarMessage: {
        paddingLeft: theme.spacing.unit
    }
});

const OverviewIcon = (props) => (
    <SvgIcon viewBox="0 0 14 14" {...props}>
        <path d="M12.5,0.8H1.5c-0.7,0-1.2,0.6-1.2,1.2v7.4c0,0.7,0.6,1.2,1.2,1.2h4.3v1.2H4.5v1.2h4.9v-1.2H8.2v-1.2h4.3
c0.7,0,1.2-0.6,1.2-1.2V2.1C13.8,1.4,13.2,0.8,12.5,0.8z M12.5,9.5H1.5V2.1h11.1V9.5z M10,9L7.9,7.8V5.5L10,4.3L12,5.5v2.4L10,9z
 M9.1,7.2L10,7.7l0.9-0.5v-1L10,5.7L9.1,6.2V7.2z M8,6.8L5.5,5.4l0.4-0.6l2.5,1.4L8,6.8z M4,7.2L1.9,6V3.7L4,2.5L6,3.7V6L4,7.2z
 M3.1,5.4L4,5.9l0.9-0.5v-1L4,3.8L3.1,4.3V5.4z"/>
    </SvgIcon>
);

const CellsIcon = (props) => (
    <SvgIcon viewBox="0 0 14 14" {...props}>
        <path d="M10.7,13.7l-3.1-1.8V8.4l3.1-1.8l3.1,1.8v3.6L10.7,13.7z M7,7.5L3.8,5.7V2L7,0.2L10.2,2v3.6L7,7.5z M5.1,4.9 l1.9,1l1.9
-1V 2.8L7,1.8l-1.9,1V4.9z M3.4,13.8L0.2,12V8.4l3.2-1.8l3.2,1.8V12L3.4,13.8z M8.8,11.2l1.9,1l1.9-1V9.1l-1.9 -1l-1.9,
1V11.2z M1.5,11.2l1.9,1l1.9-1V9.1l-1.9-1l-1.9,1V11.2z"/>
    </SvgIcon>
);
const PodIcon = (props) => (
    <SvgIcon viewBox="0 0 14 14" {...props}>
        <path
            d="M13.4,3.3C13.4,3.3,13.3,3.2,13.4,3.3c-0.1-0.3-0.3-0.4-0.4-0.4L7.2,0.1c-0.2-0.1-0.3-0.1-0.5,0L0.9,2.9 C0.7,3,0.6,3.2
,0.6,3.5v7.1c0,0.2,0.1,0.4,0.3,0.5l5.7,2.8l0,0l0,0l0.2,0c0,0,0.1,0,0,0l6-2.9c0.2-0.1,0.3-0.3,0.3-0.5L13.4,3.3 L13.4
,3.3z M6.4,6.7v5.8l-4.6-2.2V4.4L6.4,6.7z M11.5,3.5L7,5.6L2.5,3.5L7,1.3L11.5,3.5z M12.2,4.4v5.8l-4.7,2.3V6.7L12.2,4.
4z"/>
    </SvgIcon>
);

const NodeIcon = (props) => (
    <SvgIcon viewBox="0 0 14 14" {...props}>
        <path d="M12.4099731,2.5299683c0,0.2200317-0.1799927,0.4000244-0.3999634,0.4000244H6.2299805
c-0.2199707,0-0.3999634-0.1799927-0.3999634-0.4000244c0-0.2199707,0.1799927-0.3999634,0.3999634-0.3999634h5.7800293
C12.2299805,2.1300049,12.4099731,2.3099976,12.4099731,2.5299683z M6.2299805,7.0100098h5.7800293
c0.2199707,0,0.3999634-0.1800537,0.3999634-0.4000244s-0.1799927-0.4000244-0.3999634-0.4000244H6.2299805
c-0.2199707,0-0.3999634,0.1800537-0.3999634,0.4000244S6.0100098,7.0100098,6.2299805,7.0100098z M14,12.3800049
c0,0.8599854-0.7000122,1.5599976-1.5499878,1.5599976c-0.6799927,0-1.2700195-0.4400024-1.4700317-1.0599976H8.460022
C8.2600098,13.5,7.6799927,13.9400024,7,13.9400024S5.7399902,13.5,5.539978,12.8800049H3.0200195
C2.8200073,13.5,2.2299805,13.9400024,1.5499878,13.9400024C0.7000122,13.9400024,0,13.2399902,0,12.3800049
c0-0.8500366,0.7000122-1.5500488,1.5499878-1.5500488c0.6799927,0,1.2600098,0.4400024,1.4700317,1.0500488H5.539978
c0.1600342-0.4800415,0.5400391-0.8500366,1.0300293-0.9800415V9.8299561c0-0.2799683,0.2299805-0.5,0.5-0.5
c0.2799683,0,0.5,0.2200317,0.5,0.5v1.1100464c0.4199829,0.1599731,0.75,0.5100098,0.8900146,0.9400024h2.5199585
c0.210022-0.6100464,0.7900391-1.0500488,1.4700317-1.0500488C13.2999878,10.8299561,14,11.5299683,14,12.3800049z
 M2.1099854,12.3800049c0-0.3000488-0.25-0.5500488-0.5599976-0.5500488C1.25,11.8299561,1,12.0799561,1,12.3800049
c0,0.3099976,0.25,0.5599976,0.5499878,0.5599976C1.8599854,12.9400024,2.1099854,12.6900024,2.1099854,12.3800049z M13,12.3800049
c0-0.3000488-0.25-0.5500488-0.5499878-0.5500488c-0.3099976,0-0.5599976,0.25-0.5599976,0.5500488
c0,0.3099976,0.25,0.5599976,0.5599976,0.5599976C12.75,12.9400024,13,12.6900024,13,12.3800049z M3.0800171,1.6300049
c-0.5,0-0.9100342,0.3999634-0.9100342,0.8999634s0.4100342,0.9100342,0.9100342,0.9100342s0.8999634-0.4100342,0.8999634-0.9100342
S3.5800171,1.6300049,3.0800171,1.6300049z M3.0800171,7.5199585c0.5,0,0.8999634-0.4099731,0.8999634-0.9099731
S3.5800171,5.7000122,3.0800171,5.7000122S2.1699829,6.1099854,2.1699829,6.6099854S2.5800171,7.5199585,3.0800171,7.5199585z
 M14,0.7999878v3.2799683v0.9900513v3.2699585c0,0.4500122-0.3599854,0.8000488-0.7999878,0.8000488H0.7999878
C0.3599854,9.1400146,0,8.789978,0,8.3399658V5.0700073V4.0799561V0.7999878C0,0.3599854,0.3599854,0,0.7999878,0h12.4000244
C13.6400146,0,14,0.3599854,14,0.7999878z M13,5.0799561H1v3.0600586h12V5.0799561z M1,4.0700073h12V1H1V4.0700073z"/>
    </SvgIcon>
);

const MetricsIcon = (props) => (
    <SvgIcon viewBox="0 0 14 14" {...props}>
        <path d="M1.6,12.5h2.2c0.3,0,0.6-0.3,0.6-0.6V7.1c0-0.3-0.3-0.6-0.6-0.6H1.6C1.3,6.5,1,6.8,1,7.1v4.8C1,12.3,1.3,12.5,1.6,12.5z
 M2,7.5h1.4v4H2V7.5z M6,12.6h2.2c0.3,0,0.6-0.3,0.6-0.6V1c0-0.3-0.3-0.6-0.6-0.6H6C5.7,0.3,5.4,0.6,5.4,1v11
C5.4,12.3,5.7,12.6,6,12.6z M6.4,1.3h1.5v10.3H6.4V1.3z M10.4,12.6h2.2c0.3,0,0.6-0.3,0.6-0.6V4.5c0-0.3-0.3-0.6-0.6-0.6h-2.2
c-0.3,0-0.6,0.3-0.6,0.6V12C9.8,12.3,10.1,12.6,10.4,12.6z M10.8,4.9h1.4v6.7h-1.4V4.9z M13.3,13.2v0.1c0,0.2-0.2,0.4-0.4,0.4H1.5
c-0.2,0-0.4-0.2-0.4-0.4v-0.1c0-0.2,0.2-0.4,0.4-0.4h11.4C13.1,12.8,13.3,13,13.3,13.2z"/>
    </SvgIcon>
);

class AppLayout extends React.Component {

    constructor(props) {
        super(props);

        const pages = ["/", "/cells", "/tracing", "/system-metrics"];
        let selectedIndex = 0;
        for (let i = 0; i < pages.length; i++) {
            if (props.location.pathname.startsWith(pages[i])) {
                selectedIndex = i;
            }
        }

        props.globalState.addListener(StateHolder.LOADING_STATE, this.handleLoadingStateChange);
        props.globalState.addListener(StateHolder.NOTIFICATION_STATE, this.handleNotificationStateChange);

        const loadingState = props.globalState.get(StateHolder.LOADING_STATE);
        const notificationState = props.globalState.get(StateHolder.NOTIFICATION_STATE);
        this.state = {
            open: false,
            userInfo: null,
            subMenuOpen: false,
            loadingState: {
                ...loadingState
            },
            notificationState: {
                ...notificationState
            },
            selectedIndex: selectedIndex
        };
    }

    handleUserInfoMenuOpen = (event) => {
        this.setState({userInfo: event.currentTarget});
    };

    handleUserInfoMenuClose = () => {
        this.setState({userInfo: null});
    };

    handleDrawerOpen = () => {
        this.setState({open: true});
    };

    handleDrawerClose = () => {
        this.setState({open: false});
    };

    handleSystemMetricsNavSectionClick = () => {
        this.setState((prevState) => ({subMenuOpen: !prevState.subMenuOpen}));
    };

    handleNavItemClick = (nav, event) => {
        const {history} = this.props;

        this.setState({
            selectedIndex: Number(event.currentTarget.attributes.index.value)
        });
        history.push(nav, {
            hideBackButton: true
        });
    };

    handleLoadingStateChange = (loadingStateKey, oldState, newState) => {
        this.setState({
            loadingState: {
                isLoading: newState.loadingOverlayCount > 0,
                message: newState.message
            }
        });
    };

    handleNotificationStateChange = (notificationStateKey, oldState, newState) => {
        this.setState({
            notificationState: {
                isOpen: newState.isOpen,
                message: newState.message,
                notificationLevel: newState.notificationLevel
            }
        });
    };

    handleNotificationClose = () => {
        NotificationUtils.closeNotification(this.props.globalState);
    };

    generateSnackbarMessage = () => {
        const {classes} = this.props;
        const {notificationState} = this.state;

        let Icon;
        switch (notificationState.notificationLevel) {
            case NotificationUtils.Levels.SUCCESS:
                Icon = CheckCircle;
                break;
            case NotificationUtils.Levels.WARNING:
                Icon = Warning;
                break;
            case NotificationUtils.Levels.ERROR:
                Icon = Error;
                break;
            default:
                Icon = Info;
        }

        return (
            <span className={classes.snackbarMessageContainer}>
                <Icon className={classes.snackbarIcon}/>
                <span className={classes.snackbarMessage}>{notificationState.message}</span>
            </span>
        );
    };

    render = () => {
        const {classes, children, theme, globalState} = this.props;
        const {open, userInfo, loadingState, selectedIndex} = this.state;
        const userInfoOpen = Boolean(userInfo);
        return (
            <div className={classes.root}>
                <CssBaseline/>
                <AppBar position="fixed"
                    className={classNames(classes.appBar, {
                        [classes.appBarShift]: open
                    })}>
                    <Toolbar disableGutters={!open}>
                        <IconButton color="inherit" aria-label="Open drawer"
                            onClick={this.handleDrawerOpen}
                            className={classNames(classes.menuButton, {
                                [classes.hide]: open
                            })}>
                            <MenuIcon/>
                        </IconButton>
                        <Typography variant="h6" color="inherit" className={classes.grow}>
                            WSO2 VICK Observability
                        </Typography>
                        {
                            globalState.get(StateHolder.USER)
                                ? (
                                    <div>
                                        <IconButton
                                            aria-owns={userInfoOpen ? "user-info-appbar" : undefined}
                                            aria-haspopup="true"
                                            onClick={this.handleUserInfoMenuOpen}
                                            color="inherit">
                                            <AccountCircle/>
                                        </IconButton>
                                        <Menu id="user-info-appbar" anchorEl={userInfo}
                                            anchorOrigin={{
                                                vertical: "top",
                                                horizontal: "right"
                                            }}
                                            transformOrigin={{
                                                vertical: "top",
                                                horizontal: "right"
                                            }}
                                            open={userInfoOpen}
                                            onClose={this.handleUserInfoMenuClose}>
                                            {/* TODO: Implement user login */}
                                            <MenuItem onClick={this.handleUserInfoMenuClose}>
                                                Profile - {globalState.get(StateHolder.USER)}
                                            </MenuItem>
                                            <MenuItem onClick={this.handleUserInfoMenuClose}>
                                                My account
                                            </MenuItem>
                                            <MenuItem onClick={() => {
                                                AuthUtils.signOut(globalState);
                                            }}>
                                                Logout
                                            </MenuItem>
                                        </Menu>
                                    </div>
                                )
                                : null
                        }
                    </Toolbar>
                </AppBar>
                <Drawer variant="permanent"
                    className={classNames(classes.drawer, {
                        [classes.drawerOpen]: open,
                        [classes.drawerClose]: !open
                    })}
                    classes={{
                        paper: classNames({
                            [classes.drawerOpen]: open,
                            [classes.drawerClose]: !open
                        })
                    }}
                    open={open}>
                    <div className={classes.toolbar}>
                        <IconButton onClick={this.handleDrawerClose}>
                            {theme.direction === "rtl" ? <ChevronRightIcon/> : <ChevronLeftIcon/>}
                        </IconButton>
                    </div>
                    <Divider/>
                    <List className={classes.list}>
                        <Tooltip title="Overview" disableFocusListener={open} disableHoverListener={open}
                            placement="right"
                            disableTouchListener={open}>
                            <ListItem index={0} button key="Overview"
                                className={classNames({[classes.active]: selectedIndex === 0})}
                                onClick={(event) => {
                                    this.handleNavItemClick("/", event);
                                }}>
                                <ListItemIcon>
                                    <OverviewIcon className={classNames({[classes.active]: selectedIndex === 0})}/>
                                </ListItemIcon>
                                <ListItemText primary="Overview"
                                    classes={{primary: classNames({[classes.active]: selectedIndex === 0})}}/>
                            </ListItem>
                        </Tooltip>
                        <Tooltip title="Cells" disableFocusListener={open} disableHoverListener={open} placement="right"
                            disableTouchListener={open}>
                            <ListItem index={1} button key="Cells"
                                className={classNames({[classes.active]: selectedIndex === 1})}
                                onClick={(event) => {
                                    this.handleNavItemClick("/cells", event);
                                }}>
                                <ListItemIcon>
                                    <CellsIcon className={classNames({[classes.active]: selectedIndex === 1})}/>
                                </ListItemIcon>
                                <ListItemText primary="Cells"
                                    classes={{primary: classNames({[classes.active]: selectedIndex === 1})}}/>
                            </ListItem></Tooltip>
                        <Tooltip title="Distributed Tracing" disableFocusListener={open} disableHoverListener={open}
                            placement="right"
                            disableTouchListener={open}>
                            <ListItem index={2} button key="Distributed Tracing"
                                className={classNames({[classes.active]: selectedIndex === 2})}
                                onClick={(event) => {
                                    this.handleNavItemClick("/tracing", event);
                                }}>
                                <ListItemIcon>
                                    <Timeline className={classNames({[classes.active]: selectedIndex === 2})}/>
                                </ListItemIcon>
                                <ListItemText primary="Distributed Tracing"
                                    classes={{primary: classNames({[classes.active]: selectedIndex === 2})}}/>
                            </ListItem>
                        </Tooltip>
                        <Tooltip title="System Metrics" disableFocusListener={open} disableHoverListener={open}
                            placement="right"
                            disableTouchListener={open}>
                            <ListItem button onClick={this.handleSystemMetricsNavSectionClick}>
                                <ListItemIcon>
                                    <MetricsIcon/>
                                </ListItemIcon>
                                <ListItemText inset primary="System Metrics"/>
                                {this.state.subMenuOpen ? <ExpandLess/> : <ExpandMore/>}
                            </ListItem>
                        </Tooltip>
                        <Collapse in={this.state.subMenuOpen} timeout="auto" unmountOnExit>
                            <List component="div" disablePadding>
                                <Tooltip title="Global Control Plane" disableFocusListener={open}
                                    disableHoverListener={open} placement="right"
                                    disableTouchListener={open}>
                                    <ListItem index={3} button key="ControlPlane"
                                        className={classNames({[classes.active]: selectedIndex === 3},
                                            classes.nested)}
                                        onClick={(event) => {
                                            this.handleNavItemClick("/system-metrics/control-plane", event);
                                        }}>
                                        <ListItemIcon>
                                            <Settings className={classNames({[classes.active]: selectedIndex === 3})}/>
                                        </ListItemIcon>
                                        <ListItemText inset primary="Global Control Plane"
                                            classes={{
                                                primary: classNames({[classes.active]: selectedIndex === 3})
                                            }}/>
                                    </ListItem>
                                </Tooltip>
                                <Tooltip title="Pod Usage" disableFocusListener={open} disableHoverListener={open}
                                    placement="right"
                                    disableTouchListener={open}>
                                    <ListItem index={4} button key="PodUsage"
                                        className={classNames({[classes.active]: selectedIndex === 4},
                                            classes.nested)}
                                        onClick={(event) => {
                                            this.handleNavItemClick("/system-metrics/pod-usage", event);
                                        }}>
                                        <ListItemIcon>
                                            <PodIcon className={classNames({[classes.active]: selectedIndex === 4})}/>
                                        </ListItemIcon>
                                        <ListItemText inset primary="Pod Usage"
                                            classes={{
                                                primary: classNames({[classes.active]: selectedIndex === 4})
                                            }}/>
                                    </ListItem>
                                </Tooltip>
                                <Tooltip title="Node Usage" disableFocusListener={open} disableHoverListener={open}
                                    placement="right"
                                    disableTouchListener={open}>
                                    <ListItem index={5} button key="NodeUsage"
                                        className={classNames({[classes.active]: selectedIndex === 5},
                                            classes.nested)}
                                        onClick={(event) => {
                                            this.handleNavItemClick("/system-metrics/node-usage", event);
                                        }}>
                                        <ListItemIcon>
                                            <NodeIcon className={classNames({[classes.active]: selectedIndex === 5})}/>
                                        </ListItemIcon>
                                        <ListItemText inset primary="Node Usage"
                                            classes={{
                                                primary: classNames({[classes.active]: selectedIndex === 5})
                                            }}/>
                                    </ListItem>
                                </Tooltip>
                            </List>
                        </Collapse>
                    </List>
                </Drawer>
                <main className={classes.content}>
                    <div className={classes.progressOverlay} style={{
                        display: loadingState.isLoading ? "grid" : "none"
                    }}>
                        <div className={classes.progress}>
                            <CircularProgress className={classes.progressIndicator} thickness={4.5} size={45}/>
                            <div className={classes.progressContent}>
                                {loadingState.message ? loadingState.message : "Loading"}...
                            </div>
                        </div>
                    </div>
                    {children}
                </main>
                <Snackbar
                    anchorOrigin={{
                        vertical: "bottom",
                        horizontal: "left"
                    }}
                    open={this.state.notificationState.isOpen}
                    autoHideDuration={5000}
                    onClose={this.handleNotificationClose}
                    ContentProps={{"aria-describedby": "message-id"}}
                    message={this.generateSnackbarMessage()}
                    action={[
                        <IconButton key="close" aria-label="Close" color="inherit"
                            onClick={this.handleNotificationClose}>
                            <CloseIcon/>
                        </IconButton>
                    ]}
                />
            </div>
        );
    };

}

AppLayout.propTypes = {
    classes: PropTypes.object.isRequired,
    children: PropTypes.any.isRequired,
    theme: PropTypes.object.isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    history: PropTypes.any.isRequired,
    location: PropTypes.shape({
        pathname: PropTypes.string.isRequired
    })
};

export default withStyles(styles, {withTheme: true})(withRouter(withGlobalState(AppLayout)));
