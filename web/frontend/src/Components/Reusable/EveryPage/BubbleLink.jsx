import React from "react";
import { Link, useLocation } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCircleInfo, faCheckCircle, faInfoCircle, faCircleExclamation } from '@fortawesome/free-solid-svg-icons'

/*
    ########
    Example:
    ########

    // Page Title
    // ...

    // This is your code:
    <section class="hero ">
        <div class="hero-body has-text-centered">
            <div class="container">
                <div class="columns is-vcentered">
                    <div class="column">
                        <BubbleLink
                            title={`Add`}
                            subtitle={`Add clients`}
                            faIcon={faPlus}
                            url={`/admin/clients/add/step-1`}
                            bgColour={`has-background-danger-dark`}
                        />
                    </div>
                    <div className="column">
                        <BubbleLink
                            title={`Search`}
                            subtitle={`Search clients`}
                            faIcon={faSearch}
                            url={`/admin/clients/search`}
                            bgColour={`has-background-success-dark`}
                        />
                    </div>
                </div>
            </div>
        </div>
    </section>

    // Page Table
    // ...
*/
function BubbleLink({ title, subtitle, faIcon, url, bgColour="has-background-success-dark" }) {
    return (
		<Link to={url} className="has-text-centered">
			<div className="has-text-centered" style={{ display: "flex", flexDirection: "column", alignItems: "center" }}>
				<div className={`mdi has-text-white ${bgColour}`} style={{ width: "200px", height: "200px", borderRadius: "100%", padding: "10px", display: "flex", alignItems: "center", justifyContent: "center" }}>
					<FontAwesomeIcon icon={faIcon} style={{ fontSize: "100px" }} />
				</div>
				<h1 className="title is-3 pt-3">{title}</h1>
				<p className="has-text-grey">{subtitle}</p>
			</div>
		</Link>
	);
}


export default BubbleLink;
