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


import AccountCircle from "@material-ui/icons/AccountCircle";
import AppBar from "@material-ui/core/AppBar";
import AuthUtils from "../utils/api/authUtils";
import CellsIcon from "../icons/CellsIcon";
import CheckCircle from "@material-ui/icons/CheckCircle";
import ChevronLeftIcon from "@material-ui/icons/ChevronLeft";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import CircularProgress from "@material-ui/core/CircularProgress";
import CloseIcon from "@material-ui/icons/Close";
import Collapse from "@material-ui/core/Collapse";
import ColorGenerator from "./common/color/colorGenerator";
import CssBaseline from "@material-ui/core/CssBaseline";
import Divider from "@material-ui/core/Divider";
import Drawer from "@material-ui/core/Drawer";
import Error from "@material-ui/icons/Error";
import ExpandLess from "@material-ui/icons/ExpandLess";
import ExpandMore from "@material-ui/icons/ExpandMore";
import FileCopy from "@material-ui/icons/FileCopyOutlined";
import FormatColorFillOutlined from "@material-ui/icons/FormatColorFillOutlined";
import IconButton from "@material-ui/core/IconButton";
import Info from "@material-ui/icons/Info";
import InputBase from "@material-ui/core/InputBase";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import Menu from "@material-ui/core/Menu";
import MenuIcon from "@material-ui/icons/Menu";
import MenuItem from "@material-ui/core/MenuItem";
import MetricsIcon from "../icons/MetricsIcon";
import NodeIcon from "../icons/NodeIcon";
import NotificationUtils from "../utils/common/notificationUtils";
import OverviewIcon from "../icons/OverviewIcon";
import Paper from "@material-ui/core/Paper/Paper";
import PodIcon from "../icons/PodIcon";
import Popover from "@material-ui/core/Popover";
import React from "react";
import Settings from "@material-ui/icons/SettingsOutlined";
import ShareIcon from "../icons/ShareIcon";
import Snackbar from "@material-ui/core/Snackbar/Snackbar";
import Timeline from "@material-ui/icons/Timeline";
import Toolbar from "@material-ui/core/Toolbar";
import Tooltip from "@material-ui/core/Tooltip";
import Typography from "@material-ui/core/Typography";
import Warning from "@material-ui/icons/Warning";
import classNames from "classnames";
import withColor from "./common/color";
import {withRouter} from "react-router-dom";
import {withStyles} from "@material-ui/core/styles";
import withGlobalState, {StateHolder} from "./common/state";
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
    progressOverlayContainer: {
        position: "absolute",
        zIndex: 9999,
        top: 0,
        left: 0,
        width: "100%",
        height: "100%"
    },
    progressOverlay: {
        position: "relative",
        display: "grid",
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
    }
});

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
            selectedIndex: selectedIndex,
            popoverEl: null,
            showCopyText: false
        };
        this.popoverInputRef = React.createRef();
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

    handlePopoverClick = (event) => {
        this.setState({
            popoverEl: event.currentTarget
        });
    };

    handlePopoverClose = () => {
        this.setState({
            popoverEl: null
        });
    };

    copyLink = () => {
        if (this.popoverInputRef.current) {
            this.popoverInputRef.current.select();
            document.execCommand("copy");
            this.setState({
                showCopyText: true
            });
        }
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

    resetColorScheme = () => {
        this.props.colorGenerator.resetColors();
    };

    render = () => {
        const {classes, children, theme, globalState} = this.props;
        const {open, userInfo, loadingState, selectedIndex, popoverEl, showCopyText} = this.state;
        const userInfoOpen = Boolean(userInfo);
        const popoverOpen = Boolean(popoverEl);

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
                            Cellery Observability
                        </Typography>
                        {
                            globalState.get(StateHolder.USER)
                                ? (
                                    <div>
                                        <Tooltip title="Change cell colors" placement="bottom">
                                            <IconButton onClick={this.resetColorScheme} color="inherit">
                                                <FormatColorFillOutlined/>
                                            </IconButton>
                                        </Tooltip>
                                        <Tooltip title="Get shareable dashboard link" placement="bottom">
                                            <IconButton
                                                color="inherit" aria-owns={popoverOpen ? "share-dashboard" : undefined}
                                                aria-haspopup="true" variant="contained"
                                                onClick={this.handlePopoverClick}>
                                                <ShareIcon/>
                                            </IconButton>
                                        </Tooltip>
                                        <Popover
                                            id="share-dashboard" open={popoverOpen}
                                            anchorEl={popoverEl} onClose={this.handlePopoverClose}
                                            anchorOrigin={{
                                                vertical: "bottom",
                                                horizontal: "right"
                                            }}
                                            transformOrigin={{
                                                vertical: "top",
                                                horizontal: "right"
                                            }}>
                                            <div className={classes.popoverContent}>
                                                <Typography color="textSecondary" variant="subtitle2" gutterBottom
                                                    className={classes.typography}>Share Dashboard</Typography>
                                                <Divider/>
                                                <Paper className={classes.copyContainer} elevation={1}>
                                                    {/* TODO: replace InputBase value to dashboard URL*/}
                                                    <InputBase className={classes.copyInput}
                                                        placeholder="Sharable Link" value="http://cellery-dashboard/"
                                                        inputRef={this.popoverInputRef}/>
                                                    <Tooltip title="Copied!" disableFocusListener={false}
                                                        disableHoverListener={false} placement="top"
                                                        disableTouchListener={false} open={showCopyText}
                                                        onClose={() => this.setState({showCopyText: false})}>
                                                        <IconButton color="primary" className={classes.iconButton}
                                                            aria-label="Copy" onClick={this.copyLink}>
                                                            <FileCopy/>
                                                        </IconButton>
                                                    </Tooltip>
                                                </Paper>
                                            </div>
                                        </Popover>
                                        <IconButton
                                            aria-owns={userInfoOpen ? "user-info-appbar" : undefined} color="inherit"
                                            aria-haspopup="true" onClick={this.handleUserInfoMenuOpen}>
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
                <Drawer variant="permanent" open={open}
                    className={classNames(classes.drawer, {
                        [classes.drawerOpen]: open,
                        [classes.drawerClose]: !open
                    })}
                    classes={{
                        paper: classNames({
                            [classes.drawerOpen]: open,
                            [classes.drawerClose]: !open
                        })
                    }}>
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
                    <div className={classes.progressOverlayContainer} style={{
                        display: loadingState.isLoading ? "block" : "none"
                    }}>
                        <div className={classes.toolbar}/>
                        <div className={classes.progressOverlay}>
                            <div className={classes.progress}>
                                <CircularProgress className={classes.progressIndicator} thickness={4.5} size={45}/>
                                <div className={classes.progressContent}>
                                    {loadingState.message ? loadingState.message : "Loading"}...
                                </div>
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
    colorGenerator: PropTypes.instanceOf(ColorGenerator).isRequired,
    history: PropTypes.any.isRequired,
    location: PropTypes.shape({
        pathname: PropTypes.string.isRequired
    })
};

export default withStyles(styles, {withTheme: true})(withRouter(withGlobalState(withColor(AppLayout))));
