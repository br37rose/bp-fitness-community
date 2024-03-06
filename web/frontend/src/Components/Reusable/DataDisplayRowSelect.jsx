import React from "react";

function DataDisplayRowSelect(props) {
    const { label, selectedValue, options, helpText } = props;

    const option = options.find(
        (option) => option.value === selectedValue
    );

    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <p>{option && option.label}</p>
                {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default DataDisplayRowSelect;
