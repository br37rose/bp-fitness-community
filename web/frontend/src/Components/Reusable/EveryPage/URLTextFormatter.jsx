import React from "react";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowUpRightFromSquare } from '@fortawesome/free-solid-svg-icons';


function URLTextFormatter(props) {
    const { urlKey, urlValue, type=""} = props;
    return (
        <>
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
        </>
    );
}

export default URLTextFormatter;
