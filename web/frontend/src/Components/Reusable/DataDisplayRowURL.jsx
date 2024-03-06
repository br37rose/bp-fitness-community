import React from "react";
import { Link } from "react-router-dom";
import { DateTime } from "luxon";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowUpRightFromSquare } from '@fortawesome/free-solid-svg-icons';


function DataDisplayRowText(props) {
    const { label, urlKey, urlValue, helpText, type=""} = props;
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <p>
                    {urlKey && urlValue
                        ?
                        <>
                            {type === "external" &&
                                <Link target="_blank" rel="noreferrer" to={urlValue}>{urlKey}&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} /></Link>
                            }
                            {type !== "external" &&
                                <Link to={urlValue}>{urlKey}</Link>
                            }
                        </>
                        :
                        "-"
                    }
                </p>
                {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default DataDisplayRowText;
