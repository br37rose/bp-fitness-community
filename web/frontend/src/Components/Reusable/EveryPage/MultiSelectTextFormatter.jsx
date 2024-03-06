import React from "react";


function MultiSelectTextFormatter(props) {
    const { selectedValues, options, } = props;

    if (selectedValues === undefined || selectedValues === null || selectedValues.length === 0) {
        return null;
    }
    if (options === undefined || options === null || options.length === 0) {
        return null;
    }

    // Iterate through all the options and select the options vased on the `value`.
    const selectedOptions = options.filter((option) => selectedValues.includes(option.value));

    return (
        <>
            {selectedOptions && selectedOptions.map(function(datum, i){
                return <span class="tag is-success mr-2 mb-2">{datum.label}</span>;
            })}
        </>
    );
}

export default MultiSelectTextFormatter;
