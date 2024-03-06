import React from "react";

function SelectTextFormatter(props) {
    const { selectedValue, options, } = props;

    const option = options.find(
        (option) => option.value === selectedValue
    );

    return (
        <>{option && option.label}</>
    );
}

export default SelectTextFormatter;
