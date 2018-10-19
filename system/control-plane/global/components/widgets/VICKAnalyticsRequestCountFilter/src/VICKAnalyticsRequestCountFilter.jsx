/* eslint-disable object-curly-newline,no-unused-vars,no-nested-ternary,react/prop-types,no-restricted-globals,
import/prefer-default-export,react/forbid-prop-types */
/*
 * Copyright (c) 2018, WSO2 Inc. (http://wso2.com) All Rights Reserved.
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

import React from 'react';
import PropTypes from 'prop-types';
import Widget from '@wso2-dashboards/widget';
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import MenuItem from '@material-ui/core/MenuItem';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import CancelIcon from '@material-ui/icons/Cancel';
import ArrowDropUpIcon from '@material-ui/icons/ArrowDropUp';
import ClearIcon from '@material-ui/icons/Clear';
import Chip from '@material-ui/core/Chip';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import Select from 'react-select';
import JssProvider from 'react-jss/lib/JssProvider';
import { Scrollbars } from 'react-custom-scrollbars';

const darkTheme = createMuiTheme({
    palette: {
        type: 'dark',
    },
});

const lightTheme = createMuiTheme({
    palette: {
        type: 'light',
    },
});

const customStyles = {
    input: () => ({
        color: 'white',
    }),
    multiValue: () => ({
        borderRadius: 15,
        display: 'flex',
        flexWrap: 'wrap',
        color: 'black',
        fontSize: '90%',
        overflow: 'hidden',
        paddingLeft: 6,
        textOverflow: 'ellipsis',
        whiteSpace: 'nowrap',
        backgroundColor: 'darkgrey',
        minWidth: '20',
    }),
    singleValue: () => ({
        display: 'flex',
        flexWrap: 'wrap',
        color: 'white',
        fontSize: '95%',
    }),
    control: () => ({
        height: 10,
        borderRadius: 5,
        alignItems: 'center',
        minHeight: 30,
        backgroundColor: 'rgb(51, 51, 51)',
        borderColor: 'grey',
        borderStyle: 'solid',
        borderWidth: 0,
        boxShadow: '0 0 0 1px grey',
        cursor: 'default',
        display: 'flex',
        flexWrap: 'wrap',
        justifyContent: 'space-between',
        outline: '0 !important',
        position: 'relative',
        transition: 'all 100ms',
        paddingTop: 2,
    }),
    option: (styles, { data, isDisabled, isFocused }) => {
        return {
            ...styles,
            height: 30,
            backgroundColor: isDisabled
                ? null
                : isFocused ? 'rgba(255, 255, 255, 0.1)' : null,
        };
    },

    menuList: () => ({
        backgroundColor: 'rgb(51, 51, 51)',
    }),
};

const allOption = [{
    value: 'All',
    label: 'All',
    disabled: false,
}];

/**
 * Options class passed to the react-select component
 */
class Option extends React.Component {
    render() {
        const {
            children, isFocused, onFocus, isDisabled, onSelect, option,
        } = this.props;
        return (
            <MenuItem
                onFocus={onFocus}
                selected={isFocused}
                disabled={isDisabled}
                onClick={() => onSelect(option, event)}
                component="div"
            >
                {children}
            </MenuItem>
        );
    }
}

/**
 * Function to wrap react-select component
 * @param props
 * @returns <Select> componet
 * @constructor
 */
function SelectWrapped(props) {
    const { classes, isClearable, muiTheme, ...other } = props;
    return (
        <Select
            styles={muiTheme.name === 'dark' ? customStyles : {}}
            optionComponent={Option}
            noResultsText={<div>No results found</div>}
            arrowRenderer={(arrowProps) => {
                return arrowProps.isOpen ? <ArrowDropUpIcon /> : <ArrowDropDownIcon />;
            }}
            isClearable={isClearable}
            clearRenderer={() => <ClearIcon />}
            valueComponent={(valueProps) => {
                const { value, children, onRemove } = valueProps;
                const onDelete = (event) => {
                    event.preventDefault();
                    event.stopPropagation();
                    onRemove(value);
                };
                if (onRemove) {
                    return (
                        <Chip
                            tabIndex={-1}
                            label={children}
                            className={classes.chip}
                            deleteIcon={<CancelIcon onTouchEnd={onDelete} />}
                            onDelete={onDelete}
                        />
                    );
                }
                return <div className="Select-value">{children}</div>;
            }}
            {...other}
        />
    );
}

