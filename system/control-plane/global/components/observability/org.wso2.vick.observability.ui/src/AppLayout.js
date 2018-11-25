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
import AuthUtils from "./pages/common/utils/authUtils";
import BarChart from "@material-ui/icons/BarChart";
import ChevronLeftIcon from "@material-ui/icons/ChevronLeft";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import Collapse from "@material-ui/core/Collapse";
import {ConfigHolder} from "./pages/common/config/configHolder";
import CssBaseline from "@material-ui/core/CssBaseline";
import DesktopWindows from "@material-ui/icons/DesktopWindows";
import Divider from "@material-ui/core/Divider";
import Drawer from "@material-ui/core/Drawer";
import ExpandLess from "@material-ui/icons/ExpandLess";
import ExpandMore from "@material-ui/icons/ExpandMore";
import Grain from "@material-ui/icons/Grain";
import IconButton from "@material-ui/core/IconButton";
import InsertChartOutlined from "@material-ui/icons/InsertChartOutlined";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import Menu from "@material-ui/core/Menu";
import MenuIcon from "@material-ui/icons/Menu";
import MenuItem from "@material-ui/core/MenuItem";
import PropTypes from "prop-types";
import React from "react";
import Timeline from "@material-ui/icons/Timeline";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import classNames from "classnames";
import {withRouter} from "react-router-dom";
import {withStyles} from "@material-ui/core/styles";
import {ConfigConstants, withConfig} from "./pages/common/config";

const drawerWidth = 240;

const styles = (theme) => ({
    root: {
        display: "flex",
        flexGrow: 1
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
        flexGrow: 1,
        padding: theme.spacing.unit * 3
    },
    nested: {
        paddingLeft: theme.spacing.unit * 3
    }
});

class AppLayout extends React.Component {

    constructor(props) {
        super(props);

        this.handleUserInfoMenuOpen = this.handleUserInfoMenuOpen.bind(this);
        this.handleUserInfoClose = this.handleUserInfoClose.bind(this);
        this.handleDrawerOpen = this.handleDrawerOpen.bind(this);
        this.handleDrawerClose = this.handleDrawerClose.bind(this);

        this.state = {
            open: false,
            userInfo: null,
            subMenuOpen: false
        };
    }

    handleUserInfoMenuOpen = (event) => {
        this.setState({userInfo: event.currentTarget});
    };

    handleUserInfoClose = () => {
        this.setState({userInfo: null});
    };

    handleLogout = () => {
        localStorage.removeItem("username");
        ReactDOM.render(<App />, document.getElementById("root"));
    };

    handleDrawerOpen = () => {
        this.setState({open: true});
    };

    handleDrawerClose = () => {
        this.setState({open: false});
    };

    handleClick = () => {
        this.setState((state) => ({subMenuOpen: !state.subMenuOpen}));
    };

    render() {
        const {classes, history, children, theme, config} = this.props;
        const {open, userInfo} = this.state;
        const userInfoOpen = Boolean(userInfo);

        const navigationState = {
            hideBackButton: true
        };
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
                            config.get(ConfigConstants.USER)
                                ? (
                                    <div>
                                        <IconButton
                                            aria-owns={userInfoOpen ? "user-info-appbar" : undefined}
                                            aria-haspopup="true"
                                            onClick={this.handleUserInfoMenuOpen}
                                            color="inherit">
                                            <AccountCircle/>
                                        </IconButton>
                                        <Menu id="user-info-appbar" anchorEl={this.state.userInfo}
                                            anchorOrigin={{
                                                vertical: "top",
                                                horizontal: "right"
                                            }}
                                            transformOrigin={{
                                                vertical: "top",
                                                horizontal: "right"
                                            }}
                                            open={userInfoOpen}
                                            onClose={this.handleUserInfoClose}>
                                            {/* TODO: Implement user login */}
                                            <MenuItem onClick={this.handleUserInfoClose}>
                                                Profile - {config.get(ConfigConstants.USER)}
                                            </MenuItem>
                                            <MenuItem onClick={this.handleUserInfoClose}>
                                                My account
                                            </MenuItem>
                                            <MenuItem onClick={() => {
                                                AuthUtils.signOut(config);
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
                        [classes.drawerOpen]: this.state.open,
                        [classes.drawerClose]: !this.state.open
                    })}
                    classes={{
                        paper: classNames({
                            [classes.drawerOpen]: this.state.open,
                            [classes.drawerClose]: !this.state.open
                        })
                    }}
                    open={this.state.open}>
                    <div className={classes.toolbar}>
                        <IconButton onClick={this.handleDrawerClose}>
                            {theme.direction === "rtl" ? <ChevronRightIcon/> : <ChevronLeftIcon/>}
                        </IconButton>
                    </div>
                    <Divider/>
                    <List>
                        {/* TODO : Change the icons accordingly to the page menu */}
                        <ListItem button key="Overview" onClick={() => history.push("/", navigationState)}>
                            <ListItemIcon><DesktopWindows/></ListItemIcon>
                            <ListItemText primary="Overview"/>
                        </ListItem>
                        <ListItem button key="Cells" onClick={() => history.push("/cells", navigationState)}>
                            <ListItemIcon><Grain/></ListItemIcon>
                            <ListItemText primary="Cells"/>
                        </ListItem>
                        <ListItem button key="Distributed Tracing"
                            onClick={() => history.push("/tracing", navigationState)}>
                            <ListItemIcon><Timeline/></ListItemIcon>
                            <ListItemText primary="Distributed Tracing"/>
                        </ListItem>
                        <ListItem button onClick={this.handleClick}>
                            <ListItemIcon>
                                <InsertChartOutlined/>
                            </ListItemIcon>
                            <ListItemText inset primary="System Metrics"/>
                            {this.state.subMenuOpen ? <ExpandLess/> : <ExpandMore/>}
                        </ListItem>
                        <Collapse in={this.state.subMenuOpen} timeout="auto" unmountOnExit>
                            <List component="div" disablePadding>
                                <ListItem button className={classes.nested} key="ControlPlane"
                                    onClick={() => history.push("/system-metrics/control-plane", navigationState)}>
                                    <ListItemIcon>
                                        <BarChart/>
                                    </ListItemIcon>
                                    <ListItemText inset primary="Global Control Plane"/>
                                </ListItem>
                                <ListItem button className={classes.nested} key="PodUsage"
                                    onClick={() => history.push("/system-metrics/pod-usage", navigationState)}>
                                    <ListItemIcon>
                                        <BarChart/>
                                    </ListItemIcon>
                                    <ListItemText inset primary="Pod Usage"/>
                                </ListItem>
                                <ListItem button className={classes.nested} key="NodeUsage"
                                    onClick={() => history.push("/system-metrics/node-usage", navigationState)}>
                                    <ListItemIcon>
                                        <BarChart/>
                                    </ListItemIcon>
                                    <ListItemText inset primary="Node Usage"/>
                                </ListItem>
                            </List>
                        </Collapse>
                    </List>
                </Drawer>
                <main className={classes.content}>
                    <div className={classes.toolbar}/>
                    {children}
                </main>
            </div>
        );
    }

}

AppLayout.propTypes = {
    classes: PropTypes.object.isRequired,
    children: PropTypes.any.isRequired,
    theme: PropTypes.object.isRequired,
    config: PropTypes.instanceOf(ConfigHolder).isRequired,
    history: PropTypes.any.isRequired
};

export default withStyles(styles, {withTheme: true})(withRouter(withConfig(AppLayout)));
