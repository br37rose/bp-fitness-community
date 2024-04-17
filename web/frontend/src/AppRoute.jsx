import { React } from "react";
import "bulma/css/bulma.min.css";
import "./css/styles.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { RecoilRoot } from "recoil";

import AdminNutritionPlanUpdate from "./Components/Admin/Member/NutritionPlan/Update";
import AdminNutritionPlanSubmissionForm from "./Components/Admin/Member/NutritionPlan/DetailSubmissionForm";
import AdminNutritionPlanDetail from "./Components/Admin/Member/NutritionPlan/Detail";
import AdminNutritionPlanList from "./Components/Admin/Member/NutritionPlan/List";
import AdminNutritionPlanAdd from "./Components/Admin/Member/NutritionPlan/Add";
import AdminFitnessPlanUpdate from "./Components/Admin/Member/FitnessPlan/Update";
import AdminFitnessPlanSubmissionForm from "./Components/Admin/Member/FitnessPlan/DetailSubmissionForm";
import AdminFitnessPlanDetail from "./Components/Admin/Member/FitnessPlan/Detail";
import AdminFitnessPlanList from "./Components/Admin/Member/FitnessPlan/List";
import AdminFitnessPlanAdd from "./Components/Admin/Member/FitnessPlan/Add";
import AdminOfferUpdate from "./Components/Admin/Offer/Update";
import AdminOfferDetail from "./Components/Admin/Offer/Detail";
import AdminOfferList from "./Components/Admin/Offer/List";
import AdminOfferAdd from "./Components/Admin/Offer/Add";
import AdminVideoContentUpdate from "./Components/Admin/VideoCollection/VideoContent/Update";
import AdminVideoContentDetail from "./Components/Admin/VideoCollection/VideoContent/Detail";
import AdminVideoContentAdd from "./Components/Admin/VideoCollection/VideoContent/Add";
import AdminVideoContentList from "./Components/Admin/VideoCollection/VideoContent/List";
import AdminVideoCollectionUpdate from "./Components/Admin/VideoCollection/Update";
import AdminVideoCollectionDetail from "./Components/Admin/VideoCollection/Detail";
import AdminVideoCollectionAdd from "./Components/Admin/VideoCollection/Add";
import AdminVideoCollectionList from "./Components/Admin/VideoCollection/List";
import AdminVideoCategoryUpdate from "./Components/Admin/VideoCategory/Update";
import AdminVideoCategoryDetail from "./Components/Admin/VideoCategory/Detail";
import AdminVideoCategoryList from "./Components/Admin/VideoCategory/List";
import AdminVideoCategoryAdd from "./Components/Admin/VideoCategory/Add";
import AdminExerciseUpdate from "./Components/Admin/Exercise/Update";
import AdminExerciseAdd from "./Components/Admin/Exercise/Add";
import AdminExerciseDetail from "./Components/Admin/Exercise/Detail";
import AdminExerciseList from "./Components/Admin/Exercise/List";
import AdminMemberTagList from "./Components/Admin/Member/DetailForTags";
import AdminMemberUpdate from "./Components/Admin/Member/Update";
import AdminMemberDetail from "./Components/Admin/Member/Detail";
import AdminMemberList from "./Components/Admin/Member/List";
import AdminMemberAdd from "./Components/Admin/Member/Add";
import AdminDashboard from "./Components/Admin/Dashboard";
import MemberLeaderboardGlobalTabularList from "./Components/Member/Biometric/Leaderboard/Global/TabularList";
import MemberSummary from "./Components/Member/Biometric/Summary/View";
import MemberDataPointHistoricalGraphicalList from "./Components/Member/Biometric/Historical/GraphicalList";
import MemberDataPointHistoricalTabularList from "./Components/Member/Biometric/Historical/TabularList";
import MemberBiometricLaunchpad from "./Components/Member/Biometric/Launchpad";
import MemberNutritionPlanUpdate from "./Components/Member/NutritionPlan/Update";
import MemberNutritionPlanSubmissionForm from "./Components/Member/NutritionPlan/DetailSubmissionForm";
import MemberNutritionPlanDetail from "./Components/Member/NutritionPlan/Detail";
import MemberNutritionPlanList from "./Components/Member/NutritionPlan/List";
import MemberNutritionPlanAdd from "./Components/Member/NutritionPlan/Add";
import MemberFitnessPlanUpdate from "./Components/Member/FitnessPlan/Update";
import MemberFitnessPlanSubmissionForm from "./Components/Member/FitnessPlan/DetailSubmissionForm";
import MemberFitnessPlanDetail from "./Components/Member/FitnessPlan/Detail";
import MemberFitnessPlanList from "./Components/Member/FitnessPlan/List";
import MemberFitnessPlanAdd from "./Components/Member/FitnessPlan/Add";
import MemberVideoContentDetail from "./Components/Member/VideoCollection/VideoContent/Detail";
import MemberVideoContentList from "./Components/Member/VideoCollection/VideoContent/List";
import MemberVideoCollectionDetail from "./Components/Member/VideoCollection/Detail";
import MemberVideoCollectionList from "./Components/Member/Videos/VideoCategories/Collections/List";
import MemberCategoriesList from "./Components/Member/Videos/VideoCategories/List";
import MemberExerciseDetail from "./Components/Member/Exercise/Detail";
import MemberExerciseList from "./Components/Member/Exercise/List";
import MemberDashboard from "./Components/Member/Dashboard";

