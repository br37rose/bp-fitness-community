import React, { useState, useEffect } from "react";


/**
  Example:
  
  <DataDisplayRowMultiSelect
      label="Do you identify as belonging to any of the following groups?"
      selectedValues={addAssociate.identifyAs}
      options={IDENTIFY_AS_OPTIONS}
  />
 */
function DataDisplayRowMultiSelect(props) {

    ////
    //// Props.
    ////

    const {
        label="",
        selectedValues=[],
        options=[],
        helpText=""
    } = props;


    useEffect(() => {
        let mounted = true;

        if (mounted) {
        }

        return () => { mounted = false; }
    }, []);

    ////
    //// Component rendering.
    ////

    // Iterate through all the options and select the options vased on the `value`.
    const selectedOptions = options.filter((option) => selectedValues.includes(option.value));

    //TODO: IMPLEMENT
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <p>
                    {selectedOptions && selectedOptions.map(function(datum, i){
                        return <span class="tag is-success mr-2 mb-2">{datum.label}</span>;
                    })}
                </p>
                {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default DataDisplayRowMultiSelect;
