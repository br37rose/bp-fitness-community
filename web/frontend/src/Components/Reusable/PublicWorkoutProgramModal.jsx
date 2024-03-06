import React, { useState, useEffect } from "react";
import { DateTime } from "luxon";

import { getPublicWorkoutProgramDetailAPI } from "../../API/workout_program";


function PublicWorkoutProgramModal({ selectedWorkoutProgramID, setSelectedWorkoutProgramID }) {
    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [workoutProgramDetail, setWorkoutProgramDetail] = useState(null);

    ////
    //// Event handling.
    ////


    ////
    //// API.
    ////

    function onWorkoutProgramDetailSuccess(response){
        console.log("onWorkoutProgramDetailSuccess: Starting...");
        setWorkoutProgramDetail(response);
        setErrors({});
    }

    function onWorkoutProgramDetailError(apiErr) {
        console.log("onWorkoutProgramDetailError: Starting...");
        setErrors(apiErr);
    }

    function onWorkoutProgramDetailDone() {
        console.log("onWorkoutProgramDetailDone: Starting...");
        setFetching(false);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.

            // Do not fetch API if we do not have anything selected.
            if (selectedWorkoutProgramID !== undefined && selectedWorkoutProgramID !== null && selectedWorkoutProgramID !== "") {
                setFetching(true);
                getPublicWorkoutProgramDetailAPI(
                    selectedWorkoutProgramID,
                    onWorkoutProgramDetailSuccess,
                    onWorkoutProgramDetailError,
                    onWorkoutProgramDetailDone
                );
            }
        }

        return () => { mounted = false; }
    }, [selectedWorkoutProgramID]);

    return (
        <div class={`modal ${selectedWorkoutProgramID ? 'is-active' : ''}`}>
            <div class="modal-background"></div>
            <div class="modal-card">
                <header class="modal-card-head">
                    <p class="modal-card-title">
                    {isFetching
                    ?
                    <>
                    Class
                    </>
                    :
                    <>
                    Class: {workoutProgramDetail && workoutProgramDetail.name}
                    </>
                    }

                    </p>
                    <button class="delete" aria-label="close" onClick={(e,s) => setSelectedWorkoutProgramID("")}></button>
                </header>
                <section class="modal-card-body">
                    {isFetching
                        ?
                        <div class="column has-text-centered is-1">
                            <div class="loader-wrapper is-centered">
                                <br />
                                <div class="loader is-loading" style={{height: "80px", width: "80px"}}></div>
                            </div>
                            <br />
                            <div className="">Fetching...</div>
                            <br />
                        </div>
                        :
                        <>
                        {workoutProgramDetail &&
                            <div class="content">
                                <ul>
                                    <li><strong>Class Type</strong>: {workoutProgramDetail.workoutProgramTypeName}</li>
                                    <li><strong>Description</strong>: {workoutProgramDetail.description}</li>
                                    <li><strong>Start Date</strong>: {DateTime.fromISO(workoutProgramDetail.startAt).toLocaleString(DateTime.DATE_MED)}</li>
                                    <li><strong>Trainer</strong>: {workoutProgramDetail.trainerName}</li>
                                </ul>
                            </div>
                        }
                        </>
                    }

                </section>
                <footer class="modal-card-foot">
                    <button class="button is-success" onClick={(e,s) => setSelectedWorkoutProgramID("")}>Close</button>
                    {/*<button class="button" onClick={(e,s) => setSelectedWorkoutProgramID("")}>Cancel</button>*/}
                </footer>
            </div>
        </div>
    );
}

export default PublicWorkoutProgramModal;
