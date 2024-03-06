import React, { useState, useEffect } from "react";
import { getTagSelectOptionListAPI } from "../../../API/Tag";
import { getSelectedOptions } from "../../../Helpers/selectHelper";


function TagIDsTextFormatter(props) {

    ////
    //// Props.
    ////

    const {
        tags=[],
    } = props;

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [tagOptions, setTagOptions] = useState([]);
    const [selectedTagOptions, setSelectedTagOptions] = useState([]);

    ////
    //// API.
    ////

    function onSuccess(response){
        // STEP 1: Convert the API responses to be saved.

        // console.log("onTagSelectOptionsSuccess: Starting...");
        let b = [
            // {"value": "", "label": "Please select"},
            ...response
        ]

        // STEP 2: Save tag options.
        setTagOptions(b);

        // STEP 3: Get all the selected options.
        const so = getSelectedOptions(b, tags);

        // For debugging purposes only.
        console.log("tagOptions:", b);
        console.log("tags:", tags);
        console.log("so:", so);

        // STEP 4: Save the selected tag options.
        setSelectedTagOptions(so);
    }

    function onError(apiErr) {
        // console.log("onTagSelectOptionsError: Starting...");
        setErrors(apiErr);
    }

    function onDone() {
        // console.log("onTagSelectOptionsDone: Starting...");
        setFetching(false);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            setFetching(true);
            getTagSelectOptionListAPI(
                onSuccess,
                onError,
                onDone
            );
        }

        return () => { mounted = false; }
    }, []);

    ////
    //// Component rendering.
    ////

    return (
        <>
            {selectedTagOptions && selectedTagOptions.map(function(datum, i){
                return <span class="tag is-success mr-2 mb-2">{datum.label}</span>;
            })}
        </>
    );
}

export default TagIDsTextFormatter;
