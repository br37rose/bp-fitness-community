import React from "react";

function FormTextYesNoRow(props) {
    const { label, checked, helpText } = props;
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <p>{checked ? "Yes" : "No"}</p>
                {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default FormTextYesNoRow;
