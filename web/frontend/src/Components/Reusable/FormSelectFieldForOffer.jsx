import React, { useState, useEffect } from "react";
import { getOfferSelectOptionListAPI } from "../../API/Offer";

/**
EXAMPLE USAGE:

    <FormOfferField
      offerID={offerID}
      setOfferID={setOfferID}
      offerOther={offerOther}
      setOfferOther={setOfferOther}
      errorText={errors && errors.offerID}
      helpText="Please select the primary gym location this member will be using"
      maxWidth="310px"
      isHidden={true}
    />
*/
function FormSelectFieldForOffer({
    label="Offer",
    offerID,
    setOfferID,
    isOfferOther, // This variable controls whether this component detected the `Other` option or not.
    setIsOfferOther,
    errorText,
    validationText,
    helpText,
    disabled,
    isHidden
}) {
    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [offerOptions, setOfferOptions] = useState([]);

    ////
    //// Event handling.
    ////

    const setOfferIDOverride = (offerID) => {
        // CASE 1: "Other" option selected.
        for (let index in offerOptions) {
            let offerOption = offerOptions[index];
            if (offerOption.label === "Other" && offerOption.value === offerID) {
                // console.log("FormSelectFieldForOffer | offerID:", offerID, "| isOfferOther: true");
                setIsOfferOther(true);
                setOfferID(offerID);
                return;
            }
        }

        // CASE 2: Non-"Other" option selected.
        // console.log("FormSelectFieldForOffer | offerID:", offerID, "| isOfferOther: false");
        setIsOfferOther(false);
        setOfferID(offerID);
    }

    ////
    //// API.
    ////

    function onOfferSelectOptionsSuccess(response){
        // console.log("onOfferSelectOptionsSuccess: Starting...");
        let b = [
            {"value": "", "label": "Please select"},
            ...response
        ]
        setOfferOptions(b);

        // Set `isOfferOther` if the user selected the `other` label.
        for (let index in response) {
            let offerOption = response[index];
            if (offerOption.label === "Other" && offerOption.value === offerID) {
                setIsOfferOther(true);
                // console.log("FormSelectFieldForOffer | picked other | offerID:", offerID);
                return;
            }
        }
    }

    function onOfferSelectOptionsError(apiErr) {
        // console.log("onOfferSelectOptionsError: Starting...");
        setErrors(apiErr);
    }

    function onOfferSelectOptionsDone() {
        // console.log("onOfferSelectOptionsDone: Starting...");
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
                onOfferSelectOptionsSuccess,
                onOfferSelectOptionsError,
                onOfferSelectOptionsDone
            );
        }

        return () => { mounted = false; }
    }, []);

    ////
    //// Component rendering.
    ////

    // Render the JSX component.
    return (
        <>
            <div class={`field pb-4 ${isHidden && "is-hidden"}`} key={offerID}>
                <label class="label">{label}</label>
                <div class="control">
                    <span class="select">
                        {offerOptions.length > 0 &&
                            <select class={`input ${errorText && 'is-danger'} ${validationText && 'is-success'} has-text-black`}
                                     name={`offerID`}
                              placeholder={`Pick offer location`}
                                 onChange={(e,c)=>setOfferIDOverride(e.target.value)}
                                 disabled={disabled}>
                                {offerOptions && offerOptions.length > 0 && offerOptions.map(function(option, i){
                                    return <option selected={offerID === option.value} value={option.value} name={option.label}>{option.label}</option>;
                                })}
                            </select>
                        }
                    </span>
                </div>
                {helpText &&
                    <p class="help">{helpText}</p>
                }
                {errorText &&
                    <p class="help is-danger">{errorText}</p>
                }
            </div>
        </>
    );
}

export default FormSelectFieldForOffer;
