import {
  RecoilRoot,
  atom,
  selector,
  useRecoilState,
  useRecoilValue,
} from "recoil";
import { recoilPersist } from "recoil-persist";

////
//// NON-PERSISTENT STATE BELOW
////

// Control whether the hamburer menu icon was clicked or not. This state is
// needed by 'TopNavigation' an 'SideNavigation' components.
export const onHamburgerClickedState = atom({
  key: "onHamburgerClicked", // unique ID (with respect to other atoms/selectors)
  default: false, // default value (aka initial value)
});

// Control what message to display at the top as a banner in the app.
export const topAlertMessageState = atom({
  key: "topBannerAlertMessage",
  default: "",
});

// Control what type of message to display at the top as a banner in the app.
export const topAlertStatusState = atom({
  key: "topBannerAlertStatus",
  default: "success",
});

export const quizAnswersState = atom({
  key: "quizAnswersState", // unique ID (with respect to other atoms/selectors)
  default: [], // default value (aka initial value)
});

////
//// PERSISTENT STATE BELOW
////
//
// https://github.com/polemius/recoil-persist
//

const { persistAtom } = recoilPersist();
export const currentUserState = atom({
  key: "currentUser",
  default: null,
  effects_UNSTABLE: [persistAtom],
});

export const currentOTPResponseState = atom({
  key: "currentOTPResponse",
  default: null,
  effects_UNSTABLE: [persistAtom],
});

export const workoutProgramDetailState = atom({
  key: "workoutProgramDetail",
  default: null,
  effects_UNSTABLE: [persistAtom],
});

export const currentWorkoutSessionState = atom({
  key: "currentWorkoutSession",
  default: null,
  effects_UNSTABLE: [persistAtom],
});

// --- Offers --- //

// Control whether to show filters for the list.
export const offersFilterShowState = atom({
  key: "offersFilterShowState",
  default: false,
  effects_UNSTABLE: [persistAtom],
});

export const offersFilterTemporarySearchTextState = atom({
  key: "offersFilterTemporarySearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const offersFilterActualSearchTextState = atom({
  key: "offersFilterActualSearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const offersFilterStatusState = atom({
  key: "offersFilterStatusState",
  default: 2,
  effects_UNSTABLE: [persistAtom],
});

// --- Admin Members List --- //

// Control whether to show filters for the list.
export const membersFilterShowState = atom({
  key: "membersFilterShowState",
  default: false,
  effects_UNSTABLE: [persistAtom],
});

export const membersFilterTemporarySearchTextState = atom({
  key: "membersFilterTemporarySearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const membersFilterActualSearchTextState = atom({
  key: "membersFilterActualSearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const membersFilterOfferIDState = atom({
  key: "membersFilterOfferIDState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const membersFilterStatusState = atom({
  key: "membersFilterStatusState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const membersFilterSortState = atom({
  key: "membersFilterSortState",
  default: "name,-1",
  effects_UNSTABLE: [persistAtom],
});

// --- Admin Exercise List --- //

// Control whether to show filters for the list.
export const exercisesFilterShowState = atom({
  key: "exercisesFilterShowState",
  default: false,
  effects_UNSTABLE: [persistAtom],
});

export const exercisesFilterTemporarySearchTextState = atom({
  key: "exercisesFilterTemporarySearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const exercisesFilterActualSearchTextState = atom({
  key: "exercisesFilterActualSearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const exercisesFilterCategoryState = atom({
  key: "exercisesFilterCategoryState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const exercisesFilterMovementTypeState = atom({
  key: "exercisesFilterMovementTypeState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const exercisesFilterStatusState = atom({
  key: "exercisesFilterStatusState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const exercisesFilterGenderState = atom({
  key: "exercisesFilterGenderState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const exercisesFilterVideoTypeState = atom({
  key: "exercisesFilterVideoTypeState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const exercisesFilterSortState = atom({
  key: "exercisesFilterSortState",
  default: "created,-1",
  effects_UNSTABLE: [persistAtom],
});

// --- Admin Video Collection List --- //

