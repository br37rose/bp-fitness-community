import React from "react";
/*
#################
EXAMPLE OF USAGE:
#################

// STEP 1
import FormRadioField from "../Reusable/FormRadioField";

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

    <FormRadioField
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

function FormRadioField({
    label,          // The text to display the user.
    name,           // The element HTML name.
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
    onChange        // The function to call when a selection occurs.
}) {

    console.log(typeof(value), "----->", typeof(opt1Label))
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                {opt1Label &&
                    <>
                        <label class="radio">
                            <input type="radio"
                             data-label={opt1Value}
                                checked={value === opt1Value}
                                   name={name}
                                  value={opt1Value}
                               onChange={onChange} />&nbsp;
                                    {errorText
                                        ? <span style={{color:"#f14668"}} >{opt1Label}</span>
                                        : <span style={wasValidated
                                            ? {color:"#48c78e"}
                                            : {color:"#363636"} }>{opt1Label}</span>
                                    }
                        </label>
                        &nbsp;&nbsp;
                    </>
                }
                {opt2Label &&
                    <>
                        <label class="radio">
                            <input type="radio"
                             data-label={opt2Value}
                                checked={value === opt2Value}
                                   name={name}
                                  value={opt2Value}
                               onChange={onChange} />&nbsp;
                                    {errorText
                                        ? <span style={{color:"#f14668"}} >{opt2Label}</span>
                                        : <span style={wasValidated
                                            ? {color:"#48c78e"}
                                            : {color:"#363636"} }>{opt2Label}</span>
                                    }
                        </label>
                        &nbsp;&nbsp;
                    </>
                }
                {opt3Label &&
                    <>
                        <label class="radio">
                            <input type="radio"
                             data-label={opt3Value}
                                checked={value === opt3Value}
                                   name={name}
                                  value={opt3Value}
                               onChange={onChange} />&nbsp;
                                    {errorText
                                        ? <span style={{color:"#f14668"}} >{opt3Label}</span>
                                        : <span style={wasValidated
                                            ? {color:"#48c78e"}
                                            : {color:"#363636"} }>{opt3Label}</span>
                                    }
                        </label>
                        &nbsp;&nbsp;
                    </>
                }
                {opt4Label &&
                    <>
                        <label class="radio">
                            <input type="radio"
                             data-label={opt4Value}
                                checked={value === opt4Value}
                                   name={name}
                                  value={opt4Value}
                               onChange={onChange} />&nbsp;
                                    {errorText
                                        ? <span style={{color:"#f14668"}} >{opt4Label}</span>
                                        : <span style={wasValidated
                                            ? {color:"#48c78e"}
                                            : {color:"#363636"} }>{opt4Label}</span>
                                    }
                        </label>
                        &nbsp;&nbsp;
                    </>
                }
                {opt5Label &&
                    <>
                        <label class="radio">
                            <input type="radio"
                             data-label={opt5Value}
                                checked={value === opt5Value}
                                   name={name}
                                  value={opt5Value}
                               onChange={onChange} />&nbsp;
                                    {errorText
                                        ? <span style={{color:"#f14668"}} >{opt5Label}</span>
                                        : <span style={wasValidated
                                            ? {color:"#48c78e"}
                                            : {color:"#363636"} }>{opt5Label}</span>
                                    }
                        </label>
                        &nbsp;&nbsp;
                    </>
                }
                {opt6Label &&
                    <>
                        <label class="radio">
                            <input type="radio"
                             data-label={opt6Value}
                                checked={value === opt6Value}
                                   name={name}
                                  value={opt6Value}
                               onChange={onChange} />&nbsp;
                                    {errorText
                                        ? <span style={{color:"#f14668"}} >{opt6Label}</span>
                                        : <span style={wasValidated
                                            ? {color:"#48c78e"}
                                            : {color:"#363636"} }>{opt6Label}</span>
                                    }
                        </label>
                        &nbsp;&nbsp;
                    </>
                }
                {opt7Label &&
                    <>
                        <label class="radio">
                            <input type="radio"
                             data-label={opt7Value}
                                checked={value === opt7Value}
                                   name={name}
                                  value={opt7Value}
                               onChange={onChange} />&nbsp;
                                    {errorText
                                        ? <span style={{color:"#f14668"}} >{opt7Label}</span>
                                        : <span style={wasValidated
                                            ? {color:"#48c78e"}
                                            : {color:"#363636"} }>{opt7Label}</span>
                                    }
                        </label>
                        &nbsp;&nbsp;
                    </>
                }
            </div>
            {errorText &&
                <p class="help is-danger">{errorText}</p>
            }
            <p class="help">{helpText}</p>
        </div>
    );
}

export default FormRadioField;
