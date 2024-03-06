import React, { useState, useEffect } from "react";
import { getVideoCategorySelectOptionListAPI } from "../../API/VideoCategory";

/**
EXAMPLE USAGE:

    <FormVideoCategoryField
      videoCategoryID={videoCategoryID}
      setVideoCategoryID={setVideoCategoryID}
      videoCategoryOther={videoCategoryOther}
      setVideoCategoryOther={setVideoCategoryOther}
      errorText={errors && errors.videoCategoryID}
      helpText="Please select the primary gym location this member will be using"
      maxWidth="310px"
      isHidden={true}
    />
*/
function FormSelectFieldForVideoCategory({
    label="Category",
    videoCategoryID,
    setVideoCategoryID,
    isVideoCategoryOther, // This variable controls whether this component detected the `Other` option or not.
    setIsVideoCategoryOther,
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
    const [videoCategoryOptions, setVideoCategoryOptions] = useState([]);

    ////
    //// Event handling.
    ////

    const setVideoCategoryIDOverride = (videoCategoryID) => {
        // CASE 1: "Other" option selected.
        for (let index in videoCategoryOptions) {
            let videoCategoryOption = videoCategoryOptions[index];
            if (videoCategoryOption.label === "Other" && videoCategoryOption.value === videoCategoryID) {
                // console.log("FormSelectFieldForVideoCategory | videoCategoryID:", videoCategoryID, "| isVideoCategoryOther: true");
                setIsVideoCategoryOther(true);
                setVideoCategoryID(videoCategoryID);
                return;
            }
        }

        // CASE 2: Non-"Other" option selected.
        // console.log("FormSelectFieldForVideoCategory | videoCategoryID:", videoCategoryID, "| isVideoCategoryOther: false");
        setIsVideoCategoryOther(false);
        setVideoCategoryID(videoCategoryID);
    }

    ////
    //// API.
    ////

    function onVideoCategorySelectOptionsSuccess(response){
        // console.log("onVideoCategorySelectOptionsSuccess: Starting...");
        let b = [
            {"value": "", "label": "Please select"},
            ...response
        ]
        setVideoCategoryOptions(b);

        // Set `isVideoCategoryOther` if the user selected the `other` label.
        for (let index in response) {
            let videoCategoryOption = response[index];
            if (videoCategoryOption.label === "Other" && videoCategoryOption.value === videoCategoryID) {
                setIsVideoCategoryOther(true);
                // console.log("FormSelectFieldForVideoCategory | picked other | videoCategoryID:", videoCategoryID);
                return;
            }
        }
    }

    function onVideoCategorySelectOptionsError(apiErr) {
        // console.log("onVideoCategorySelectOptionsError: Starting...");
        setErrors(apiErr);
    }

    function onVideoCategorySelectOptionsDone() {
        // console.log("onVideoCategorySelectOptionsDone: Starting...");
        setFetching(false);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            setFetching(true);
            getVideoCategorySelectOptionListAPI(
                onVideoCategorySelectOptionsSuccess,
                onVideoCategorySelectOptionsError,
                onVideoCategorySelectOptionsDone
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
            <div class={`field pb-4 ${isHidden && "is-hidden"}`} key={videoCategoryID}>
                <label class="label">{label}</label>
                <div class="control">
                    <span class="select">
                        {videoCategoryOptions.length > 0 &&
                            <select class={`input ${errorText && 'is-danger'} ${validationText && 'is-success'} has-text-black`}
                                     name={`videoCategoryID`}
                              placeholder={`Pick videoCategory location`}
                                 onChange={(e,c)=>setVideoCategoryIDOverride(e.target.value)}
                                 disabled={disabled}>
                                {videoCategoryOptions && videoCategoryOptions.length > 0 && videoCategoryOptions.map(function(option, i){
                                    return <option selected={videoCategoryID === option.value} value={option.value} name={option.label}>{option.label}</option>;
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

export default FormSelectFieldForVideoCategory;
