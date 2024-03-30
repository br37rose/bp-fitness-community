import React, { useState } from 'react';
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faArrowLeftLong, faArrowRightLong } from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { quizAnswersState, currentUserState } from "../../../AppState"; // Ensure this path matches your project structure

// Assuming Card and Title are imported correctly above this component

const OnBoardingQuestionWizard = ({ questions }) => {
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
            alert('Wizard Finished!');
            // Here you can implement additional logic upon finishing the wizard
        }
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
                ? answers[questionId].filter(id => id !== selectedId)
                : [...(answers[questionId] || []), selectedId];
            setAnswers({ ...answers, [questionId]: updatedSelections });
        } else {
            setAnswers({ ...answers, [questionId]: selectedId });
        }
    };

    const renderQuestionContent = () => {
        const currentQuestion = questions[currentQuestionIndex];
        if (!currentQuestion || !currentQuestion.content) {
            return <div>Question data is incomplete or missing.</div>;
        }
        // Assuming content is a function that returns JSX
        return currentQuestion.content({
            onSelect: handleSelect,
            selectedAnswers: answers[currentQuestion.questionId],
            isMultiSelect: currentQuestion.isMultiSelect
        });
    };

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
                    <div className="progress-container" style={{ width: '100%', backgroundColor: '#f5f5f5' }}>
                        <div
                            className="progress-bar"
                            style={{
                                height: '4px',
                                width: `${progress}%`,
                                backgroundColor: '#735827',
                                transition: 'width 0.5s'
                            }}
                        ></div>
                    </div>
                </header>
                <section className="modal-card-body has-background-dark">
                    {/* Render the current step component here */}
                    <div>{renderQuestionContent()}</div>
                </section>
                <footer className="modal-card-foot is-flex is-justify-content-space-between has-background-dark">
                    <button className="button is-danger" onClick={handlePrevious} disabled={isFirstQuestion}>
                        <FontAwesomeIcon icon={faArrowLeftLong} />&nbsp;Previous
                    </button>
                    {isLastQuestion ? (
                        <button className="button is-success" onClick={handleNext}>
                            Finish&nbsp;<FontAwesomeIcon icon={faArrowRightLong} />
                        </button>
                    ) : (
                        <button className="button is-info" onClick={handleNext}>
                            Next&nbsp;<FontAwesomeIcon icon={faArrowRightLong} />
                        </button>
                    )}
                </footer>
            </div>
        </div>
    );
};

export default OnBoardingQuestionWizard;
