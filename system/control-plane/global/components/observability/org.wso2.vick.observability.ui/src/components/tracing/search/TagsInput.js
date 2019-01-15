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

import ChipInput from "material-ui-chip-input";
import React from "react";
import * as PropTypes from "prop-types";

class TagsInput extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            tagsTempInput: {
                content: "",
                errorMessage: ""
            },
            tags: props.defaultTags
        };
    }


    render() {
        const {tagsTempInput, tags} = this.state;

        // Generating the chips from the tags object
        const tagChips = [];
        for (const tagKey in tags) {
            if (tags.hasOwnProperty(tagKey)) {
                tagChips.push(`${tagKey}=${tags[tagKey]}`);
            }
        }

        const helperText = tagsTempInput.content ? "Press Enter to Add the Tag" : null;
        return (
            <ChipInput label="Tags" InputLabelProps={{shrink: true}}
                onBeforeAdd={(chip) => Boolean(TagsInput.parseChip(chip))}
                error={Boolean(tagsTempInput.errorMessage)}
                helperText={tagsTempInput.errorMessage ? tagsTempInput.errorMessage : helperText}
                onAdd={this.handleTagAdd}
                placeholder={"Eg: http.status_code=200"} value={tagChips}
                onUpdateInput={this.handleTagsTempInputUpdate}
                inputValue={tagsTempInput.content} onDelete={this.handleTagRemove}
                onBlur={() => this.setState({
                    tagsTempInput: {
                        content: "",
                        errorMessage: ""
                    }
                })}/>
        );
    }

    handleTagsTempInputUpdate = (event) => {
        const value = event.currentTarget.value;
        this.setState({
            tagsTempInput: {
                content: value,
                errorMessage: !value || TagsInput.parseChip(value)
                    ? ""
                    : "Invalid tag filter format. Expected \"tagKey=tagValue\""
            }
        });
    };

    /**
     * Handle a tag being added to the tag filter.
     *
     * @param {string} chip The chip representing the tag that was added
     */
    handleTagAdd = (chip) => {
        const {onTagsUpdate} = this.props;
        const tag = TagsInput.parseChip(chip);
        if (tag) {
            this.setState((prevState) => {
                const newTags = {
                    ...prevState.tags,
                    [tag.key]: tag.value
                };
                if (onTagsUpdate) {
                    onTagsUpdate(newTags);
                }
                return {
                    tags: newTags,
                    tagsTempInput: {
                        ...prevState.tagsTempInput,
                        content: "",
                        errorMessage: ""
                    }
                };
            });
        }
    };

    /**
     * Handle a tag being removed from the tag filter.
     *
     * @param {string} chip The chip representing the tag that was removed
     */
    handleTagRemove = (chip) => {
        const {onTagsUpdate} = this.props;
        const tag = TagsInput.parseChip(chip);
        if (tag) {
            this.setState((prevState) => {
                const newTags = {...prevState.tags};
                Reflect.deleteProperty(newTags, tag.key);
                if (onTagsUpdate) {
                    onTagsUpdate(newTags);
                }
                return {
                    tags: newTags
                };
            });
        }
    };

    static parseChip = (chip) => {
        let tag = null;
        if (chip) {
            const chipContent = chip.split("=");
            if (chipContent.length === 2 && chipContent[0] && chipContent[1]) {
                tag = {
                    key: chipContent[0].trim(),
                    value: chipContent[1].trim()
                };
            }
        }
        return tag;
    };

}

TagsInput.propTypes = {
    onTagsUpdate: PropTypes.func.isRequired,
    defaultTags: PropTypes.object.isRequired
};

export default TagsInput;
