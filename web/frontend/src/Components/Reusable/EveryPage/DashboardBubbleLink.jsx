import React, { useState } from "react";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

function DashboardBubbleLink({ title, subtitle, faIcon, url, bgColour = "has-background-success-dark", notificationCount }) {
    return (
        <Link to={url} className="has-text-centered">
            <div className="has-text-centered" style={{ display: "flex", flexDirection: "column", alignItems: "center" }}>
                <div className={`mdi has-text-white ${bgColour}`} style={{ width: "200px", height: "200px", borderRadius: "50%", padding: "10px", position: "relative", display: "flex", justifyContent: "center", alignItems: "center" }}>
                    {notificationCount !== undefined && notificationCount !== null && notificationCount !== "" && notificationCount !== 0 &&
                        <div style={{ position: "absolute", top: "4px", right: "-5px", backgroundColor: "red", color: "white", borderRadius: "50%", padding: "10px" }}>
                            <span style={{ fontSize: "26px" }}>{notificationCount}</span>
                        </div>
                    }
                    <FontAwesomeIcon icon={faIcon} style={{ fontSize: "100px", zIndex: 1 }} />
                </div>
                <h1 className="title is-3 pt-3">{title}</h1>
                <p className="has-text-grey">{subtitle}</p>
            </div>
        </Link>
    );
}

export default DashboardBubbleLink;
