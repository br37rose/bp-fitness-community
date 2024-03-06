import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faUserCircle,
  faFileInvoiceDollar,
  faReceipt,
  faEye,
  faPencil,
  faTrashCan,
  faPlus,
  faGauge,
  faArrowRight,
  faTable,
  faArrowUpRightFromSquare,
  faRefresh,
  faFilter,
  faSearch,
  faArrowLeft,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { DateTime } from "luxon";

import FormErrorBox from "../../../Reusable/FormErrorBox";
import { getPaymentProcessorStripeInvoiceListAPI } from "../../../../API/PaymentProcessor";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
} from "../../../../AppState";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../../Reusable/FormInputFieldWithButton";
import { PAGE_SIZE_OPTIONS } from "../../../../Constants/FieldOptions";
import AccountInvoiceListDesktop from "./ListDesktop";
import AccountInvoiceListMobile from "./ListMobile";
import Footer from "../../../Menu/Footer";
import Layout from "../../../Menu/Layout";


function AccountInvoiceList() {
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [currentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [listData, setListData] = useState("");
  const [selectedMemberForDeletion, setSelectedMemberForDeletion] =
    useState("");
  const [isFetching, setFetching] = useState(false);
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination
  const [showFilter, setShowFilter] = useState(false); // Filtering + Searching
  const [sortField, setSortField] = useState("created"); // Sorting
  const [temporarySearchText, setTemporarySearchText] = useState(""); // Searching - The search field value as your writes their query.
  const [actualSearchText, setActualSearchText] = useState(""); // Searching - The actual search query value to submit to the API.
  const [branchID, setBranchID] = useState(""); // Filtering

  ////
  //// API.
  ////

  function onMemberListSuccess(response) {
    console.log("onMemberListSuccess: Starting...");
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

  function onMemberListError(apiErr) {
    console.log("onMemberListError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onMemberListDone() {
    console.log("onMemberListDone: Starting...");
    setFetching(false);
  }

  ////
  //// BREADCRUMB
  ////
  const breadcrumbItems = {
    items: [
      { text: 'Dashboard', link: '/dashboard', isActive: false, icon: faGauge },
      { text: 'Account', link: '/account', icon: faUserCircle, isActive: false },
      { text: 'Subscription', link: '#', icon: faFileInvoiceDollar, isActive: true }
    ],
    mobileBackLinkItems: {
      link: "/account",
      text: "Back to Account",
      icon: faArrowLeft
    }
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
      onMemberListSuccess,
      onMemberListError,
      onMemberListDone
    );
  }

  const onNextClicked = (e) => {
    console.log("onNextClicked");
    let arr = [...previousCursors];
    arr.push(currentCursor);
    setPreviousCursors(arr);
    setCurrentCursor(nextCursor);
  };

  const onPreviousClicked = (e) => {
    console.log("onPreviousClicked");
    let arr = [...previousCursors];
    const previousCursor = arr.pop();
    setPreviousCursors(arr);
    setCurrentCursor(previousCursor);
  };

  const onSearchButtonClick = (e) => {
    // Searching
    console.log("Search button clicked...");
    setActualSearchText(temporarySearchText);
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      fetchList(currentUser.id, currentCursor, pageSize);
    }

    return () => {
      mounted = false;
    };
  }, [currentUser, currentCursor, pageSize]);

  ////
  //// Component rendering.
  ////

  return (
    <Layout breadcrumbItems={breadcrumbItems}>
      <div className="box">
        <div className="columns is-mobile">
          <div className="column">
            <h1 class="title is-4 is-hidden-touch"><FontAwesomeIcon className="fas" icon={faReceipt} />
              &nbsp;Past Invoices
            </h1>
            <h1 class="title is-6 is-hidden-desktop mt-2"><FontAwesomeIcon className="fas" icon={faReceipt} />
              &nbsp;Past Invoices
            </h1>
          </div>
          {/*
              <div className="column has-text-right">
                <button
                  onClick={() =>
                    fetchList(
                      currentCursor,
                      pageSize,
                      actualSearchText,
                      branchID
                    )
                  }
                  class="button is-small is-info"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faRefresh} />
                </button>
                &nbsp;
                <button
                  onClick={(e) => setShowFilter(!showFilter)}
                  class="button is-small is-success"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faFilter} />
                  &nbsp;Filter
                </button>
              </div>
              */}
        </div>

        {/*
            {showFilter && (
              <div
                class="columns has-background-white-bis"
                style={{ borderRadius: "15px", padding: "20px" }}
              >
                <div class="column">
                  <FormInputFieldWithButton
                    label={"Search"}
                    name="temporarySearchText"
                    type="text"
                    placeholder="Search by name"
                    value={temporarySearchText}
                    helpText=""
                    onChange={(e) => setTemporarySearchText(e.target.value)}
                    isRequired={true}
                    maxWidth="100%"
                    buttonLabel={
                      <>
                        <FontAwesomeIcon className="fas" icon={faSearch} />
                      </>
                    }
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
              </div>
            )}
            */}

        {isFetching ? (
          <PageLoadingContent displayMessage={"Please wait..."} />
        ) : (
          <>
            <FormErrorBox errors={errors} />
            {listData &&
              listData.results &&
              (listData.results.length > 0 || previousCursors.length > 0) ? (
              <div className="container">

                {/*
                        ##################################################################
                        EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                        ##################################################################
                    */}
                <div class="is-hidden-touch" >
                  <AccountInvoiceListDesktop
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
                  <AccountInvoiceListMobile
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
            ) : (
              <section className="hero is-medium has-background-white-ter">
                <div className="hero-body">
                  <p className="title">
                    <FontAwesomeIcon className="fas" icon={faTable} />
                    &nbsp;No Invoices
                  </p>
                  <p className="subtitle">
                    No invoices.{" "}
                    <b>
                      <Link to="/subscriptions">
                        Click here&nbsp;
                        <FontAwesomeIcon
                          className="mdi"
                          icon={faArrowRight}
                        />
                      </Link>
                    </b>{" "}
                    to get started enrolling in a subscription.
                  </p>
                </div>
              </section>
            )}
          </>
        )}

        <div class="columns pt-5">
          <div class="column is-half">
            <Link class="button is-medium is-fullwidth-mobile" to={"/account/subscription"}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Subscription</Link>
          </div>
          <div class="column is-half has-text-right">
          </div>
        </div>

      </div>
    </Layout>
  );
}

export default AccountInvoiceList;