import TwoFactorAuthenticationWizardStep1 from "./Components/Gateway/2FA/Step1";
import TwoFactorAuthenticationWizardStep2 from "./Components/Gateway/2FA/Step2";
import TwoFactorAuthenticationWizardStep3 from "./Components/Gateway/2FA/Step3";
import TwoFactorAuthenticationValidateOnLogin from "./Components/Gateway/2FA/ValidateOnLogin";
import LogoutRedirector from "./Components/Gateway/LogoutRedirector";
import Login from "./Components/Gateway/Login";
import Register from "./Components/Gateway/Register";
import RegisterSuccessful from "./Components/Gateway/RegisterSuccessful";
import Index from "./Components/Gateway/Index/Index";
import AnonymousCurrentUserRedirector from "./Components/Misc/AnonymousCurrentUserRedirector";
import TwoFactorAuthenticationRedirector from "./Components/Misc/TwoFactorAuthenticationRedirector";
import TopAlertBanner from "./Components/Misc/TopAlertBanner";
import Sidebar from "./Components/Menu/Sidebar";
import Topbar from "./Components/Menu/Top";
import NotFoundError from "./Components/Misc/NotFoundError";
import NotImplementedError from "./Components/Misc/NotImplementedError";
import AccountTwoFactorAuthenticationDetail from "./Components/Account/2FA/View";
import AccountEnableTwoFactorAuthentication from "./Components/Account/2FA/EnableView";
import EmailVerification from "./Components/Gateway/EmailVerification";
import AccountFriendList from "./Components/Account/Friend/Friend";
import AccountMoreOperationAvatar from "./Components/Account/More/Operation/Avatar/Avatar";
import AccountMoreOperationSubscribe from "./Components/Account/More/Operation/Subscribe/Subscribe";
import AccountMoreLaunchpad from "./Components/Account/More/Launchpad";
import AccountWearableTechLaunchpad from "./Components/Account/WearableTech/Launchpad";
import AccountSubscriptionDetailAndCancel from "./Components/Account/Subscription/Subscription";
import AccountInvoiceList from "./Components/Account/Subscription/Invoice/List";
import AccountTagList from "./Components/Account/DetailForTags";
import AccountDetail from "./Components/Account/Detail";
import AccountUpdate from "./Components/Account/Update";
import AccountChangePassword from "./Components/Account/ChangePassword";
import ForgotPassword from "./Components/Gateway/ForgotPassword";
import PasswordReset from "./Components/Gateway/PasswordReset";
import MemberVideoCollectionContentList from "./Components/Member/Videos/VideoCategories/Collections/VideoContent/List";
import MemberLeaderboardPersonal from "./Components/Member/Biometric/Leaderboard/Personal/Personal";
// import MemberHistoricalDashboard from "./Components/Member/Biometric/Historical/Dashboard";
import PaymentProcessoePurchaseCanceled from "./Components/Member/PaymentProcessor/PurchaseCanceled";
import PaymentProcessorPurchaseSuccess from "./Components/Member/PaymentProcessor/PurchaseSuccess";
import UserProfile from "./Components/Account/User/UserProfile";
import AdminTrainingProgramList from "./Components/Admin/TrainingProgram/list";
import AdminTrainingProgramAdd from "./Components/Admin/TrainingProgram/add";
import AdminWokoutList from "./Components/Admin/Workouts/list";
import AdminWorkoutAdd from "./Components/Admin/Workouts/add";
import AdminWorkoutDetail from "./Components/Admin/Workouts/Detail";
import AdminWorkoutUpdate from "./Components/Admin/Workouts/update";
import MemberWorkoutList from "./Components/Member/Workouts/list";
import MemberWorkoutAdd from "./Components/Member/Workouts/add";
import MemberWorkoutDetail from "./Components/Member/Workouts/Detail";
import MemberWorkoutEdit from "./Components/Member/Workouts/update";
import AdminTPDetail from "./Components/Admin/TrainingProgram/Detail";
import MemberTrainingProgramList from "./Components/Member/TrainingProgram/list";
import MemberTrainingProgramAdd from "./Components/Member/TrainingProgram/add";
import MemberTPDetail from "./Components/Member/TrainingProgram/Detail";