// This is the workaround suggested in https://github.com/marmelab/react-admin/issues/1782
const escapeRegex = /([[\].#*$><+~=|^:(),"'`\s])/g;
let classCounter = 0;

export const generateClassName = (rule, styleSheet) => {
    classCounter += 1;

    if (process.env.NODE_ENV === 'production') {
        return `c${classCounter}`;
    }

    if (styleSheet && styleSheet.options.classNamePrefix) {
        let prefix = styleSheet.options.classNamePrefix;
        // Sanitize the string as will be used to prefix the generated class name.
        prefix = prefix.replace(escapeRegex, '-');

        if (prefix.match(/^Mui/)) {
            return `${prefix}-${rule.key}`;
        }

        return `${prefix}-${rule.key}-${classCounter}`;
    }

    return `${rule.key}-${classCounter}`;
};


/**
 * VICKAnalyticsRequestCountFilter which renders the perspective and filter in home page
 */
class VICKAnalyticsRequestCountFilter extends Widget {
    constructor(props) {
        super(props);
        this.state = {
            width: this.props.glContainer.width,
            height: this.props.glContainer.height,
            perspective: 0,
            cells:[],
            cellOptions: [],
            selectedCellValues: null,
            servers: [],
            serverOptions: [],
            selectedServerValues: null,
            services: [],
            serviceOptions: [],
            selectedServiceValues: null,
            methods: [],
            methodOptions: [],
            selectedMethodValues: null,
            selectedSingleServiceValue: null,
            faultyProviderConf: false,
        };

        this.props.glContainer.on('resize', () => this.setState({
            width: this.props.glContainer.width,
            height: this.props.glContainer.height,
        }));
        this.handleChange = this.handleChange.bind(this);
        this.handleDataReceived = this.handleDataReceived.bind(this);
        this.publishUpdate = this.publishUpdate.bind(this);
    }

    /**
     * Publish user selection to other widgets
     */
    publishUpdate() {
        const filterOptions = {
            perspective: this.state.perspective,
            selectedCellValues: this.state.selectedCellValues,
            selectedServerValues: this.state.selectedServerValues,
            selectedServiceValues: this.state.selectedServiceValues,
            selectedMethodValues: this.state.selectedMethodValues,
            selectedSingleServiceValue: this.state.selectedSingleServiceValue,
        };
        super.publish(filterOptions);
    }

    /**
     * Set the state of the widget after metadata and data is received from SiddhiAppProvider
     * @param data
     */
    handleDataReceived(data) {
        let cells = []; let servers = []; let services = []; let methods = [];
        data.data.forEach((dataUnit) => {
            cells.push(dataUnit[0]);
            servers.push(dataUnit[1]);
            services.push(dataUnit[2]);
            methods.push(dataUnit[3]);
        });
        cells.push('All');
        servers.push('All');
        services.push('All');
        methods.push('All');

        cells = cells.filter((v, i, a) => a.indexOf(v) === i);
        servers = servers.filter((v, i, a) => a.indexOf(v) === i);
        services = services.filter((v, i, a) => a.indexOf(v) === i);
        methods = methods.filter((v, i, a) => a.indexOf(v) === i);

        const cellOptions = cells.map(cell => ({
            value: cell,
            label: cell,
            disabled: false,
        }));

        const serverOptions = servers.map(server => ({
            value: server,
            label: server,
            disabled: false,
        }));

        const serviceOptions = services.map(service => ({
            value: service,
            label: service,
            disabled: false,
        }));


        const methodOptions = methods.map(method => ({
            value: method,
            label: method,
            disabled: false,
        }));

        this.setState({
            cells,
            servers,
            services,
            methods,
            cellOptions,
            serviceOptions,
            serverOptions,
            methodOptions,
            selectedCellValues: allOption,
            selectedServerValues: allOption,
            selectedServiceValues: allOption,
            selectedMethodValues: allOption,
            selectedSingleServiceValue: allOption[0],
        }, this.publishUpdate);
    }

    /**
     * Publish user selection in filters
     * @param values
     */
    handleChange = name => (values) => {
        let options;
        let selectedOptionLabelName;
        let selectedOptionsName;
        let selectedValues;
        if (name === 0) {
            options = this.state.cells;
            selectedOptionLabelName = 'selectedCellValues';
            selectedOptionsName = 'cellOptions';
            selectedValues = values;
        } else if (name === 2) {
            options = this.state.services;
            selectedOptionLabelName = 'selectedServiceValues';
            selectedOptionsName = 'serviceOptions';
            selectedValues = values;
        } else if (name === 1){
            options = this.state.servers;
            selectedOptionLabelName = 'selectedServerValues';
            selectedOptionsName = 'serverOptions';
            selectedValues = values;
        } else{
            options = this.state.methods;
            selectedOptionLabelName = 'selectedMethodValues';
            selectedOptionsName = 'methodOptions';
            selectedValues = values;
        }

        let updatedOptions;
        if (selectedValues.some(value => value.value === 'All')) {
            updatedOptions = options.map(option => ({
                value: option,
                label: option,
                disabled: true,
            }));
            this.setState({
                [selectedOptionLabelName]: [{
                    value: 'All',
                    label: 'All',
                    disabled: false,
                }],
                [selectedOptionsName]: updatedOptions,
            }, this.publishUpdate);
        } else {
            updatedOptions = options.map(option => ({
                value: option,
                label: option,
                disabled: false,
            }));
            this.setState({
                [selectedOptionLabelName]: values,
                [selectedOptionsName]: updatedOptions,
            }, this.publishUpdate);
        }
    };

    componentDidMount() {
        super.getWidgetConfiguration(this.props.widgetID)
            .then((message) => {
                super.getWidgetChannelManager()
                    .subscribeWidget(this.props.id, this.handleDataReceived, message.data.configs.providerConfig);
            })
            .catch((error) => {
                this.setState({
                    faultyProviderConf: true,
                });
            });
    }

    componentWillUnmount() {
        super.getWidgetChannelManager().unsubscribeWidget(this.props.id);
    }

    render() {
        const { classes } = this.props;
        console.log("this.state.perspective "+this.state.perspective+"      "+this.state.selectedServerValues)
        return (
            <JssProvider generateClassName={generateClassName}>
                <MuiThemeProvider theme={this.props.muiTheme.name === 'dark' ? darkTheme : lightTheme}>
                    <Scrollbars style={{ height: this.state.height }}>
                        <div style={{ paddingLeft: 24, paddingRight: 16 }}>
                            <Tabs
                                value={this.state.perspective}
                                onChange={(evt, value) => this.setState({ perspective: value }, this.publishUpdate)}
                            >
                                <Tab label="Cell" />
                                <Tab label="Pods" />
                                <Tab label="Service" />
                                <Tab label="Method" />
                            </Tabs>
                            <div style={{ paddingLeft: 10, paddingRight: 16, paddingTop: 3 }}>
                                {
                                    this.state.perspective === 0
                                    && (
                                        <TextField
                                            fullWidth
                                            value={this.state.selectedCellValues}
                                            onChange={this.handleChange(0)}
                                            placeholder="Filter by Cell"
                                            label=""
                                            InputLabelProps={{
                                                shrink: false,
                                            }}
                                            InputProps={{
                                                inputComponent: SelectWrapped,
                                                inputProps: {
                                                    classes,
                                                    isMulti: true,
                                                    simpleValue: true,
                                                    options: this.state.cellOptions,
                                                    muiTheme: this.props.muiTheme,
                                                    isClearable: true,
                                                },
                                            }}
                                        />
                                    )
                                }
                                {
                                    this.state.perspective === 1
                                    && (
                                        <TextField
                                            fullWidth
                                            value={this.state.selectedServerValues}
                                            onChange={this.handleChange(1)}
                                            placeholder="Filter by Pod"
                                            label=""
                                            InputLabelProps={{
                                                shrink: false,
                                            }}
                                            InputProps={{
                                                inputComponent: SelectWrapped,
                                                inputProps: {
                                                    classes,
                                                    isMulti: true,
                                                    simpleValue: true,
                                                    options: this.state.serverOptions,
                                                    muiTheme: this.props.muiTheme,
                                                    isClearable: true,
                                                },
                                            }}
                                        />
                                    )
                                }
                                {
                                    this.state.perspective === 2
                                    && (
                                        <TextField
                                            fullWidth
                                            value={this.state.selectedServiceValues}
                                            onChange={this.handleChange(2)}
                                            placeholder="Filter by Service"
                                            label=""
                                            InputLabelProps={{
                                                shrink: false,
                                            }}
                                            InputProps={{
                                                inputComponent: SelectWrapped,
                                                inputProps: {
                                                    classes,
                                                    isMulti: true,
                                                    simpleValue: true,
                                                    options: this.state.serviceOptions,
                                                    muiTheme: this.props.muiTheme,
                                                    isClearable: true,
                                                },
                                            }}
                                        />
                                    )
                                }
                                {
                                    this.state.perspective === 3
                                    && (
                                        <TextField
                                            fullWidth
                                            value={this.state.selectedMethodValues}
                                            onChange={this.handleChange(3)}
                                            placeholder="Filter by Method"
                                            label=""
                                            InputLabelProps={{
                                                shrink: false,
                                            }}
                                            InputProps={{
                                                inputComponent: SelectWrapped,
                                                inputProps: {
                                                    classes,
                                                    isMulti: true,
                                                    simpleValue: true,
                                                    options: this.state.methodOptions,
                                                    muiTheme: this.props.muiTheme,
                                                    isClearable: true,
                                                },
                                            }}
                                        />
                                    )
                                }
                            </div>
                        </div>
                    </Scrollbars>
                </MuiThemeProvider>
            </JssProvider>
        );
    }
}

VICKAnalyticsRequestCountFilter.propTypes = {
    classes: PropTypes.object.isRequired,
};
global.dashboard.registerWidget('VICKAnalyticsRequestCountFilter', VICKAnalyticsRequestCountFilter);
