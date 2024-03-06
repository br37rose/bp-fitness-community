import React from "react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTimesCircle, faCheckCircle, } from '@fortawesome/free-solid-svg-icons'

function DataDisplayRowCheckbox(props) {
    const { label, checked, helpText } = props;
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <p>
                {checked
                    ?
                    <>
                       <FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Yes
                    </>
                    :
                    <>
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />&nbsp;No
                    </>
                }
                </p>
                {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default DataDisplayRowCheckbox;
