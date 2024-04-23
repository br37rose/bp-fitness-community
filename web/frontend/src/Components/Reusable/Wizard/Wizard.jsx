import React, { useEffect, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeftLong,
  faArrowRightLong,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import {
  quizAnswersState,
  currentUserState,
  topAlertMessageState,
  topAlertStatusState,
} from "../../../AppState"; // Ensure this path matches your project structure
import { getQuestionnaireListApi } from "../../../API/questionnaire";
import Scroll from "react-scroll";
import PageLoadingContent from "../PageLoadingContent";
import { SelectableOption, Title } from "./Questions";
import { putMemberUpdateAPI } from "../../../API/member";
import { Navigate } from "react-router-dom";

// Assuming Card and Title are imported correctly above this component

export const Onboarding = () => {
  const [currentUser] = useRecoilState(currentUserState);
  const [isFetching, setFetching] = useState(false);
  const [errors, setErrors] = useState({});
  const [listData, setListData] = useState(null);
  const [forceURL, setForceURL] = useState("");

  function OnListSuccess(response) {
    if (response && response.results && response.results.length > 0) {
      console.log("OnListSuccess");
      setListData(response);
      //   if (response.hasNextPage) {
      //     setNextCursor(response.nextCursor); // For pagination purposes.
      //   }
    } else {
      setListData(null);
      setForceURL("/dashboard");
      //   setNextCursor("");
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

  const fetchList = (cur, limit, keywords, st, sbv) => {
    setFetching(true);
    setErrors({});

    let params = new Map();
    params.set("status", true);
    getQuestionnaireListApi(params, OnListSuccess, OnListErr, onListDone);
  };

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

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      {isFetching ? (
        <PageLoadingContent displayMessage={"Please wait..."} />
      ) : (
        <>
          {listData && listData.results.length > 0 && (
            <OnBoardingQuestionWizard questions={listData.results} />
          )}
        </>
      )}
    </>
  );
};

export const OnBoardingQuestionWizard = ({ questions }) => {
  //

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  const [forceURL, setForceURL] = useState("");
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);

  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [answers, setAnswers] = useRecoilState(quizAnswersState);
  const [currentUser] = useRecoilState(currentUserState);

  const isLastQuestion = currentQuestionIndex === questions.length - 1;
  const isFirstQuestion = currentQuestionIndex === 0;
  const progress = ((currentQuestionIndex + 1) / questions.length) * 100;

  const handleNext = () => {
    if (!isLastQuestion) {
      setCurrentQuestionIndex(currentQuestionIndex + 1);
    } else {
      console.log("answers", answers);
      alert("Wizard Finished!");
      // Here you can implement additional logic upon finishing the wizard
    }
  };

  function onAdminMemberUpdateSuccess(response) {
    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Member updated");
    setTopAlertStatus("Workout Member");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Redirect the user to a new page.
    setForceURL("/admin/member/" + response.id);
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
  const handleSubmit = () => {
    setFetching(true);
    setErrors({});

    const onboardingAnswers = questions.map((question) => ({
      question_id: question.id,
      answers: answers[question.id] || [],
    }));
    const decamelizedData = {
      id: currentUser.id,
      organization_id: currentUser.organizationID,
      first_name: currentUser.firstName,
      last_name: currentUser.lastName,
      email: currentUser.email,
      phone: currentUser.phone,
      postal_code: currentUser.postalCode,
      address_line_1: currentUser.addressLine1,
      address_line_2: currentUser.addressLine2,
      city: currentUser.city,
      region: currentUser.region,
      country: currentUser.country,
      status: currentUser.status,
      password: currentUser.password,
      password_repeated: currentUser.passwordRepeated,
      how_did_you_hear_about_us: currentUser.howDidYouHearAboutUs,
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

  const handlePrevious = () => {
    if (!isFirstQuestion) {
      setCurrentQuestionIndex(currentQuestionIndex - 1);
    }
  };

  // Adjusted to correctly handle selection with Recoil state
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

  // const renderQuestionContent = () => {
  //   const currentQuestion = questions[currentQuestionIndex];
  //   if (!currentQuestion || !currentQuestion.content) {
  //     return <div>Question data is incomplete or missing.</div>;
  //   }
  //   // Assuming content is a function that returns JSX
  //   return currentQuestion.content({
  //     onSelect: handleSelect,
  //     selectedAnswers: answers[currentQuestion.questionId],
  //     isMultiSelect: currentQuestion.isMultiSelect,
  //   });
  // };

  const renderQuestionContent = () => {
    const currentQuestion = questions[currentQuestionIndex];
    if (!currentQuestion) {
      return <div>Question data is incomplete or missing.</div>;
    }

    return (
      <div>
        <Title
          text={currentQuestion.title}
          subtitle={currentQuestion.subtitle}
        />
        <div className="columns">
          <div className="column">
            {currentQuestion.options.map((option) => (
              <SelectableOption
                key={option}
                option={option}
                isSelected={
                  Array.isArray(answers[currentQuestion.id])
                    ? answers[currentQuestion.id].includes(option)
                    : answers[currentQuestion.id] === option
                }
                onSelect={() =>
                  handleSelect(
                    currentQuestion.id,
                    option,
                    currentQuestion.isMultiselect
                  )
                }
              />
            ))}
          </div>
          <div className="column">
            <figure className="image py-0 px-6">
              <img
                src="https://www.transparentlabs.com/cdn/shop/articles/how_long_muscle_1200x1200.jpg?v=1602607728"
                alt="Fitness"
                style={{ maxHeight: "70vh" }}
              />
            </figure>
          </div>
        </div>
      </div>
    );
  };

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <div className="modal is-active">
      <div className="modal-background"></div>
      <div className="modal-card is-fullwidth is-fullheight">
        <header className="modal-card-head has-background-dark">
          <nav
            class="navbar is-dark"
            role="navigation"
            aria-label="main navigation"
          >
            <div class="navbar-brand">
              {currentUser && (
                <>
                  {currentUser.role === 1 && <>TODO - SaaSify</>}
                  {currentUser.role === 2 && (
                    <a
                      class="navbar-item"
                      href="/admin/dashboard"
                      style={{ color: "white" }}
                    >
                      <img
                        src="/static/logo.png"
                        style={{
                          maxWidth: "8rem",
                          maxHeight: "3rem",
                          width: "100%", // This ensures the image scales within the given maxWidth and maxHeight
                          height: "auto", // This maintains the aspect ratio of the image
                        }}
                        alt="BP8_Fitness_Logo"
                      />
                    </a>
                  )}
                  {currentUser.role === 3 && (
                    <a
                      class="navbar-item"
                      href="/dashboard"
                      style={{ color: "white" }}
                    >
                      <img
                        src="/static/logo.png"
                        style={{
                          maxWidth: "8rem",
                          maxHeight: "3rem",
                          width: "100%", // This ensures the image scales within the given maxWidth and maxHeight
                          height: "auto", // This maintains the aspect ratio of the image
                        }}
                        alt="BP8_Fitness_Logo"
                      />
                    </a>
                  )}
                  {currentUser.role === 4 && (
                    <a
                      class="navbar-item"
                      href="/dashboard"
                      style={{ color: "white" }}
                    >
                      <img
                        src="/static/logo.png"
                        style={{
                          maxWidth: "8rem",
                          maxHeight: "3rem",
                          width: "100%", // This ensures the image scales within the given maxWidth and maxHeight
                          height: "auto", // This maintains the aspect ratio of the image
                        }}
                        alt="BP8_Fitness_Logo"
                      />
                    </a>
                  )}
                </>
              )}
            </div>
          </nav>
          <div
            className="progress-container"
            style={{ width: "100%", backgroundColor: "#f5f5f5" }}
          >
            <div
              className="progress-bar"
              style={{
                height: "4px",
                width: `${progress}%`,
                backgroundColor: "#735827",
                transition: "width 0.5s",
              }}
            ></div>
          </div>
        </header>
        <section className="modal-card-body has-background-dark">
          {/* Render the current step component here */}
          <div>{renderQuestionContent()}</div>
        </section>
        <footer className="modal-card-foot is-flex is-justify-content-space-between has-background-dark">
          <button
            className="button is-danger"
            onClick={handlePrevious}
            disabled={isFirstQuestion}
          >
            <FontAwesomeIcon icon={faArrowLeftLong} />
            &nbsp;Previous
          </button>
          {isLastQuestion ? (
            <button className="button is-success" onClick={handleSubmit}>
              Finish&nbsp;
              <FontAwesomeIcon icon={faArrowRightLong} />
            </button>
          ) : (
            <button className="button is-info" onClick={handleNext}>
              Next&nbsp;
              <FontAwesomeIcon icon={faArrowRightLong} />
            </button>
          )}
        </footer>
      </div>
    </div>
  );
};

export default Onboarding;
