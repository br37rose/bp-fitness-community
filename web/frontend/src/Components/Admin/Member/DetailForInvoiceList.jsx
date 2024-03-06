import React, { useState, useEffect } from "react";
import { Link, useSearchParams, useParams } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faDumbbell, faGauge, faSearch, faEye, faPencil, faTrashCan, faPlus, faArrowRight, faTable, faArrowUpRightFromSquare, faFilter, faRefresh, faArchive, faUsers } from '@fortawesome/free-solid-svg-icons';
import { useRecoilState } from 'recoil';

import FormErrorBox from "../../Reusable/FormErrorBox";
import { getPaymentProcessorStripeInvoiceListAPI } from "../../../API/PaymentProcessor";
import { getWorkoutSessionDetailAPI } from "../../../API/WorkoutSession";
import { topAlertMessageState, topAlertStatusState, currentUserState, currentWorkoutSessionState } from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import FormSelectFieldForBranch from "../../Reusable/FormSelectFieldForBranch";
import FormSelectFieldForWorkoutProgram from "../../Reusable/FormSelectFieldForWorkoutProgram";
import FormSelectFieldForTrainer from "../../Reusable/FormSelectFieldForTrainer";
import FormDateField from "../../Reusable/FormDateField";
import { PAGE_SIZE_OPTIONS } from "../../../Constants/FieldOptions";
import { DateTime } from "luxon";
import AdminMemberDetailForInvoiceListDesktop from "./DetailForInvoiceListDesktop";
import AdminMemberDetailForInvoiceListMobile from "./DetailForInvoiceListMobile";


