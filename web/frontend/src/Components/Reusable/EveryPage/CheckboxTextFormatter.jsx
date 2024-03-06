import React from "react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTimesCircle, faCheckCircle, } from '@fortawesome/free-solid-svg-icons'

function CheckboxTextFormatter(props) {
    const { checked } = props;
    if (checked) {
        return (
            <><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Yes</>
        );
    } else {
        return (
            <><FontAwesomeIcon className="fas" icon={faTimesCircle} />&nbsp;No</>
        );
    }
}

export default CheckboxTextFormatter;
