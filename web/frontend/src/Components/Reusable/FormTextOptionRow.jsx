import React from "react";

function FormTextOptionRow(props) {
    const { label, selectedValue, options, helpText } = props;
    console.log("selectedValue",selectedValue);
    console.log("options",options);

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

export default FormTextOptionRow;
