import React from "react";

function VehicleTypesTextFormatter(props) {

    ////
    //// Props.
    ////

    const {
        vehicleTypes=[],
    } = props;


    ////
    //// Component rendering.
    ////

    return (
        <>
            {vehicleTypes && vehicleTypes.map(function(datum, i){
                return <span class="tag is-success mr-2 mb-2">{datum.name}</span>;
            })}
        </>
    );
}

export default VehicleTypesTextFormatter;
