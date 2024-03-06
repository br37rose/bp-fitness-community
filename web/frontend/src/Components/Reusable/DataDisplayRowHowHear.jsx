import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";

import { getHowHearAboutUsItemSelectOptionListAPI } from "../../API/HowHearAboutUsItem";
import { getSelectedOptions } from "../../Helpers/selectHelper";


function DataDisplayRowHowHearAboutUsItem(props) {

    ////
    //// Props.
    ////

    const {
        label="How did you hear about us?",
        howDidYouHearAboutUsID,
        helpText,
    } = props;

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [hhOptions, setHhOptions] = useState([]);
    const [hhOption, setHHOption] = useState(null);

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
        setHhOptions(b);

        // STEP 3: Get all the selected options.
        const hhs = getSelectedOptions(b, [howDidYouHearAboutUsID]);

        // For debugging purposes only.
        // console.log("tagOptions:", b);
        // console.log("hhID:", [howDidYouHearAboutUsID]);
        // console.log("so:", so);

        // STEP 4: Save the selected tag options.
        if (hhs && hhs.length > 0) {
            setHHOption(hhs[0]);
        }
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
            getHowHearAboutUsItemSelectOptionListAPI(
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
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <p>
                    {hhOption !== undefined && hhOption !== null && hhOption !== "" && <>{hhOption.label}</>}
                </p>
                {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default DataDisplayRowHowHearAboutUsItem;
