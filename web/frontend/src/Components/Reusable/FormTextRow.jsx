import React from "react";

function FormTextRow(props) {
    const { label, value, helpText } = props;
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <p>{value ? value : "-"}</p>
                {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default FormTextRow;
