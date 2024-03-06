// Version 1.0.0
import React from "react";

/*
    // USAGE EXAMPLE

    // STEP 1: Create our state
    const [gender, setGender] = useState("");

    // ...

    // STEP 2: Use our component
    <DataDisplayRowRadio
        label="Gender"
        value={datum.gender}
        opt1Value="Male"
        opt1Label="Male"
        opt2Value="Female"
        opt2Label="Female"
        opt3Value="Other"
        opt3Label="Other"
    />
*/
function DataDisplayRowRadio(props) {
    const {
        label,          // The text to display the user.
        value,          // The selected value.
        opt1Value,      // The value to use for option #1.
        opt1Label,      // The label to display for option #1.
        opt2Value,      // ...
        opt2Label,      // ...
        opt3Value,      // ...
        opt3Label,      // ...
        opt4Value,      // ...
        opt4Label,      // ...
        opt5Value,      // ...
        opt5Label,      // ...
        opt6Value,      // ...
        opt6Label,      // ...
        opt7Value,      // ...
        opt7Label,      // ...
        errorText,      // The error message to display
        wasValidated,   // Boolean indicates if this element was successfully validated or not.
        helpText,       // The special help task to include.
    } = props;
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                {opt1Value === value && <p>{opt1Label}</p>}
                {opt2Value === value && <p>{opt2Label}</p>}
                {opt3Value === value && <p>{opt3Label}</p>}
                {opt4Value === value && <p>{opt4Label}</p>}
                {opt5Value === value && <p>{opt5Label}</p>}
                {opt6Value === value && <p>{opt6Label}</p>}
                {opt7Value === value && <p>{opt7Label}</p>}
                {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default DataDisplayRowRadio;
