import React from "react";
import { Link } from "react-router-dom";
import parsePhoneNumber from 'libphonenumber-js'


function PhoneTextFormatter(props) {
    const { value } = props;

    // Special thanks to: https://www.npmjs.com/package/libphonenumber-js#difference-from-googles-libphonenumber
    const phoneNumber = parsePhoneNumber(value)
    return (
        <Link to={phoneNumber.getURI()}>{phoneNumber.formatNational()}</Link>
    );
}

export default PhoneTextFormatter;
