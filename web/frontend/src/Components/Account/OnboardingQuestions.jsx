import { useState, useEffect } from "react";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faQuestionCircle, faList } from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import FormErrorBox from "../Reusable/FormErrorBox";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
  questionnaireFilterStatus,
  questionnaireFilterShowState,
  questionnaireFilterSortState,
  quizAnswersState,
} from "../../AppState";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import { getQuestionnaireListApi } from "../../API/questionnaire";
import { SelectableOption, Title } from "../Reusable/Wizard/Questions";
import { putMemberUpdateAPI } from "../../API/member";

function MemberQuestionnaireList() {
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [currentUser] = useRecoilState(currentUserState);
  const [showFilter, setShowFilter] = useRecoilState(
    questionnaireFilterShowState
  ); // Filtering + Searching
  const [status, setStatus] = useRecoilState(questionnaireFilterStatus);
  const [sort, setSort] = useRecoilState(questionnaireFilterSortState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [listData, setListData] = useState("");

  const [isFetching, setFetching] = useState(false);
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination
  const [answers, setAnswers] = useRecoilState(quizAnswersState);
  // const [forceURL, setForceURL] = useState("");

  ////
  //// API.
  ////

  //   if (currentUser && currentUser.on)

  function OnListSuccess(response) {
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

  function OnListErr(apiErr) {
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onListDone() {
    setFetching(false);
  }

  ////
  //// Event handling.
  ////

  const fetchList = () => {
    setFetching(true);
    setErrors({});

    let params = new Map();

    params.set("status", true);
    getQuestionnaireListApi(params, OnListSuccess, OnListErr, onListDone);
  };

  const handleSubmit = () => {
    setFetching(true);
    setErrors({});

    const onboardingAnswers = listData.results.map((question) => ({
      question_id: question.id,
      answers: answers[question.id] || [],
    }));
    const decamelizedData = {
      id: currentUser.id,
      organization_id: currentUser.organization_id,
      first_name: currentUser.first_name,
      last_name: currentUser.last_name,
      email: currentUser.email,
      phone: currentUser.phone,
      postal_code: currentUser.postal_code,
      address_line_1: currentUser.address_line_1,
      address_line_2: currentUser.address_line_2,
      city: currentUser.city,
      region: currentUser.region,
      country: currentUser.country,
      status: currentUser.status,
      password: currentUser.password,
      password_repeated: currentUser.passwordRepeated,
      how_did_you_hear_about_us: currentUser.how_did_you_hear_about_us,
      how_did_you_hear_about_us_other: currentUser.howDidYouHearAboutUsOther,
      agree_promotions_email: currentUser.agreePromotionsEmail,
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

  ////
  //// Misc.
  ////

  useEffect(() => {
    if (currentUser?.onboardingAnswers?.length > 0) {
      const userAnswers = {};
      currentUser.onboardingAnswers.forEach((answer) => {
        userAnswers[answer.questionId] = answer.answers || [];
      });
      setAnswers(userAnswers);
    }
  }, []);

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      fetchList();
    }

    return () => {
      mounted = false;
    };
  }, []);

  ////
  //// Component rendering.
  ////

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
  console.log("Beginning component rendering ...");
  
  return (
    <div className="section">
      {/* Page */}
      <nav className="mb-4">
        <div className="columns">
          <div className="column">
            <h1 className="title is-4">
              <FontAwesomeIcon className="fas" icon={faList} />
              &nbsp;Questions
            </h1>
          </div>
          <div className="column has-text-right">
            <button className="button is-success" onClick={handleSubmit}>
              Submit&nbsp;
            </button>
          </div>
        </div>

        {isFetching ? (
          <PageLoadingContent displayMessage={"Please wait..."} />
        ) : (
          <>
            <FormErrorBox errors={errors} />
            {listData &&
            listData.results &&
            (listData.results.length > 0 || previousCursors.length > 0) ? (
              <div>
                {listData.results.map((datum, i) => (
                  <div className="mb-6">
                    <div className="mb-3">
                      <h1 className="title has-text-centered">
                        <span className="mr-4">{i + 1}.</span>
                        {datum.title}
                      </h1>
                      {datum.subtitle && (
                        <h2 className="subtitle is-size-3 has-text-centered mb-5">
                          {datum.subtitle}
                        </h2>
                      )}
                    </div>
                    <div className="columns">
                      <div className="column">
                        {datum.options.map((option) => (
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
                  <button className="button is-success" onClick={handleSubmit}>
                    Submit&nbsp;
                  </button>
                </div>
              </div>
            ) : (
              <section className="hero is-medium has-background-white-ter">
                <div className="hero-body">
                  <p className="title">
                    <FontAwesomeIcon className="fas" icon={faQuestionCircle} />
                    &nbsp;No Questions Available
                  </p>
                  <p className="subtitle">
                    There are currently no questions available.&nbsp;
                  </p>
                </div>
              </section>
            )}
          </>
        )}
      </nav>
    </div>
  );
}

export default MemberQuestionnaireList;
