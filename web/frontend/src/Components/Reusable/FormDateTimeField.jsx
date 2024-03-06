import React from "react";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";

function FormDateTimeField({ label, name, placeholder, value, type="text", errorText, validationText, helpText, onChange, maxWidth, disabled=false, withPortal=true, filterDate=null }) {
    let dt = null;
    if (value === undefined || value === null || value === "") {
        // Do nothing...
    } else {
        const valueMilliseconds = Date.parse(value);
        dt = new Date(valueMilliseconds);
    }

    let classNameText = "input";
    if (errorText) {
        classNameText = "input is-danger";
    }
    
    // SPECIAL THANKS:
    // https://reactdatepicker.com/

    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control" style={{maxWidth:maxWidth}}>
                <DatePicker className={classNameText}
                         selected={dt}
                             name={name}
                      placeholder={placeholder}
                         disabled={disabled}
                       withPortal={withPortal}
                         portalId={name}
                       filterDate={filterDate}
                     autoComplete="off"
                         onChange={(date) => onChange(date)}
                     showTimeSelect
                     dateFormat="MMMM d, yyyy h:mm aa">
                         <div style={{ color: "red" }}>{errorText}</div>
                </DatePicker>
            </div>
            {errorText &&
                <p class="help is-danger">{errorText}</p>
            }
            {helpText &&
                <p class="help">{helpText}</p>
            }
        </div>
    );
}

export default FormDateTimeField;
