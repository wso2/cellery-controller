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
import ChevronLeftIcon from "@material-ui/icons/ChevronLeft";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import CssBaseline from "@material-ui/core/CssBaseline";
import Divider from "@material-ui/core/Divider";
import Drawer from "@material-ui/core/Drawer";
import IconButton from "@material-ui/core/IconButton";
import InboxIcon from "@material-ui/icons/MoveToInbox";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import Menu from "@material-ui/core/Menu";
import MenuIcon from "@material-ui/icons/Menu";
import MenuItem from "@material-ui/core/MenuItem";
import PropTypes from "prop-types";
import React from "react";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import Utils from "./utils";
import classNames from "classnames";
import {withRouter} from "react-router-dom";
import {withStyles} from "@material-ui/core/styles";
import {ConfigConstants, withConfig} from "./pages/utils/config";

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
    drawerPaper: {
        position: "relative",
        whiteSpace: "nowrap",
        width: drawerWidth,
        transition: theme.transitions.create("width", {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen
        })
    },
    drawerPaperClose: {
        overflowX: "hidden",
        transition: theme.transitions.create("width", {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen
        }),
        width: theme.spacing.unit * 7,
        [theme.breakpoints.up("sm")]: {
            width: theme.spacing.unit * 9
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
            userInfo: null
        };
    }

    handleUserInfoMenuOpen = (event) => {
        this.setState({userInfo: event.currentTarget});
    };

    handleUserInfoClose = () => {
        this.setState({userInfo: null});
    };

    handleDrawerOpen = () => {
        this.setState({open: true});
    };

    handleDrawerClose = () => {
        this.setState({open: false});
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
                                            <MenuItem onClick={() => Utils.signOut()}>
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
                    classes={{
                        paper: classNames(classes.drawerPaper, !open && classes.drawerPaperClose)
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
                            <ListItemIcon><InboxIcon/></ListItemIcon>
                            <ListItemText primary="Overview"/>
                        </ListItem>
                        <ListItem button key="Cells" onClick={() => history.push("/cells", navigationState)}>
                            <ListItemIcon><InboxIcon/></ListItemIcon>
                            <ListItemText primary="Cells"/>
                        </ListItem>
                        <ListItem button key="Micro Services"
                            onClick={() => history.push("/microservices", navigationState)}>
                            <ListItemIcon><InboxIcon/></ListItemIcon>
                            <ListItemText primary="Micro Services"/>
                        </ListItem>
                        <ListItem button key="Distributed Tracing"
                            onClick={() => history.push("/tracing", navigationState)}>
                            <ListItemIcon><InboxIcon/></ListItemIcon>
                            <ListItemText primary="Distributed Tracing"/>
                        </ListItem>
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
    config: PropTypes.any.isRequired,
    history: PropTypes.any.isRequired
};

export default withStyles(styles, {withTheme: true})(withRouter(withConfig(AppLayout)));