function AppRoute() {
  return (
    <div class="is-widescreen">
      <RecoilRoot>
        <Router>
          <AnonymousCurrentUserRedirector />
          <TwoFactorAuthenticationRedirector />
          <TopAlertBanner />
          <Topbar />
          <div class="columns pt-3">
            <Sidebar />
            <div class="column">
              <section class="main-content columns is-fullheight">
                <Routes>
                  {/*
                                        -----------------------------------------------
                                        EVERYTHING BELOW BELONGS TO THE ADMINISTRATION.
                                        -----------------------------------------------
                                    */}

                  <Route
                    exact
                    path="/admin/offer/:id/update"
                    element={<AdminOfferUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/offer/:id"
                    element={<AdminOfferDetail />}
                  />
                  <Route
                    exact
                    path="/admin/offers/add"
                    element={<AdminOfferAdd />}
                  />
                  <Route
                    exact
                    path="/admin/offers"
                    element={<AdminOfferList />}
                  />

                  <Route
                    exact
                    path="/admin/video-collection/:vcid/video-content/:vconid/update"
                    element={<AdminVideoContentUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/video-collection/:vcid/video-content/:vconid"
                    element={<AdminVideoContentDetail />}
                  />
                  <Route
                    exact
                    path="/admin/video-collection/:vcid/video-contents/add"
                    element={<AdminVideoContentAdd />}
                  />
                  <Route
                    exact
                    path="/admin/video-collection/:vcid/video-contents"
                    element={<AdminVideoContentList />}
                  />

                  <Route
                    exact
                    path="/admin/video-collection/:vcid/update"
                    element={<AdminVideoCollectionUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/video-collection/:vcid"
                    element={<AdminVideoCollectionDetail />}
                  />
                  <Route
                    exact
                    path="/admin/video-collections/add"
                    element={<AdminVideoCollectionAdd />}
                  />
                  <Route
                    exact
                    path="/admin/video-collections"
                    element={<AdminVideoCollectionList />}
                  />

                  <Route
                    exact
                    path="/admin/video-category/:id/update"
                    element={<AdminVideoCategoryUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/video-category/:id"
                    element={<AdminVideoCategoryDetail />}
                  />
                  <Route
                    exact
                    path="/admin/video-categories/add"
                    element={<AdminVideoCategoryAdd />}
                  />
                  <Route
                    exact
                    path="/admin/video-categories"
                    element={<AdminVideoCategoryList />}
                  />

                  <Route
                    exact
                    path="/admin/exercise/:id/update"
                    element={<AdminExerciseUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/exercise/:id"
                    element={<AdminExerciseDetail />}
                  />
                  <Route
                    exact
                    path="/admin/exercises/add"
                    element={<AdminExerciseAdd />}
                  />
                  <Route
                    exact
                    path="/admin/exercises"
                    element={<AdminExerciseList />}
                  />

                  <Route
                    exact
                    path="/admin/member/:uid/nutrition-plan/:id/update"
                    element={<AdminNutritionPlanUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/member/:uid/nutrition-plan/:id/submission-form"
                    element={<AdminNutritionPlanSubmissionForm />}
                  />
                  <Route
                    exact
                    path="/admin/member/:uid/nutrition-plan/:id"
                    element={<AdminNutritionPlanDetail />}
                  />
                  <Route
                    exact
                    path="/admin/member/:uid/nutrition-plans/add"
                    element={<AdminNutritionPlanAdd />}
                  />
                  <Route
                    exact
                    path="/admin/member/:uid/nutrition-plans"
                    element={<AdminNutritionPlanList />}
                  />

                  <Route
                    exact
                    path="/admin/member/:uid/fitness-plan/:id/update"
                    element={<AdminFitnessPlanUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/member/:uid/fitness-plan/:id/submission-form"
                    element={<AdminFitnessPlanSubmissionForm />}
                  />
                  <Route
                    exact
                    path="/admin/member/:uid/fitness-plan/:id"
                    element={<AdminFitnessPlanDetail />}
                  />
                  <Route
                    exact
                    path="/admin/member/:uid/fitness-plans/add"
                    element={<AdminFitnessPlanAdd />}
                  />
                  <Route
                    exact
                    path="/admin/member/:uid/fitness-plans"
                    element={<AdminFitnessPlanList />}
                  />

                  <Route
                    exact
                    path="/admin/member/:id/update"
                    element={<AdminMemberUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/member/:id/tags"
                    element={<AdminMemberTagList />}
                  />
                  <Route
                    exact
                    path="/admin/member/:id"
                    element={<AdminMemberDetail />}
                  />
                  <Route
                    exact
                    path="/admin/members"
                    element={<AdminMemberList />}
                  />
                  <Route
                    exact
                    path="/admin/members/add"
                    element={<AdminMemberAdd />}
                  />

                  <Route
                    exact
                    path="/admin/dashboard"
                    element={<AdminDashboard />}
                  />
                  <Route
                    exact
                    path="/admin/training-program"
                    element={<AdminTrainingProgramList />}
                  />
                  <Route
                    exact
                    path="/admin/training-program/add"
                    element={<AdminTrainingProgramAdd />}
                  />
                  <Route
                    exact
                    path="/admin/training-program/:id"
                    element={<AdminTPDetail />}
                  />
                  <Route
                    exact
                    path="/admin/workouts"
                    element={<AdminWokoutList />}
                  />
                  <Route
                    exact
                    path="/admin/workouts/add"
                    element={<AdminWorkoutAdd />}
                  />
                  <Route
                    exact
                    path="/admin/workouts/:id"
                    element={<AdminWorkoutDetail />}
                  />
                  <Route
                    exact
                    path="/admin/workouts/:id/update"
                    element={<AdminWorkoutUpdate />}
                  />

                  {/*
                                        ---------------------------------------
                                        EVERYTHING BELOW BELONGS TO THE MEMBER.
                                        ---------------------------------------
                                    */}

                  <Route
                    exact
                    path="/biometrics/leaderboard/global"
                    element={<MemberLeaderboardGlobalTabularList />}
                  />
                  <Route
                    exact
                    path="/biometrics/leaderboard/personal"
                    element={<MemberLeaderboardPersonal />}
                  />
                  <Route
                    exact
                    path="/biometrics/summary"
                    element={<MemberSummary />}
                  />
                  <Route
                    exact
                    path="/biometrics/history/graphview"
                    element={<MemberDataPointHistoricalGraphicalList />}
                  />
                  <Route
                    exact
                    path="/biometrics/history/tableview"
                    element={<MemberDataPointHistoricalTabularList />}
                  />
                  <Route
                    exact
                    path="/biometrics"
                    element={<MemberBiometricLaunchpad />}
                  />

                  <Route
                    exact
                    path="/nutrition-plan/:id/update"
                    element={<MemberNutritionPlanUpdate />}
                  />
                  <Route
                    exact
                    path="/nutrition-plan/:id/submission-form"
                    element={<MemberNutritionPlanSubmissionForm />}
                  />
                  <Route
                    exact
                    path="/nutrition-plan/:id"
                    element={<MemberNutritionPlanDetail />}
                  />
                  <Route
                    exact
                    path="/nutrition-plans/add"
                    element={<MemberNutritionPlanAdd />}
                  />
                  <Route
                    exact
                    path="/nutrition-plans"
                    element={<MemberNutritionPlanList />}
                  />

                  <Route
                    exact
                    path="/fitness-plan/:id/update"
                    element={<MemberFitnessPlanUpdate />}
                  />
                  <Route
                    exact
                    path="/fitness-plan/:id/submission-form"
                    element={<MemberFitnessPlanSubmissionForm />}
                  />
                  <Route
                    exact
                    path="/fitness-plan/:id"
                    element={<MemberFitnessPlanDetail />}
                  />
                  <Route
                    exact
                    path="/fitness-plans/add"
                    element={<MemberFitnessPlanAdd />}
                  />
                  <Route
                    exact
                    path="/fitness-plans"
                    element={<MemberFitnessPlanList />}
                  />

                  {/* <Route exact path="/video-collection/:vcid/video-content/:vconid" element={<MemberVideoContentDetail />} /> */}
                  <Route
                    exact
                    path="/video-collection/:vcid/video-content/:vconid"
                    element={<MemberVideoCollectionContentList />}
                  />

                  <Route
                    exact
                    path="/video-collection/:vcid/video-contents"
                    element={<MemberVideoContentList />}
                  />

                  <Route
                    exact
                    path="/video-collection/:vcid"
                    element={<MemberVideoCollectionDetail />}
                  />
                  <Route
                    exact
                    path="/video-categories"
                    element={<MemberCategoriesList />}
                  />
                  <Route
                    exact
                    path="/video-category/:vcatid/video-collections"
                    element={<MemberVideoCollectionList />}
                  />

                  <Route
                    exact
                    path="/exercise/:id"
                    element={<MemberExerciseDetail />}
                  />
                  <Route
                    exact
                    path="/exercises"
                    element={<MemberExerciseList />}
                  />

                  <Route
                    exact
                    path="/dashboard"
                    element={<MemberDashboard />}
                  />
                  <Route
                    exact
                    path="/workouts"
                    element={<MemberWorkoutList />}
                  />
                  <Route
                    exact
                    path="/workouts/add"
                    element={<MemberWorkoutAdd />}
                  />
                  <Route
                    exact
                    path="/workouts/:id"
                    element={<MemberWorkoutDetail />}
                  />
                  <Route
                    exact
                    path="/workouts/:id/update"
                    element={<MemberWorkoutEdit />}
                  />

                  <Route
                    exact
                    path="/training-program"
                    element={<MemberTrainingProgramList />}
                  />
                  <Route
                    exact
                    path="/training-program/add"
                    element={<MemberTrainingProgramAdd />}
                  />
                  <Route
                    exact
                    path="/training-program/:id"
                    element={<MemberTPDetail />}
                  />

                  {/*
                                        -----------------------------------------------
                                        EVERYTHING BELOW BLONGS TO USER PROFILE.
                                        -----------------------------------------------
                                    */}

                  <Route
                    exact
                    path="/purchase/success"
                    element={<PaymentProcessorPurchaseSuccess />}
                  />
                  <Route
                    exact
                    path="/purchase/canceled"
                    element={<PaymentProcessoePurchaseCanceled />}
                  />
                  <Route
                    exact
                    path="/account/more/subscribe"
                    element={<AccountMoreOperationSubscribe />}
                  />
                  <Route
                    exact
                    path="/account/more/avatar"
                    element={<AccountMoreOperationAvatar />}
                  />
                  <Route
                    exact
                    path="/account/more"
                    element={<AccountMoreLaunchpad />}
                  />
                  <Route
                    exact
                    path="/account/wearable-tech"
                    element={<AccountWearableTechLaunchpad />}
                  />
                  <Route
                    exact
                    path="/account/friends"
                    element={<AccountFriendList />}
                  />
                  <Route
                    exact
                    path="/account/tags"
                    element={<AccountTagList />}
                  />
                  <Route exact path="/account" element={<UserProfile />} />
                  <Route
                    exact
                    path="/account/update"
                    element={<AccountUpdate />}
                  />
                  <Route
                    exact
                    path="/account/change-password"
                    element={<AccountChangePassword />}
                  />
                  <Route
                    exact
                    path="/account/subscription"
                    element={<AccountSubscriptionDetailAndCancel />}
                  />
                  <Route
                    exact
                    path="/account/subscription/invoices"
                    element={<AccountInvoiceList />}
                  />
                  <Route
                    exact
                    path="/account/2fa"
                    element={<AccountTwoFactorAuthenticationDetail />}
                  />
                  <Route
                    exact
                    path="/account/2fa/enable"
                    element={<AccountEnableTwoFactorAuthentication />}
                  />

                  {/*
                                        -----------------------------------------------
                                        EVERYTHING BELOW BELONGS TO THE GATEWAY OR INDEX.
                                        -----------------------------------------------
                                     */}
                  <Route exact path="/register" element={<Register />} />
                  <Route
                    exact
                    path="/register-successful"
                    element={<RegisterSuccessful />}
                  />
                  <Route exact path="/login" element={<Login />} />
                  <Route
                    exact
                    path="/login/2fa/step-1"
                    element={<TwoFactorAuthenticationWizardStep1 />}
                  />
                  <Route
                    exact
                    path="/login/2fa/step-2"
                    element={<TwoFactorAuthenticationWizardStep2 />}
                  />
                  <Route
                    exact
                    path="/login/2fa/step-3"
                    element={<TwoFactorAuthenticationWizardStep3 />}
                  />
                  <Route
                    exact
                    path="/login/2fa"
                    element={<TwoFactorAuthenticationValidateOnLogin />}
                  />
                  <Route exact path="/logout" element={<LogoutRedirector />} />
                  <Route exact path="/verify" element={<EmailVerification />} />
                  <Route
                    exact
                    path="/forgot-password"
                    element={<ForgotPassword />}
                  />
                  <Route
                    exact
                    path="/password-reset"
                    element={<PasswordReset />}
                  />
                  <Route exact path="/" element={<Index />} />
                  <Route path="*" element={<NotFoundError />} />
                </Routes>
              </section>
              <div>
                {/* DEVELOPERS NOTE: Mobile tab-bar menu can go here */}
              </div>
              <footer class="footer is-hidden">
                <div class="container">
                  <div class="content has-text-centered">
                    <p>Hello</p>
                  </div>
                </div>
              </footer>
            </div>
          </div>
        </Router>
      </RecoilRoot>
    </div>
  );
}

export default AppRoute;
