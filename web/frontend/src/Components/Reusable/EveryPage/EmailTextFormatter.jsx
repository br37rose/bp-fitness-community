import React from "react";
import { Link } from "react-router-dom";


function EmailTextFormatter(props) {
    const { value } = props;
    return (
        <Link to={`mailto:${value}`}>{value}</Link>
    );
}

export default EmailTextFormatter;