function AdminMemberDetailForInvoiceList() {
    ////
    //// URL Parameters.
    ////

    const { bid, uid } = useParams()

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser] = useRecoilState(currentUserState);
    const [currentWorkoutSession, setCurrentWorkoutSession] = useRecoilState(currentWorkoutSessionState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [listData, setListData] = useState("");
    const [isFetching, setFetching] = useState(false);
    const [pageSize, setPageSize] = useState(10);                       // Pagination
    const [previousCursors, setPreviousCursors] = useState([]);         // Pagination
    const [nextCursor, setNextCursor] = useState("");                   // Pagination
    const [currentCursor, setCurrentCursor] = useState("");             // Pagination
    const [showFilter, setShowFilter] = useState(false);                // Filtering + Searching
    const [sortField, setSortField] = useState("created");              // Sorting
    const [temporarySearchText, setTemporarySearchText] = useState(""); // Searching - The search field value as your writes their query.
    const [actualSearchText, setActualSearchText] = useState("");       // Searching - The actual search query value to submit to the API.
    const [branchID, setBranchID] = useState("");                       // Filtering
    const [workoutProgramID, setWorkoutProgramID] = useState("");       // Filtering
    const [trainerID, setTrainerID] = useState("");                     // Filtering
    const [startAt, setStartAt] = useState(null);                       // Filtering

    ////
    //// API.
    ////

    const onWorkoutSessionListSuccess = (response) => {
        console.log("onWorkoutSessionListSuccess: Starting...");
        if (response.results !== null) {
            setListData(response);
            if (response.hasNextPage) {
                setNextCursor(response.nextCursor); // For pagination purposes.
            }
        } else {
            setListData([]);
            setNextCursor("");
        }
    }

    const onWorkoutSessionListError = (apiErr) => {
        console.log("onWorkoutSessionListError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    const onWorkoutSessionListDone = () => {
        console.log("onWorkoutSessionListDone: Starting...");
        setFetching(false);
    }

    ////
    //// Event handling.
    ////

    function onDownloadClick(e, url) {
        console.log("Downloading Invoice PDF...");
        const downloadLink = document.createElement("a");
        downloadLink.href = `${url}`;
        downloadLink.click();
    }

    const fetchList = (userid, cur, limit) => {
        setFetching(true);
        setErrors({});

        getPaymentProcessorStripeInvoiceListAPI(
            userid,
            cur,
            limit,
            onWorkoutSessionListSuccess,
            onWorkoutSessionListError,
            onWorkoutSessionListDone
        );
    }

    const onNextClicked = (e) => {
        console.log("onNextClicked");
        let arr = [...previousCursors];
        arr.push(currentCursor);
        setPreviousCursors(arr);
        setCurrentCursor(nextCursor);
    }

    const onPreviousClicked = (e) => {
        console.log("onPreviousClicked");
        let arr = [...previousCursors];
        const previousCursor = arr.pop();
        setPreviousCursors(arr);
        setCurrentCursor(previousCursor);
    }

    const onSearchButtonClick = (e) => { // Searching
        console.log("Search button clicked...");
        setActualSearchText(temporarySearchText);
    }

    const onWorkoutSessionDetailSuccess = (response) => {
        // For debugging purposes only.
        console.log("onAdminWorkoutProgramSessionDetailSuccess: Starting...");
        console.log(response);
        setCurrentWorkoutSession(response);
    }

    const onWorkoutSessionDetailError = (apiErr) => {
        console.log("onAdminWorkoutProgramSessionDetailError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onAdminWorkoutProgramSessionDetailError: Delayed for 2 seconds.");
            console.log("onAdminWorkoutProgramSessionDetailError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    const onWorkoutSessionDetailDone = () => {
        console.log("onAdminWorkoutProgramSessionDetailDone: Starting...");
        setFetching(false);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.
            fetchList(uid, currentCursor, pageSize);
        }

        return () => { mounted = false; }
    }, [uid, currentCursor, pageSize]);

    ////
    //// Component rendering.
    ////

    return (
        <>
            <div class="container">
                <section class="section">
                    <nav class="breadcrumb" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/admin/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Dashboard</Link></li>
                            <li class=""><Link to="/admin/members" aria-current="page"><FontAwesomeIcon className="fas" icon={faUsers} />&nbsp;Members</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail (Invoices)</Link></li>
                        </ul>
                    </nav>
                    <nav class="box">

                        <div class="columns is-mobile">
                            <div class="column">
                                <h1 class="title is-4"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Member</h1>
                            </div>
                            <div class="column has-text-right">
                                <button onClick={()=>fetchList(uid, currentCursor, pageSize)} class="button is-small is-info" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faRefresh} />
                                </button>
                                &nbsp;
                                {/*
                                <Link to={`/admin/branch/${bid}/class/${pid}/session/${sid}/Invoices/add`} class="button is-small is-primary" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faPlus} />&nbsp;New Invoice
                                </Link>
                                */}
                            </div>
                        </div>

                        {showFilter &&
                            <div class="columns has-background-white-bis" style={{borderRadius:"15px", padding:"20px"}}>
                                <div class="column">
                                    <FormInputFieldWithButton
                                        label={"Search"}
                                        name="temporarySearchText"
                                        type="text"
                                        placeholder="Search by name"
                                        value={temporarySearchText}
                                        helpText=""
                                        onChange={(e)=>setTemporarySearchText(e.target.value)}
                                        isRequired={true}
                                        maxWidth="100%"
                                        buttonLabel={<><FontAwesomeIcon className="fas" icon={faSearch} /></>}
                                        onButtonClick={onSearchButtonClick}
                                    />
                                </div>
                                <div class="column">
                                    <FormSelectFieldForBranch
                                        organizationID={currentUser.organizationID}
                                        branchID={branchID}
                                        setBranchID={setBranchID}
                                        errorText={errors && errors.branchId}
                                        helpText=""
                                    />
                                </div>
                                <div class="column">
                                    <FormSelectFieldForWorkoutProgram
                                        branchID={branchID}
                                        workoutProgramID={workoutProgramID}
                                        setWorkoutProgramID={setWorkoutProgramID}
                                        errorText={errors && errors.workoutProgramId}
                                        disabled={branchID === ""}
                                    />
                                </div>
                                <div class="column">
                                    <FormSelectFieldForTrainer
                                        branchID={branchID}
                                        trainerID={trainerID}
                                        setTrainerID={setTrainerID}
                                        errorText={errors && errors.trainerId}
                                        disabled={branchID === ""}
                                    />
                                </div>
                                <div class="column">
                                    <FormDateField
                                        label="Start At"
                                        name="startAt"
                                        placeholder="Text input"
                                        value={startAt}
                                        errorText={errors && errors.startAt}
                                        helpText=""
                                        onChange={(date)=>setStartAt(date)}
                                        isRequired={true}
                                        maxWidth="120px"
                                    />
                                </div>
                            </div>
                        }

                        {isFetching
                            ? <>
                                <PageLoadingContent displayMessage={"Please wait..."} />
                            </>
                            : <>
                                <div class= "tabs is-medium is-size-7-mobile">
                                  <ul>
                                    <li>
                                        <Link to={`/admin/branch/${bid}/member/${uid}`}>Detail</Link>
                                    </li>
                                    <li>
                                        <Link to={`/admin/branch/${bid}/member/${uid}/bookings`}>Bookings</Link>
                                    </li>
                                    <li>
                                        <Link to={`/admin/branch/${bid}/member/${uid}/waitlist`}>Waitlist</Link>
                                    </li>
                                    <li class="is-active">
                                        <Link>Invoices</Link>
                                    </li>
                                  </ul>
                                </div>

                                {listData && listData.results && (listData.results.length > 0 || previousCursors.length > 0)
                                    ?
                                    <>
                                        <FormErrorBox errors={errors} />
                                        <div class="container">

                                            {/*
                                                ##################################################################
                                                EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                                                ##################################################################
                                            */}
                                            <div class="is-hidden-touch" >
                                                <AdminMemberDetailForInvoiceListDesktop
                                                    listData={listData}
                                                    setPageSize={setPageSize}
                                                    pageSize={pageSize}
                                                    previousCursors={previousCursors}
                                                    onPreviousClicked={onPreviousClicked}
                                                    onNextClicked={onNextClicked}
                                                    onDownloadClick={onDownloadClick}
                                                />
                                            </div>

                                            {/*
                                                ###########################################################################
                                                EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                                                ###########################################################################
                                            */}
                                            <div class="is-fullwidth is-hidden-desktop">
                                                <AdminMemberDetailForInvoiceListMobile
                                                    listData={listData}
                                                    setPageSize={setPageSize}
                                                    pageSize={pageSize}
                                                    previousCursors={previousCursors}
                                                    onPreviousClicked={onPreviousClicked}
                                                    onNextClicked={onNextClicked}
                                                    onDownloadClick={onDownloadClick}
                                                />
                                            </div>

                                        </div>
                                    </>
                                    :
                                    <section class="hero is-medium has-background-white-ter">
                                          <div class="hero-body">
                                            <p class="title">
                                                <FontAwesomeIcon className="fas" icon={faTable} />&nbsp;No Invoice
                                            </p>
                                            <p class="subtitle">
                                                Member has no invoices.
                                            </p>
                                          </div>
                                    </section>
                                }
                            </>
                        }
                    </nav>
                </section>
            </div>
        </>
    );
}

export default AdminMemberDetailForInvoiceList;
