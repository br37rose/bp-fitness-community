import React from "react";
import { DateTime } from "luxon";


function FormTextDateTimeRow(props) {
    const { label, value, helpText } = props;
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
            <p>{value ? DateTime.fromISO(value).toLocaleString(DateTime.DATETIME_MED) : "-"}</p>
            {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default FormTextDateTimeRow;
