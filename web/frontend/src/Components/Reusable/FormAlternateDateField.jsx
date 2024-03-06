// Version 1.0.0
import React, { useState, useEffect } from "react";


/*
    // EXAMPLE USAGE:


    // Step 1: Create your state.
    const [birthday, setBirthday] = useState(null)

    // ...

    // Step 2: Use this componen.
    <FormAlternateDateField
        label="Birthday"
        name="birthday"
        placeholder="Text input"
        value={birthday}
        helpText=""
        onChange={(date)=>setBirthDay(date)}
        errorText={errors && errors.birthday}
        isRequired={true}
        maxWidth="180px"
        maxDate={new Date()}
    />

    // ...

    // NOTES:
    // Special thanks to the following URL:
    // https://www.smashingmagazine.com/2021/05/frustrating-design-patterns-birthday-picker/
 */
function FormAlternateDateField({
    label,
    name,
    placeholder,
    value,
    errorText,
    validationText,
    helpText,
    onChange,
    maxWidth,
    disabled=false,
    withPortal=true,
    filterDate=null,
    minDate=null,
    maxDate=null,
    monthsShown=1
}) {

    ////
    //// Component states.
    ////

    const [day, setDay] = useState(null);
    const [month, setMonth] = useState(null);
    const [year, setYear] = useState(null);

    useEffect(() => {
        if (value === '0001-01-01T00:00:00Z') {
            // Handle the case where 'value' is '0001-01-01T00:00:00Z'
            // You can avoid rendering or set the component to an initial state.
            return;
        }

        let dt = null;

        if (value !== undefined && value !== null && value !== "") {
          const valueDate = new Date(value);
          valueDate.setHours(0, 0, 0, 0);
          dt = valueDate;
        }

        setDay(dt ? dt.getDate() : null);
        setMonth(dt ? dt.getMonth() + 1 : null);
        setYear(dt ? dt.getFullYear() : null);
    }, [value]);

    ////
    //// Event handling.
    ////

    // Utility funciton to check if number is actually a number.
    const isNumeric = (str) => {
          // Use a regular expression to check if the string contains only numbers
          return /^[0-9]+$/.test(str);
    }

    const onDayChange = (d) => {
        // Step 1: Convert into an integer.
        const di = parseInt(d);
        if (!isNumeric(d)) {
            setDay("");
            return;
        }

        // Prevent futuredates if max date exists.
        if (maxDate) {
            if (year >= maxDate.getFullYear() && month >= (maxDate.getMonth()+1) && di >= maxDate.getDate() ){
                setDay(maxDate.getDate());
                return;
            }
        }

        // Step 2: Defensive Code
        if (di < 1) {
            setDay(1);
            onDateChange(1, month, year);
            return;
        }
        if (di > 31) {
            setDay(31);
            onDateChange(31, month, year);
            return;
        }
        // Step 3: Set day.
        setDay(di);
        onDateChange(di, month, year);
    }

    const onMonthChange = (m) => {
        // Step 1: Convert into an integer.
        const mi = parseInt(m);
        if (!isNumeric(m)) {
            setMonth("");
            return;
        }

        // Prevent futuredates if max date exists.
        if (maxDate) {
            if (year >= maxDate.getFullYear() && mi >= (maxDate.getMonth()+1) ){
                setMonth(maxDate.getMonth()+1);
                return;
            }
        }

        // Step 2: Defensive Code
        if (mi < 1) {
            setMonth(1);
            onDateChange(day, 1, year);
            return;
        }
        if (mi > 12) {
            setMonth(12);
            onDateChange(day, 12, year);
            return;
        }

        // Step 3: Set month.
        setMonth(mi);
        onDateChange(day, mi, year);
    }

    const onYearChange = (y) => {
        // Convert into an integer.
        const yi = parseInt(y);

        // If user entered a value which is not a number, then remove.
        if (!isNumeric(y)) {
            setYear("");
            return;
        }

        // Prevent futuredates if max date exists.
        if (maxDate) {
            if (y > maxDate.getFullYear()){
                setYear(maxDate.getFullYear());
                return;
            }
        }

        // Save the validated year.
        setYear(yi);

        // Update the main date.
        onDateChange(day, month, yi);
    }

    const onDateChange = (d, m, y) => {
        if (d > 0 && m > 0 && y > 0) {
            if (y > 1000) {
                // Convert day, month, and year into a JavaScript Date object
                const date = new Date(y, m - 1, d); // Months are 0-based, so subtract 1 from the month

                console.log("Date formatted successfully:", date);
                onChange(date); // Set the date up the component to the parent because we successfully generated a date.
            }
        }
    }

    ////
    //// Rendering
    ////

    let classNameText = "input";
    if (errorText) {
        classNameText = "input is-danger";
    }

    return (
        <div className="field pb-4">
            <label className="label">{label}</label>
            <div className="control">
                <div className="columns is-mobile">

                    <div className="column" style={{maxWidth:"75px"}}>
                        <span className="is-fullwidth">
                            Day
                            <input className={classNameText} type="text" placeholder="DD" maxlength="2" value={day} onChange={(e)=>{onDayChange(e.target.value)}} />
                        </span>

                    </div>
                    <div className="column" style={{maxWidth:"75px"}}>
                        <span className="is-fullwidth">
                            Month
                            <input className={classNameText} type="text" placeholder="MM" maxlength="2" value={month} onChange={(e)=>{onMonthChange(e.target.value)}} />
                        </span>

                    </div>
                    <div className="column" style={{maxWidth:"90px"}}>
                        <span className="is-fullwidth">
                            Year
                            <input className={classNameText} type="text" placeholder="YYYY" maxlength="4" value={year} onChange={(e)=>{onYearChange(e.target.value)}} />
                        </span>
                    </div>
                </div>
            </div>
            {helpText &&
                <p className="help">{helpText}</p>
            }
            {errorText &&
                <p className="help is-danger">{errorText}</p>
            }
        </div>
    );
}

export default FormAlternateDateField;
