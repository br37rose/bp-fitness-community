import React from "react";

function FormInputFieldWithButton({ label, name, placeholder, value, type="text", errorText, validationText, helpText, onChange, maxWidth, disabled=false, onButtonClick, buttonLabel }) {
    let classNameText = "input";
    if (errorText) {
        classNameText = "input is-danger";
    }
    return (
        <>
        <label class="label">{label}:</label>
        <div class="field has-addons pb-4">

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
            <div class="control">
                <button class="button is-info" onClick={onButtonClick}>{buttonLabel}</button>
            </div>
            {errorText &&
                <p class="help is-danger">{errorText}</p>
            }
            {helpText &&
                <p class="help">{helpText}</p>
            }
        </div>
        </>
    );
}

export default FormInputFieldWithButton;
