import React from "react";
import { startCase } from 'lodash';

/*
#################
EXAMPLE OF USAGE:
#################

// STEP 1
import FormCheckboxField from "../Element/FormCheckboxField";

// STEP 2
const [errors, setErrors] = useState({
    // "okToEmail": "This field is required",
});

// STEP 3
const [validation, setValidation] = useState({
    "okToEmail": false,
});

// STEP 4
const [okToEmail, setOKToEmail] = useState("");

// STEP 5
function onOKToEmailChange(e) {
    const label = e.target.dataset.label; // Note: 'dataset' is a react data via https://stackoverflow.com/a/20383295
    setOKToEmail(label === true || label === "true");
    validation["okToEmail"] = false // Clear validation
    setValidation(validation);
    // setErrors(errors["email"]="");
}

// STEP 6
function onSubmit() {
    if (organizationName !== "Over55" && organizationName !== "" && organizationType !== "" && firstName !== "" && lastName !== "" && primaryPhone !== "" && email !== "" && okToEmail !== "" && okToText !== "") {
        setForceURL("/setup-employer-profile-step-1"); // (For example)
    } else {
        var newErrors = {};
        var newValidation = {};

        // ...

        if (okToEmail === "") {
            newErrors["okToEmail"] = "missing choice";
        } else {
            newValidation["okToEmail"] = true
        }

        // ...

        setErrors(newErrors);
        setValidation(newValidation);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();

        // For debugging purposes only.
        console.log(newErrors)
    }
}

// Step 7:
return (
    // ...

    <FormCheckboxField
        label="OK TO EMAIL? (OPTIONAL)"
        name="okToEmail"
        value={okToEmail}
        opt1Value={true}
        opt1Label="Yes"
        opt2Value={false}
        opt2Label="No"
        errorText={errors && errors.okToEmail}
        wasValidated={validation && validation.okToEmail}
        helpText="By selecting YES, important communication will occur through email"
        onChange={onOKToEmailChange}
    />

    // ...
)
*/

function FormCheckboxField({
    label,          // The text to display the user.
    name,           // The element HTML name.
    checked,        // The value to use for option.
    errorText,      // The error message to display
    wasValidated,   // Boolean indicates if this element was successfully validated or not.
    helpText,       // The special help task to include.
    onChange,       // The function to call when a selection occurs.
    disabled,
    paddingBottom = "pb-4"         //
}) {
    return (
        <div class={`field ${paddingBottom}`}>

            <div class="control">

                <label class="label checkbox">
                    <input type="checkbox"
                        checked={checked}
                        name={name}
                        disabled={disabled}
                        onChange={onChange} />&nbsp;
                    {errorText
                        ? <span style={{ color: "#f14668" }} >{label}</span>
                        : <span style={wasValidated
                            ? { color: "#48c78e" }
                            : { color: "#363636" }}>&nbsp;<strong>{label}</strong></span>
                    }
                </label>

            </div>
            {errorText &&
                <p class="help is-danger">{errorText}</p>
            }
            <p class="help">{helpText}</p>
        </div>
    );
}

export default FormCheckboxField;
