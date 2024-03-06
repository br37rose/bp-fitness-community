import React from "react";
import { CountryDropdown } from 'react-country-region-selector';


function FormCountrySelectField({ label, name, placeholder, selectedCountry, errorText, validationText, helpText, onChange, disabled, maxWidth, priorityOptions }) {
    // DEVELOPERS NOTE:
    // https://github.com/country-regions/react-country-region-selector
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control" style={{maxWidth:maxWidth}}>
                <span class="select">
                <CountryDropdown
                    priorityOptions={priorityOptions}
                    name={name}
                    placeholder={placeholder}
                    disabled={disabled}
                    class={`input ${errorText && 'is-danger'} ${validationText && 'is-success'} has-text-black`}
                    value={selectedCountry}
                    onChange={onChange}
                />
                </span>
            </div>
            {helpText &&
                <p class="help">{helpText}</p>
            }
            {errorText &&
                <p class="help is-danger">{errorText}</p>
            }
        </div>
    );
}

export default FormCountrySelectField;
