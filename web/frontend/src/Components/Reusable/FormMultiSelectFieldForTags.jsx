import React, { useState, useEffect } from "react";
import { startCase } from 'lodash';
import Select from 'react-select'

import { getTagSelectOptionListAPI } from "../../API/Tag";
import { getSelectedOptions } from "../../Helpers/selectHelper";


function FormMultiSelectFieldForTags({
    label="Tags (Optional)",
    name="tags",
    placeholder="Please select tags",
    tenantID,
    tags,
    setTags,
    errorText,
    validationText,
    helpText,
    maxWidth,
    disabled=false })
{
    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [tagSelectOptions, setTagSelectOptions] = useState([]);

    ////
    //// API.
    ////

    function onTagSelectOptionsSuccess(response){
        // console.log("onTagSelectOptionsSuccess: Starting...");
        let b = [
            // {"value": "", "label": "Please select"},
            ...response
        ]
        setTagSelectOptions(b);
    }

    function onTagSelectOptionsError(apiErr) {
        // console.log("onTagSelectOptionsError: Starting...");
        setErrors(apiErr);
    }

    function onTagSelectOptionsDone() {
        // console.log("onTagSelectOptionsDone: Starting...");
        setFetching(false);
    }

    ////
    //// Event handling.
    ////

    const onTagsChange = (e) => {
        // console.log("onTagsChange, e:",e); // For debugging purposes only.
        let values = [];
        for (let option of e) {
            // console.log("option:",option); // For debugging purposes only.
            values.push(option.value);
        }
        // console.log("onTagsChange, values:",values); // For debugging purposes only.
        setTags(values);
    }


    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            setFetching(true);
            getTagSelectOptionListAPI(
                onTagSelectOptionsSuccess,
                onTagSelectOptionsError,
                onTagSelectOptionsDone
            );
        }

        return () => { mounted = false; }
    }, []);

    ////
    //// Component rendering.
    ////

    let style = {maxWidth:maxWidth};
    if (errorText) {
        style = {maxWidth:maxWidth, borderColor:"red", borderStyle: "solid", borderWidth: "1px"};
    }
    return (
        <div className="field pb-4">
            <label className="label">{label}</label>
            <div className="control" style={style}>
                <Select isMulti
                    placeholder={placeholder}
                           name="tags"
                        options={tagSelectOptions}
                          value={getSelectedOptions(tagSelectOptions, tags)}
                    isClearable={false}
                       onChange={onTagsChange}
                     isDisabled={disabled}
                     isLoading={isFetching}
                />
            </div>
            {errorText &&
                <p className="help is-danger">{errorText}</p>
            }
            {helpText &&
                <p className="help">{helpText}</p>
            }
        </div>


    );
}

export default FormMultiSelectFieldForTags;
