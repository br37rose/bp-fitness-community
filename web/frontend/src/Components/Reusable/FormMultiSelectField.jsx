// Version 1.0.0
import React from "react";
import { startCase } from 'lodash';
import Select from 'react-select'

import { getSelectedOptions } from "../../Helpers/selectHelper";


/*
// EXAMPLE USAGE:

// Step 1: Create your options.
export const TAG_OPTIONS = [
    { value: 1, label: 'Education' },
    { value: 2, label: 'Enterainment' },
    { value: 3, label: 'Sporst ' },
];

// ...

// Step 2: Create your state.
const [tags, setTags] = useState([])

// ...

// Step 3: Use this componen.
<FormMultiSelectField
    label="Please pick any of the tags? (Optional)"
    name="tags"
    placeholder="Text input"
    options={TAG_OPTIONS}
    selectedValues={tags}
    onChange={(e)=>{
        let values = [];
        for (let option of e) {
            values.push(option.value);
        }
        setTags(values);
    }}
    errorText={errors && errors.tags}
    helpText=""
    isRequired={true}
    maxWidth="320px"
/>
*/

function FormMultiSelectField({
    label,
    name,
    placeholder,
    options,
    selectedValues,
    onChange,
    errorText,
    validationText,
    helpText,
    maxWidth,
    disabled=false })
{
    let style = {maxWidth:maxWidth};
    if (errorText) {
        style = {maxWidth:maxWidth, borderColor:"red", borderStyle: "solid", borderWidth: "1px"};
    }
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control" style={style}>
                <Select isMulti
                           name="name"
                        options={options}
                          value={getSelectedOptions(options, selectedValues)}
                    isClearable={false}
                       onChange={onChange}
                     isDisabled={disabled}
               menuPortalTarget={document.body}
                         styles={{
                             menuPortal: base => ({ ...base, zIndex: 9999 }),
                             menu: provided => ({ ...provided, zIndex: 9999 })
                         }}
                />
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

export default FormMultiSelectField;

// DEVELOPERS NOTES:
// The `styles` options were taken from this[0] link because there was an error
// that certain GUI overlapped the drowpdown. Therefore we used this library to
// fix our issue.
// [0] https://stackoverflow.com/questions/55830799/how-to-change-zindex-in-react-select-drowpdown
