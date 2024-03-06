import React from "react";
import { RegionDropdown } from 'react-country-region-selector';


function FormRegionField({ label, name, placeholder, selectedCountry, selectedRegion, errorText, validationText, helpText, onChange, disabled, maxWidth }) {
    // DEVELOPERS NOTE:
    // https://github.com/country-regions/react-country-region-selector
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control" style={{maxWidth:maxWidth}}>
                <span class="select">
                <RegionDropdown
                    name={name}
                    placeholder={placeholder}
                    disabled={disabled}
                    class={`input ${errorText && 'is-danger'} ${validationText && 'is-success'} has-text-black`}
                    country={selectedCountry}
                    value={selectedRegion}
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

export default FormRegionField;
