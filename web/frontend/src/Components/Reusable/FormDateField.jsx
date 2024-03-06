// Version 1.0.0
import React from "react";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";

/*
    // EXAMPLE USAGE:

    // Step 1: Create your state.
    const [birthDate, setBirthDate] = useState(null);

    // ...

    // Step 2: Use this componen.
    <FormDateField
        label="Birth Date"
        name="birthDate"
        placeholder="Text input"
        value={birthDate}
        helpText=""
        onChange={(date)=>setBirthDate(date)}
        errorText={errors && errors.birthDate}
        isRequired={true}
        maxWidth="180px"
    />
 */
function FormDateField({ label, name, placeholder, value, type="text", errorText, validationText, helpText, onChange, maxWidth, disabled=false, withPortal=true, filterDate=null, minDate=null, maxDate=null, monthsShown=1 }) {
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
                     autoComplete="off"
                       withPortal={withPortal}
                      isClearable={true}
                showMonthDropdown={true}
                 showYearDropdown={true}
                         portalId={name}
                       filterDate={filterDate}
                          minDate={minDate}
                          maxDate={maxDate}
                      monthsShown={monthsShown}
                     onChange={(date) => onChange(date)}
                       dateFormat="MM/d/yyyy">
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

export default FormDateField;
