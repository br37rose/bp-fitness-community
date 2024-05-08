import React, { useState, useEffect } from "react";
import { Link, useParams, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faUsers,
  faEye,
  faArrowLeft,
  faGauge,
  faQuestionCircle,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { getMemberDetailAPI, putMemberUpdateAPI } from "../../../API/member";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { SelectableOption } from "../../Reusable/Wizard/Questions";
import { getQuestionnaireListApi } from "../../../API/questionnaire";

function AdminMemberProfileDetail() {
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
  const [datum, setDatum] = useState({});

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [listData, setListData] = useState("");
  const [answers, setAnswers] = useState([]);

  ////
  //// Event handling.
  ////

  ////
  //// API.
  ////

  // --- Detail --- //

  function onAccountDetailSuccess(response) {
    setDatum(response);
    if (response?.onboardingAnswers?.length > 0) {
      const userAnswers = {};
      response.onboardingAnswers.forEach((answer) => {
        userAnswers[answer.questionId] = answer.answers || [];
      });
      setAnswers(userAnswers);
    }
  }

  function onAccountDetailError(apiErr) {
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onAccountDetailDone() {
    setFetching(false);
  }

  // --- Update --- //
  const handleSubmit = () => {
    setFetching(true);
    setErrors({});

    const onboardingAnswers = listData.results.map((question) => ({
      question_id: question.id,
      answers: answers[question.id] || [],
    }));
    const decamelizedData = {
      id: datum.id,
      organization_id: datum.organizationId,
      first_name: datum.firstName,
      last_name: datum.lastName,
      email: datum.email,
      phone: datum.phone,
      postal_code: datum.postalCode,
      address_line_1: datum.addressLine1,
      address_line_2: datum.addressLine2,
      city: datum.city,
      region: datum.region,
      country: datum.country,
      status: datum.status,
      password: datum.password,
      password_repeated: datum.passwordRepeated,
      how_did_you_hear_about_us: datum.howDidYouHearAboutUs,
      how_did_you_hear_about_us_other: datum.howDidYouHearAboutUsOther,
      agree_promotions_email: datum.agreePromotionsEmail,
      onboarding_answers: onboardingAnswers,
      onboarding_completed: true,
    };
    putMemberUpdateAPI(
      decamelizedData,
      onAdminMemberUpdateSuccess,
      onAdminMemberUpdateError,
      onAdminMemberUpdateDone
    );
  };

  function onAdminMemberUpdateSuccess(response) {
    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Member updated");
    setTopAlertStatus("Workout Member");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Redirect the user to a new page.
    // setForceURL("/admin/member/" + response.id);
  }

  function onAdminMemberUpdateError(apiErr) {
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onAdminMemberUpdateDone() {
    setFetching(false);
  }
  const handleSelect = (questionId, selectedId, isMultiSelect) => {
    if (isMultiSelect) {
      const updatedSelections = answers[questionId]?.includes(selectedId)
        ? answers[questionId].filter((id) => id !== selectedId)
        : [...(answers[questionId] || []), selectedId];
      setAnswers({ ...answers, [questionId]: updatedSelections });
    } else {
      setAnswers({ ...answers, [questionId]: [selectedId] }); // Wrap selectedId in an array
    }
  };
  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      setFetching(true);
      setErrors({});
      getMemberDetailAPI(
        id,
        onAccountDetailSuccess,
        onAccountDetailError,
        onAccountDetailDone
      );
    }

    return () => {
      mounted = false;
    };
  }, []);

  useEffect(() => {
    if (!isFetching && !listData && Object.keys(datum).length > 0) {
      fetchList();
    }
    return () => {};
  }, [datum]);

  const fetchList = () => {
    setFetching(true);
    setErrors({});

    let params = new Map();

    params.set("status", true);
    getQuestionnaireListApi(params, OnListSuccess, OnListErr, onListDone);
  };
  function OnListSuccess(response) {
    if (response.results !== null) {
      setListData(response);
    } else {
      setListData([]);
    }
  }

  function OnListErr(apiErr) {
    setErrors(apiErr);
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onListDone() {
    setFetching(false);
  }

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

          {/* Page */}
          <nav class="box">
            {/* Title + Options */}
            {datum && (
              <div class="columns">
                <div class="column">
                  <p class="title is-4">
                    <FontAwesomeIcon className="fas" icon={faUsers} />
                    &nbsp;Member
                  </p>
                </div>
              </div>
            )}

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Please wait..."} />
            ) : (
              <>
                {datum && (
                  <div class="container">
                    {/* Tab Navigation */}
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li>
                          <Link to={`/admin/member/${id}`}>Detail</Link>
                        </li>
                        <li>
                          <Link to={`/admin/member/${id}/tags`}>Tags</Link>
                        </li>
                        <li className="is-active">
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

                    <p class="title is-6">
                      <FontAwesomeIcon
                        className="fas"
                        icon={faQuestionCircle}
                      />
                      &nbsp;Questions
                    </p>
                    <hr />

                    {/* data here */}
                    {listData &&
                    listData.results &&
                    listData.results.length > 0 ? (
                      <div>
                        {listData.results.map((datum, i) => (
                          <div className="mb-6">
                            <div>
                              <h1 className="has-text-centered is-size-4">
                                <span className="mr-4">{i + 1}.</span>
                                {datum.title}
                              </h1>
                              {datum.subtitle && (
                                <h2 className="subtitle is-size-5 has-text-centered mt-1">
                                  {datum.subtitle}
                                </h2>
                              )}
                            </div>
                            <div className="columns mt-3">
                              <div className="column">
                                {answers &&
                                  datum.options.map((option) => (
                                    <SelectableOption
                                      isFullScreen={false}
                                      key={option}
                                      option={option}
                                      isSelected={
                                        Array.isArray(answers[datum.id])
                                          ? answers[datum.id].includes(option)
                                          : answers[datum.id] === option
                                      }
                                      onSelect={() =>
                                        handleSelect(
                                          datum.id,
                                          option,
                                          datum.isMultiselect
                                        )
                                      }
                                    />
                                  ))}
                              </div>
                            </div>
                          </div>
                        ))}

                        <div className="column has-text-right">
                          <button
                            className="button is-success"
                            onClick={handleSubmit}
                          >
                            Submit&nbsp;
                          </button>
                        </div>
                      </div>
                    ) : (
                      <section className="hero is-medium has-background-white-ter">
                        <div className="hero-body">
                          <p className="title">
                            <FontAwesomeIcon
                              className="fas"
                              icon={faQuestionCircle}
                            />
                            &nbsp;No Questions Available
                          </p>
                          <p className="subtitle">
                            There are currently no questions available.&nbsp;
                          </p>
                        </div>
                      </section>
                    )}

                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link
                          class="button is-medium is-fullwidth-mobile"
                          to={"/dashboard"}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Dashboard
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

export default AdminMemberProfileDetail;
