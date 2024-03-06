import React, { useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCheckCircle, faInfoCircle, faCircleExclamation } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import { topAlertMessageState, topAlertStatusState } from "../../AppState";


function TopAlertBanner() {
    const [topAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus] = useRecoilState(topAlertStatusState);

    if (topAlertMessage === "") {
        return null;
    } else {
        switch (topAlertStatus) {
            case "success":
                return (
                    <div className="has-background-success topbar is-fullwidth p-2">
                        <div className="has-text-white is-size-4 has-text-centered has-text-weight-bold">
                            <FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;{topAlertMessage}
                        </div>
                    </div>
                );
                break;
            case "danger":
                return (
                    <div className="has-background-danger is-fullwidth p-2">
                        <div className="has-text-white is-size-4 has-text-centered has-text-weight-bold">
                            <FontAwesomeIcon className="fas" icon={faCircleExclamation} />&nbsp;{topAlertMessage}
                        </div>
                    </div>
                );
                break;
            default:
                return (
                    <div className="has-background-primary is-fullwidth p-2">
                        <div className="has-text-white is-size-4 has-text-centered">
                            <FontAwesomeIcon className="fas" icon={faInfoCircle} />&nbsp;{topAlertMessage}
                        </div>
                    </div>
                );
        }
    }

    // Render the following component GUI.

}

export default TopAlertBanner;