// Control whether to show filters for the list.
export const videoCollectionsFilterShowState = atom({
  key: "videoCollectionsFilterShowState",
  default: false,
  effects_UNSTABLE: [persistAtom],
});

export const videoCollectionsFilterTemporarySearchTextState = atom({
  key: "videoCollectionsFilterTemporarySearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const videoCollectionsFilterActualSearchTextState = atom({
  key: "videoCollectionsFilterActualSearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const videoCollectionsFilterStatusState = atom({
  key: "videoCollectionsFilterStatusState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const videoCollectionsFilterVideoTypeState = atom({
  key: "videoCollectionsFilterVideoTypeState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const videoCollectionsFilterSortState = atom({
  key: "evideoCollectionsFilterSortState",
  default: "created,-1",
  effects_UNSTABLE: [persistAtom],
});

// --- Admin Video Content List --- //

// Control whether to show filters for the list.
export const videoContentsFilterShowState = atom({
  key: "videoContentsFilterShowState",
  default: false,
  effects_UNSTABLE: [persistAtom],
});

export const videoContentsFilterTemporarySearchTextState = atom({
  key: "videoContentsFilterTemporarySearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const videoContentsFilterActualSearchTextState = atom({
  key: "videoContentsFilterActualSearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const videoContentsFilterStatusState = atom({
  key: "videoContentsFilterStatusState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const videoContentsFilterVideoTypeState = atom({
  key: "videoContentsFilterVideoTypeState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const videoContentsFilterSortState = atom({
  key: "videoContentsFilterSortState",
  default: "created,-1",
  effects_UNSTABLE: [persistAtom],
});

export const videoContentsFilterOfferIDState = atom({
  key: "videoContentsFilterOfferIDState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const videoContentsFilterCategoryIDState = atom({
  key: "videoContentsFilterCategoryIDState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

// --- Admin Video Category List --- //

// Control whether to show filters for the list.
export const videoCategoryFilterShowState = atom({
  key: "videoCategoryFilterShowState",
  default: false,
  effects_UNSTABLE: [persistAtom],
});

export const videoCategoryFilterTemporarySearchTextState = atom({
  key: "videoCategoryFilterTemporarySearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const videoCategoryFilterActualSearchTextState = atom({
  key: "videoCategoryFilterActualSearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const videoCategoryFilterSortState = atom({
  key: "videoCategoryFilterSortState",
  default: "no,-1",
  effects_UNSTABLE: [persistAtom],
});

export const videoCategoryFilterStatusState = atom({
  key: "videoCategoryFilterStatusState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

// --- Admin Data Point List --- //

// Control whether to show filters for the list.
export const dataPointFilterShowState = atom({
  key: "dataPointFilterShowState",
  default: false,
  effects_UNSTABLE: [persistAtom],
});

export const dataPointFilterTemporarySearchTextState = atom({
  key: "dataPointFilterTemporarySearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const dataPointFilterActualSearchTextState = atom({
  key: "dataPointFilterActualSearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const dataPointFilterSortState = atom({
  key: "dataPointFilterSortState",
  default: "timestamp,DESC",
  effects_UNSTABLE: [persistAtom],
});

export const dataPointFilterStatusState = atom({
  key: "dataPointFilterStatusState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const dataPointFilterIsHeartRateState = atom({
  key: "dataPointFilterIsHeartRate",
  default: true,
  effects_UNSTABLE: [persistAtom],
});

export const dataPointFilterIsStepsCounterState = atom({
  key: "dataPointFilterIsStepsCounter",
  default: true,
  effects_UNSTABLE: [persistAtom],
});

export const questionnaireFilterStatus = atom({
  key: "questionnaireFilterStatus",
  default: 1,
  effects_UNSTABLE: [persistAtom],
});

export const questionnaireFilterShowState = atom({
  key: "questionnaireFilterShowState",
  default: false,
  effects_UNSTABLE: [persistAtom],
});

export const questionnaireFilterSortState = atom({
  key: "questionnaireFilterSortState",
  default: "created,-1",
  effects_UNSTABLE: [persistAtom],
});
