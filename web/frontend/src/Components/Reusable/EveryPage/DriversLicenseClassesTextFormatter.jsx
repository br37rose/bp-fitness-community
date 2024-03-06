import React from "react";


function DriversLicenseClassesTextFormatter(props) {

    ////
    //// Props.
    ////

    const {
        driversLicenseClass=[],
    } = props;


    ////
    //// Component rendering.
    ////

    return driversLicenseClass; // BUGFIX.

    return (
        <>
            {driversLicenseClass && driversLicenseClass.map(function(datum, i){
                return <span class="tag is-success mr-2 mb-2">{datum.text}</span>;
            })}
        </>
    );
}

export default DriversLicenseClassesTextFormatter;
