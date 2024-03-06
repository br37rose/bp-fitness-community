import React, { useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCircleInfo, faCheckCircle, faInfoCircle, faCircleExclamation } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';


function AlertBanner({ message, status }) {


    if (message === "" || status === "") {
        return null;
    } else {
        switch (status) {
            case "success":
                return (
                    <div className="is-fullwidth pb-4">
                        <div className="has-text-white is-size-4 has-text-centered has-text-weight-bold has-background-success  p-2">
                            <FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;{message}
                        </div>
                    </div>
                );
                break;
            case "danger":
                return (
                    <>
                        <div className="is-fullwidth pb-4">
                            <div className="has-text-white is-size-4 has-text-centered has-text-weight-bold has-background-danger p-2">
                                <FontAwesomeIcon className="fas" icon={faCircleExclamation} />&nbsp;{message}
                            </div>
                        </div>
                    </>
                );
                break;
            case "info":
                return (
                    <div className="is-fullwidth pb-4">
                        <div className="has-text-white is-size-4 has-text-centered has-text-weight-bold has-background-info p-2">
                            <FontAwesomeIcon className="fas" icon={faCircleInfo} />&nbsp;{message}
                        </div>
                    </div>
                );
                break;
            default:
                return (
                    <div className="is-fullwidth pb-4">
                        <div className="has-text-white is-size-4 has-text-centered has-text-weight-bold has-background-danger p-2">
                            <FontAwesomeIcon className="fas" icon={faInfoCircle} />&nbsp;{message}
                        </div>
                    </div>
                );
        }
    }

    // Render the following component GUI.

}

export default AlertBanner;
