import React from "react";
import { startCase } from 'lodash';


function ErrorBox({ errors }) {
    // STEP 1: Check to see if we have ANY errors returned from the API
    //         web-service and if no errors were returned then our stateless
    //         component will return nothing.
    if (errors === null || errors === undefined) {
        // console.log("ErrorBox | null"); // For debugging purposes only.
        return null;
    }
    if (Object.keys(errors).length === 0) {
        // console.log("ErrorBox | null"); // For debugging purposes only.
        return null;
    }

    // STEP 2: If the error result is not an object then return the following GUI.
    if (typeof errors !== 'object') {
        return (
            <article class="message is-danger">
                <div class="message-header">
                    Error(s):
                </div>
                <div class="message-body">
                    <p style={{margin: '10px'}}>
                        <strong>Error:</strong>&nbsp;{errors}
                    </p>
                    <hr />
                    <p style={{margin: '10px'}}><i>Please make sure the above error(s) have been fixed before submitting again</i></p>
                </div>
            </article>
        );
    }

    // STEP 3: If the result is an object then run the following GUI.
    return (
        <article class="message is-danger">
            <div class="message-header">
                Error(s):
            </div>
            <div class="message-body">
                {Object.keys(errors).map(key => {
                    // STEP 3: Process a single "field" or "non_field_errors" and
                    //         then get our value.
                    let startKey = startCase(key);

                    // DEVELOPERS NOTE:
                    // The following code will remove any "Id" related keys as it was added
                    // due to "Golang" naming convention in the database. Ex: `how_hear_id`.
                    startKey = startKey.replace(" Id", "");

                    let value = errors[key];
                    // console.log(key, value); // For debugging purposes only.

                    // STEP 4: Generate the error row if the value accomponying it is not blank.
                    if (value !== undefined && value !== null) {
                        return (
                            <p style={{margin: '10px'}}>
                                <strong>{startKey}:</strong>&nbsp;{value}
                            </p>
                        );
                    }
                    return null;
                })}
                <hr />
                <p style={{margin: '10px'}}><i>Please make sure the above error(s) have been fixed before submitting again</i></p>
            </div>
        </article>
    );
}

export default ErrorBox;
