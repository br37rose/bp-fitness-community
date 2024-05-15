import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faRepeat,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faUserCircle,
  faGauge,
  faPencil,
  faUsers,
  faEye,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faCogs,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import useLocalStorage from "../../../Hooks/useLocalStorage";
import { getMemberDetailAPI, deleteMemberAPI } from "../../../API/member";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import FormCountryField from "../../Reusable/FormCountryField";
import FormRegionField from "../../Reusable/FormRegionField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import {
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
  MEMBER_STATUS_WITH_EMPTY_OPTIONS,
} from "../../../Constants/FieldOptions";
import {
  SUBSCRIPTION_STATUS_WITH_EMPTY_OPTIONS,
  SUBSCRIPTION_TINTERVAL_WITH_EMPTY_OPTIONS,
} from "../../../Constants/FieldOptions";
import FormTextRow from "../../Reusable/FormTextRow";
import FormTextTagRow from "../../Reusable/FormTextTagRow";
import FormTextYesNoRow from "../../Reusable/FormTextYesNoRow";
import FormTextOptionRow from "../../Reusable/FormTextOptionRow";

function AdminMemberDetail() {
  ////
  //// URL Parameters.
  ////

  const { id } = useParams();

  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [datum, setDatum] = useState({});
  const [tabIndex, setTabIndex] = useState(1);
  const [selectedMemberForDeletion, setSelectedMemberForDeletion] =
    useState(null);

  ////
  //// Event handling.
  ////

  const onDeleteConfirmButtonClick = () => {
    deleteMemberAPI(
      selectedMemberForDeletion.id,
      onMemberDeleteSuccess,
      onMemberDeleteError,
      onMemberDeleteDone
    );
    setSelectedMemberForDeletion(null);
  };

  ////
  //// API.
  ////

  // --- Detail --- //

  function onMemberDetailSuccess(response) {
    setDatum(response);
  }

  function onMemberDetailError(apiErr) {
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onMemberDetailDone() {
    setFetching(false);
  }

  // --- Delete --- //

  function onMemberDeleteSuccess(response) {
    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("Member deleted");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Redirect back to the members page.
    setForceURL("/admin/members");
  }

  function onMemberDeleteError(apiErr) {
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onMemberDeleteDone() {
    setFetching(false);
  }

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      setFetching(true);
      getMemberDetailAPI(
        id,
        onMemberDetailSuccess,
        onMemberDetailError,
        onMemberDetailDone
      );
    }

    return () => {
      mounted = false;
    };
  }, []);
  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div class="container">
        <section class="section">
          {/* Desktop Breadcrumbs */}
          <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/admin/members" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faUsers} />
                  &nbsp;Members
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/admin/members" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Members
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modal */}
          <nav>
            {/* Delete modal */}
            <div
              class={`modal ${
                selectedMemberForDeletion !== null ? "is-active" : ""
              }`}
            >
              <div class="modal-background"></div>
              <div class="modal-card">
                <header class="modal-card-head">
                  <p class="modal-card-title">Are you sure?</p>
                  <button
                    class="delete"
                    aria-label="close"
                    onClick={(e, ses) => setSelectedMemberForDeletion(null)}
                  ></button>
                </header>
                <section class="modal-card-body">
                  You are about to delete this member and all the data
                  associated with it. This action is cannot be undone. Are you
                  sure you would like to continue?
                </section>
                <footer class="modal-card-foot">
                  <button
                    class="button is-success"
                    onClick={onDeleteConfirmButtonClick}
                  >
                    Confirm
                  </button>
                  <button
                    class="button"
                    onClick={(e, ses) => setSelectedMemberForDeletion(null)}
                  >
                    Cancel
                  </button>
                </footer>
              </div>
            </div>
          </nav>

          {/* Page */}
          <nav class="box">
            {datum && (
              <div class="columns">
                <div class="column">
                  <p class="title is-4">
                    <FontAwesomeIcon className="fas" icon={faUsers} />
                    &nbsp;Member
                  </p>
                </div>
                <div class="column has-text-right">
                  <Link
                    to={`/admin/member/${id}/update`}
                    class="button is-warning is-fullwidth-mobile is-small"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPencil} />
                    &nbsp;Edit Member
                  </Link>
                  &nbsp;
                  <Link
                    onClick={(e, s) => {
                      setSelectedMemberForDeletion(datum);
                    }}
                    class="button is-danger is-fullwidth-mobile is-small"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPencil} />
                    &nbsp;Delete Member
                  </Link>
                </div>
              </div>
            )}
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Please wait..."} />
            ) : (
              <>
                {datum && (
                  <div class="container" key={datum.id}>
                    {/* Tab navigation */}
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li class="is-active">
                          <Link>Detail</Link>
                        </li>
                        <li>
                          <Link to={`/admin/member/${datum.id}/tags`}>
                            Tags
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/member/${datum.id}/profile`}>
                            Profile
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/member/${datum.id}/fitness-plans`}>
                            Fitness Plans
                          </Link>
                        </li>
                        <li>
                          <Link
                            to={`/admin/member/${datum.id}/nutrition-plans`}
                          >
                            Nutrition Plans
                          </Link>
                        </li>
                      </ul>
                    </div>

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faIdCard} />
                      &nbsp;Full-Name
                    </p>
                    <hr />

                    <FormTextRow label="Name" value={datum.name} />

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faContactCard} />
                      &nbsp;Contact Information
                    </p>
                    <hr />

                    <FormTextRow label="Email" value={datum.email} />

                    <FormTextRow label="Phone" value={datum.phone} />

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faAddressBook} />
                      &nbsp;Address
                    </p>
                    <hr />

                    <FormTextRow label="Country" value={datum.country} />

                    <FormTextRow
                      label="Province/Territory"
                      value={datum.region}
                    />

                    <FormTextRow label="City" value={datum.city} />

                    <FormTextRow
                      label="Address Line 1"
                      value={datum.addressLine1}
                    />

                    <FormTextRow
                      label="Address Line 2"
                      value={datum.addressLine2}
                    />

                    <FormTextRow label="Postal Code" value={datum.postalCode} />

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faChartPie} />
                      &nbsp;Metrics
                    </p>
                    <hr />

                    <FormTextOptionRow
                      label="How did you hear about us?"
                      selectedValue={datum.howDidYouHearAboutUs}
                      options={HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS}
                    />

                    {datum.howDidYouHearAboutUs === 1 && (
                      <FormTextRow
                        label="Other (Please specify):"
                        value={datum.howDidYouHearAboutUsOther}
                      />
                    )}

                    <FormTextYesNoRow
                      label="I agree to receive electronic updates from my local gym"
                      value={datum.agreePromotionsEmail}
                    />

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faCogs} />
                      &nbsp;Settings
                    </p>
                    <hr />

                    <FormTextTagRow
                      label="Status"
                      value={datum.status}
                      opt2Value={1}
                      opt2Label="Active"
                      opt4Value={100}
                      opt4Label="Archived"
                      opt4Code={`is-danger`}
                    />

                    {datum !== undefined && datum !== null && datum !== "" && (
                      <>
                        {datum.stripeSubscription !== undefined &&
                          datum.stripeSubscription !== null &&
                          datum.stripeSubscription !== "" && (
                            <>
                              <p class="subtitle is-6">
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faRepeat}
                                />
                                &nbsp;Subscription
                              </p>
                              <hr />
                              <FormTextOptionRow
                                label="Interval"
                                selectedValue={
                                  datum.stripeSubscription.interval
                                }
                                options={
                                  SUBSCRIPTION_TINTERVAL_WITH_EMPTY_OPTIONS
                                }
                              />
                              <FormTextOptionRow
                                label="Status"
                                selectedValue={datum.stripeSubscription.status}
                                options={SUBSCRIPTION_STATUS_WITH_EMPTY_OPTIONS}
                              />
                            </>
                          )}
                      </>
                    )}

                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link
                          class="button is-fullwidth-mobile"
                          to={`/admin/members`}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Members
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        <Link
                          to={`/admin/member/${id}/update`}
                          class="button is-warning is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPencil} />
                          &nbsp;Edit
                        </Link>
                      </div>
                    </div>
                  </div>
                )}
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminMemberDetail;
