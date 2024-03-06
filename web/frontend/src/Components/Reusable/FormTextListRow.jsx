import React from "react";

function FormTextListRow(props) {
    const { label, value, helpText } = props;
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                {value && value.map(function(line, i){
                    return <>
                        <p>{line}</p>
                    </>
                })}
                {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default FormTextListRow;
