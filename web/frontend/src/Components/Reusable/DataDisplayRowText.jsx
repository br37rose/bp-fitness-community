import React from "react";
import { Link } from "react-router-dom";
import PhoneTextFormatter from "./PhoneTextFormatter";
import DateTextFormatter from "./DateTextFormatter";
import DateTimeTextFormatter from "./DateTimeTextFormatter";
import TextWithLineBreaksFormatter from "./TextWithLineBreaksFormatter";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTimesCircle, faCheckCircle, faCircle } from '@fortawesome/free-solid-svg-icons'


function DataDisplayRowText(props) {
    const { label, value, helpText, type="text"} = props;
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <p>
                    {value
                        ?
                        <>
                        {type === "text" &&
                            value
                        }
                        {type === "text_with_linebreaks" &&
                            <><TextWithLineBreaksFormatter text={value} /></>
                        }
                        {type === "email" &&
                            <Link to={`mailto:${value}`}>{value}</Link>
                        }
                        {type === "phone" &&
                            <PhoneTextFormatter value={value} />
                        }
                        {type === "datetime" &&
                            <DateTimeTextFormatter value={value} />
                        }
                        {type === "date" &&
                            <DateTextFormatter value={value} />
                        }
                        {type === "currency" &&
                            <>${value}</>
                        }
                        {type === "textlist" &&
                            <>
                            {value && value.map(function(datum, i){
                                return <div class="pb-3" key={i}><FontAwesomeIcon className="fas" icon={faCircle} />&nbsp;{datum}</div>
                            })}
                            </>
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
