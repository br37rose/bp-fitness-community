import React, {useState} from "react";
import { startCase } from 'lodash';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faMinus } from '@fortawesome/free-solid-svg-icons';
/*
    DATA-STRUCTURE
    ---------------
    data - needs to be an array of strings. For example:
        ["hi", "how are you?", "thanks"].

    FUNCTIONS
    ---------------
    onListInputFieldChange - needs to look something like this:
        const onListInputFieldChange = (e, i) => {
            // For debugging purposes.
            console.log(e, i);

            // Make a copy of the "array of strings" into a mutable array.
            const copyOfArr = [...previewDescription];

            // Update record.
            copyOfArr[i] = e

            // Save to our react state.
            setPreviewDescription(copyOfArr);
        }

    onRemoveListInputFieldChange - this function is used to remove user selected data:
        const onRemoveListInputFieldChange = (i) => {
            // For debugging purposes.
            console.log(i);

            // Make a copy of the "array of strings" into a mutable array.
            const copyOfArr = [...previewDescription];

            // Delete record.
            const x = copyOfArr.splice(i, 1);

            // For debugging purposes.
            console.log(x);

            // Save to our react state.
            setPreviewDescription(copyOfArr);
        }

    onAddListInputFieldClick - this function is used to add data:
        const onAddListInputFieldClick = () => {
            // For debugging purposes.
            console.log("add");

            // Make a copy of the "array of strings" into a mutable array.
            const copyOfArr = [...previewDescription];

            // Add empty record.
            copyOfArr.push("");

            // For debugging purposes.
            console.log(copyOfArr);

            // Save to our react state.
            setPreviewDescription(copyOfArr);
        }


*/
function FormListInputField({
    label, name, placeholder, value, type="text", errorText, validationText, helpText,
    onListInputFieldChange, onRemoveListInputFieldChange, onAddListInputFieldClick, maxWidth, disabled=false
}) {
    let classNameText = "input";
    if (errorText) {
        classNameText = "input is-danger";
    }

    return (
        <div class="pb-4">
            <label class="label">{label}&nbsp;
                <button class="button is-success is-small" onClick={onAddListInputFieldClick} disabled={disabled}><FontAwesomeIcon className="fas" icon={faPlus} /></button>
            </label>
            {value && value.map(function(datum, i){
                return <div class="field has-addons pb-4" key={i}  style={{maxWidth:maxWidth}}>
                    <p class="control is-expanded">
                        <input class={`${classNameText} is-centered`}
                                name={name}
                                type={type}
                         placeholder={placeholder}
                               value={datum}
                            onChange={(e)=>onListInputFieldChange(e.target.value, i)}
                            disabled={disabled}
                        autoComplete="off"
                            disabled={disabled} />
                    </p>
                    <p class="control">
                        <button class="button is-danger" onClick={(e)=>onRemoveListInputFieldChange(i)} disabled={disabled}>
                            <FontAwesomeIcon className="fas" icon={faMinus} />
                        </button>
                    </p>
                </div>
            })}
            {errorText &&
                <p class="help is-danger">{errorText}</p>
            }
            {helpText &&
                <p class="help">{helpText}</p>
            }
        </div>
    );
}

export default FormListInputField;
