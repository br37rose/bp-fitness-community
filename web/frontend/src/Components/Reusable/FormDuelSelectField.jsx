// Version 1.0.0
import React from "react";


/*
    // EXAMPLE USAGE:

    // Step 1: Create your options.
    const FEET_OPTIONS = [
        { value: 0, label: '-' },
        { value: 7, label: '7 feet'' },
        { value: 6, label: '6 feet' },
        { value: 5, label: '5 feet' }
    ];
    const INCHES_OPTIONS = [
        { value: 0, label: '0\"' },
        { value: 1, label: '1\"' },
        { value: 2, label: '2\"' },
        { value: 3, label: '3\"' },
        { value: 4, label: '4\"' },
        { value: 5, label: '5\"' },
        { value: 6, label: '6\"' },
        { value: 7, label: '7\"' },
        { value: 8, label: '8\"' },
        { value: 9, label: '9\"' },
        { value: 10, label: '10\"' },
        { value: 11, label: '11\"' },
        { value: 12, label: '12\"' }
    ];

    // ...


    // Step 2: Create your state.
    const [heightFeet, setHeightFeet] = useState(0);
    const [heightInches, setHeightInches] = useState(0);

    // ...

    // Step 3: Use this componen.
    <FormDuelSelectField
        label="Please enter your height"
        oneName="heightFeet"
        onePlaceholder="Pick"
        oneSelectedValue={heightFeet}
        oneErrorText={errors && errors.heightFeet}
        oneOnChange={(e) => setHeightFeet(e.target.value)}
        oneOptions={FEET_WITH_EMPTY_OPTIONS}
        oneDisabled={false}
        oneMaxWidth={{maxWidth:"100px"}}
        twoLabel="Height"
        twoName="heightInches"
        twoPlaceholder="Pick"
        twoSelectedValue={heightInches}
        twoErrorText={errors && errors.heightInches}
        twoOnChange={(e) => setHeightInches(e.target.value)}
        twoOptions={INCHES_WITH_EMPTY_OPTIONS}
        twoDisabled={false}
        twoMaxWidth={{maxWidth:"100px"}}
        helpText={heightFeet > -1 && heightInches > -1 && <>(Your height is {heightFeet} ft and {heightInches} inches)</>}
    />

 */
function FormDuelSelectField({
    label,
    oneName, onePlaceholder, oneSelectedValue, oneErrorText, oneValidationText, oneOnChange, oneOptions, oneDisabled, oneMaxWidth,
    twoName, twoPlaceholder, twoSelectedValue, twoErrorText, twoValidationText, twoOnChange, twoOptions, twoDisabled, twoMaxWidth,
    helpText,
}) {
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <div class="columns is-mobile">
                    <div class="column" style={oneMaxWidth}>
                        <span class="select is-fullwidth">
                            <select class={`input ${oneErrorText && 'is-danger'} ${oneSelectedValue && 'is-success'} has-text-black`}
                                     name={oneName}
                              placeholder={onePlaceholder}
                                 onChange={oneOnChange}
                                 disabled={oneDisabled}>
                                {oneOptions.map(function(option, i){
                                    return <option selected={oneSelectedValue === option.value} value={option.value}>{option.label}</option>;
                                })}
                            </select>
                        </span>
                        {oneErrorText &&
                            <p class="help is-danger">{oneErrorText}</p>
                        }
                    </div>
                    <div class="column" style={twoMaxWidth}>
                        <span class="select is-fullwidth">
                            <select class={`input ${twoErrorText && 'is-danger'} ${twoSelectedValue && 'is-success'} has-text-black`}
                                     name={twoName}
                              placeholder={twoPlaceholder}
                                 onChange={twoOnChange}
                                 disabled={twoDisabled}>
                                {twoOptions.map(function(option, i){
                                    return <option selected={twoSelectedValue === option.value} value={option.value}>{option.label}</option>;
                                })}
                            </select>
                        </span>
                        {twoErrorText &&
                            <p class="help is-danger">{twoErrorText}</p>
                        }
                    </div>
                </div>
            </div>
            {helpText &&
                <p class="help">{helpText}</p>
            }
        </div>
    );
}

export default FormDuelSelectField;
