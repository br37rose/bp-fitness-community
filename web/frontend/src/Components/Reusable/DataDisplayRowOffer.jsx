import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";

import { getOfferSelectOptionListAPI } from "../../API/Offer";
import { getSelectedOptions } from "../../Helpers/selectHelper";


function DataDisplayRowOffer(props) {

    ////
    //// Props.
    ////

    const {
        label="Offer",
        offerID,
        helpText,
    } = props;

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [offerOptions, setOfferOptions] = useState([]);
    const [offerOption, setOfferOption] = useState(null);

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
        setOfferOptions(b);

        // STEP 3: Get all the selected options.
        const offers = getSelectedOptions(b, [offerID]);

        // For debugging purposes only.
        // console.log("tagOptions:", b);
        // console.log("offerID:", [offerID]);
        // console.log("so:", so);

        // STEP 4: Save the selected tag options.
        if (offers && offers.length > 0) {
            setOfferOption(offers[0]);
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
            getOfferSelectOptionListAPI(
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
                    {offerOption !== undefined && offerOption !== null && offerOption !== "" && <>{offerOption.label}</>}
                </p>
                {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default DataDisplayRowOffer;
