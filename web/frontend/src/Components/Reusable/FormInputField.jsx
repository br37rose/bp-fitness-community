import React from "react";

function FormInputField({ label, name, placeholder, value, type="text", errorText, validationText, helpText, onChange, maxWidth, disabled=false }) {
    let classNameText = "input";
    if (errorText) {
        classNameText = "input is-danger";
    }
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <input class={classNameText}
                        name={name}
                        type={type}
                 placeholder={placeholder}
                       value={value}
                    onChange={onChange}
                       style={{maxWidth:maxWidth}}
                    disabled={disabled}
                autoComplete="off" />
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

export default FormInputField;
